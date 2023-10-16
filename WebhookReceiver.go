package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func WebhookReceiver(w http.ResponseWriter, r *http.Request) {
	// Create an IncommingWebhook struct
	var err error
	var incommingWebhookRequest *IncommingWebhookRequest
	{
		incommingWebhookRequest, err = NewIncommingWebhookRequest(r)
		if err != nil {
			errorMsg := fmt.Sprintf("Error in NewIncommingWebhookRequest: %v", err)
			log.Error(errorMsg)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
	}

	var iJson string
	{
		iJson, err = incommingWebhookRequest.asJson()
		if err != nil {
			errorMsg := fmt.Sprintf("Error with json of incommingWebhookRequest: %v", err)
			log.Error(errorMsg)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
	}
	log.Infof("WebhookReceiver: Received new incommingWebhookRequest:\n%s", iJson)

	// Send a response back with the appropriate return code
	w.WriteHeader(http.StatusOK)

	// Call WebhookMatcher passing the IncommingWebhook instance by reference
	WebhookMatcher(incommingWebhookRequest)
}
