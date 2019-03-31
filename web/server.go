package web

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/tupass/tupass-backend/api"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
)

//localBuild specifies whether the frontend should be included into the Go build
var localBuild = "false"

// StartServer starts a http server accepting incoming requests.
func StartServer(serverPort string) {

	router := mux.NewRouter()

	// set requestHandler as handler for every incoming API request
	router.HandleFunc("/api/", api.RequestHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api", api.RequestHandler).Methods("GET", "OPTIONS")

	if localBuild == "true" {
		//handle language redirection
		router.HandleFunc("/", RedirectLanguageHandler)

		//server frontend for en/de
		box := rice.MustFindBox("frontend").HTTPBox()
		fileServer := http.FileServer(box)
		router.PathPrefix("/").Handler(fileServer)
		log.Println("Serving locally bundled frontend")

		//open browser after frontend is ready, only when '-d' flag was not specified
		daemonized := flag.Bool("d", false, "when set browser will not open if started")
		flag.Parse()
		if !*daemonized {
			//Wait for server to start and open browser
			go func() {
				for {
					time.Sleep(time.Second)

					resp, err := http.Get("http://localhost:" + serverPort)
					if err != nil {
						log.Println("Failed to reach webserver:", err)
						continue
					}

					err = resp.Body.Close()
					if err != nil {
						log.Println("Unable to close response for request to webserver")
					}

					if resp.StatusCode != http.StatusOK {
						log.Println("Not OK:", resp.StatusCode)
						continue
					}

					// Reached this point: server is up and running!
					break
				}

				log.Println("Opening browser for TUPass")
				err := open("http://localhost:" + serverPort)
				if err != nil {
					log.Println("Unable to open browser")
				}
			}()
		}
	}

	// construct the server
	s := &http.Server{
		Handler:        router,
		Addr:           "127.0.0.1:" + serverPort, // only on localhost (prod gets proxied by nginx)
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1024, // limit header to 1KB to prevent requests with too long passwords (could exceed memory)
	}

	// start serving
	log.Printf("Starting TUPass API Server on port %s\n", serverPort)
	log.Fatal(s.ListenAndServe())
}

//RedirectLanguageHandler takes incoming requests and redirects them based on the users language
func RedirectLanguageHandler(w http.ResponseWriter, r *http.Request) {
	serverLangs := []language.Tag{
		language.English, // en fallback
		language.German,  // de
	}
	matcher := language.NewMatcher(serverLangs)

	lang, _ := r.Cookie("lang")
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, lang.String(), accept)

	langCode := tag.String()[:2]
	target := "http://" + r.Host + "/" + langCode
	http.Redirect(w, r, target, http.StatusFound)
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
