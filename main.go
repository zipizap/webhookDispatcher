package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize the database connection
	// db := initDB()

	// Set up the HTTP server
	log.SetLevel(log.DebugLevel)
	log.Info("Starting web server on :8080")
	http.HandleFunc("/", WebhookReceiver)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
