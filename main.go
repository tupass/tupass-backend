package main

import (
	"os"

	"github.com/tupass/tupass-backend/api"
	"github.com/tupass/tupass-backend/web"
)

// main starts the in API included webserver.
func main() {

	// load passwordList from file to heap for predictability calculation
	api.SetupPasswordList()

	// listen on port 8000 for staging/development
	serverPort := "8000"
	if os.Getenv("APP_ENV") == "prod" {
		// listen on port 8001 for production
		serverPort = "8001"
	}
	web.StartServer(serverPort)
}
