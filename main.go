package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/zipizap/webhookDispatcher/packages/common"
)

func main() {
	// Initialize the database connection
	// db := initDB()

	// Set up the HTTP server
	log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
	log.Info("Starting web server on :8080")
	http.HandleFunc("/", common.WebhookReceiver)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
