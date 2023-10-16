package web

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/zipizap/webhookDispatcher/packages/common"
)

func HandlerWebhookReceiver(w http.ResponseWriter, r *http.Request) {
	// ignore (return) if GET /favicon.ico
	{
		if r.Method == "GET" && r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Create an IncommingWebhook struct
	var err error
	var incommingWebhookRequest *common.IncommingWebhookRequest
	{
		incommingWebhookRequest, err = common.NewIncommingWebhookRequest(r)
		if err != nil {
			errorMsg := fmt.Sprintf("Error in NewIncommingWebhookRequest: %v", err)
			log.Error(errorMsg)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
	}

	var iJson string
	{
		iJson, err = incommingWebhookRequest.AsJson()
		if err != nil {
			errorMsg := fmt.Sprintf("Error with json of incommingWebhookRequest: %v", err)
			log.Error(errorMsg)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
	}
	log.Infof("WebhookReceiver: Received new incommingWebhookRequest '%s'", incommingWebhookRequest.Timestamp)
	log.Debug(iJson)

	// Send a response back with the appropriate return code
	w.WriteHeader(http.StatusOK)

	// Call WebhookMatcher passing the IncommingWebhook instance by reference
	go common.WebhookMatcher(incommingWebhookRequest)
}
