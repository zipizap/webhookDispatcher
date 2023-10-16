package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zipizap/webhookDispatcher/packages/web"
)

func init_log() (err error) {
	// envvar LOGLEVEL can be one of: "panic", "fatal", "error", "warning", "info", "debug", "trace"
	var logLevel log.Level
	{
		envVarLogLevel, exists := os.LookupEnv("LOGLEVEL")
		if !exists {
			envVarLogLevel = "info"
		}
		logLevel, err = log.ParseLevel(envVarLogLevel)
		if err != nil {
			return err
		}
	}

	log.SetLevel(logLevel)
	return nil
}

func init() {
	// init_log
	{
		err := init_log()
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// Set up the HTTP server and routes
	{
		log.Info("Starting web server on :8080")
		http.HandleFunc("/web/data", web.HandlerWebData)
		http.HandleFunc("/web", web.HandlerWeb)
		http.HandleFunc("/", web.HandlerWebhookReceiver)

		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
