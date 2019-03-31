package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tupass/tupass-backend/metric"
)

//validateInput only returns true if the given password is valid (contains only ASCII characters and is not empty nor too long)
func validateInput(pw string) bool {
	notASCII := func(r rune) bool {
		return r < ' ' || r > '~'
	}

	isASCII := strings.IndexFunc(pw, notASCII) == -1
	length := metric.CalculateLength(pw)
	isInRange := 0 < length && length <= 100
	return isASCII && isInRange
}

//validateInput only returns true if the given language is valid (de, en)
func validateInputLanguage(lg string) bool {
	return lg == "en" || lg == "de"
}

//RequestHandler takes incoming requests and writes response for cors or for the model calculations length, complexity and predictability
func RequestHandler(w http.ResponseWriter, r *http.Request) {

	// send cors headers in development mode (different ports on localhost)
	if os.Getenv("APP_ENV") == "dev" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "language, password")
	}

	if r.Method == "OPTIONS" {
		// only handle cors preflight and return 200
		w.WriteHeader(http.StatusOK)
		return
	}

	start := time.Now()
	if os.Getenv("APP_ENV") == "dev" {
		log.Println("Processing request...")
	}

	// get json password from header
	jsonPassword := r.Header.Get("password")
	// get language from header
	language := r.Header.Get("language")

	var password string

	// unquote password
	errMarshal := json.Unmarshal([]byte(jsonPassword), &password)
	if errMarshal != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Could not decode input string")

	} else if !validateInput(password) || !validateInputLanguage(language) {
		// a bad request simply returns http status 400 without body
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error: Input password or language invalid")

	} else {
		// return actual response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		result := CalculateResult(password, language)

		err := json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error: Could not encode result: %s\n", err)
		}
	}

	if os.Getenv("APP_ENV") == "dev" {
		elapsed := time.Since(start)
		log.Printf("... Done (took %s)\n", elapsed)
	}
}
