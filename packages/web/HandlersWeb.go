package web

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/zipizap/webhookDispatcher/packages/common"
)

func HandlerWeb(w http.ResponseWriter, r *http.Request) {
	log.Infof("HandlerWeb - got web request %s %s", r.Method, r.URL.Path)

	// If  LOGLEVEL=="debug" and localfile exists, then use localfile to read indexBytes
	// Else use file embedded in golang-binary
	localFile := "packages/web/website/index.html"
	var indexBytes []byte
	{
		var err error
		var logLevelIsDebug bool
		{
			if log.GetLevel() == log.DebugLevel {
				logLevelIsDebug = true
			} else {
				logLevelIsDebug = false
			}
		}
		var localFileExists bool
		{
			f, err := os.OpenFile(localFile, os.O_RDONLY, 0)
			f.Close()
			if err == nil {
				localFileExists = true
			} else {
				localFileExists = false
			}
		}

		if logLevelIsDebug && localFileExists {
			indexBytes, err = os.ReadFile(localFile)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			indexBytes, err = WebsiteFs.ReadFile("website/index.html")
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Write(indexBytes)
}

func HandlerWebData(w http.ResponseWriter, r *http.Request) {
	log.Debugf("HandlerWebData - got web request %s %s", r.Method, r.URL.Path)
	var allForwardedWebhooksAsJson string
	{
		var err error
		allForwardedWebhooksAsJson, err = common.AllForwardedWebhooks.AsJson()
		if err != nil {
			errorStr := fmt.Sprintf("Error unexpected from common.AllForwardedWebhooks.AsJson(): '%s'", err)
			http.Error(w, errorStr, http.StatusInternalServerError)
			log.Error(errorStr)
			return
		}
		// log.Debugf("AllForwardedWebhooks is now:\n%s\n", allForwardedWebhooksAsJson)
	}

	// Set the header to content type 'application/json'
	w.Header().Set("Content-Type", "application/json")

	// Write JSON data to response body
	{
		_, err := w.Write([]byte(allForwardedWebhooksAsJson))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
