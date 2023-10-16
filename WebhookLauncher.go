package main

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func WebhookLauncher(forwardedWebhook *ForwardedWebhook) {
	log.Info("Launching ForwardedWebhook from IncommingWebhookRequest+Rule")
	err := forwardRequest(forwardedWebhook)
	if err != nil {
		log.Errorf("Error launching ForwardedWebhook: '%s'", err)
	}
	fAsJson, err := forwardedWebhook.asJson()
	if err != nil {
		log.Errorf("Error from forwardedWebhook.asJson(): '%s'", err)
	} else {
		log.Debugf("Launched ForwardedWebhook\n%s\n", fAsJson)
	}
}

// forwardRequest creates a new POST request to the given URL with the same headers and body as the original request
// func forwardRequest(headers http.Header, body []byte, url string) error {
func forwardRequest(forwardedWebhook *ForwardedWebhook) (err error) {
	incReq := forwardedWebhook.IncommingWebhookRequest
	rule := forwardedWebhook.Rule
	fwdMethod := incReq.Method
	fwdUrl := rule.ForwardDestination.Url
	fwdBodyBytes := incReq.bodyBytes
	fwdHeaders := incReq.Header

	// Create a new request
	var fwdReq *http.Request
	{
		fwdReq, err = http.NewRequest(fwdMethod, fwdUrl, bytes.NewBuffer(fwdBodyBytes))
		if err != nil {
			return err
		}
	}

	// Copy the headers from the original request
	{
		for name, values := range fwdHeaders {
			for _, value := range values {
				fwdReq.Header.Add(name, value)
			}
		}
	}

	// Send the request, define fwdResp
	var fwdResp *http.Response
	{
		client := &http.Client{}
		fwdResp, err = client.Do(fwdReq)
		if err != nil {
			return err
		}
		defer fwdResp.Body.Close()
	}

	// Set forwardedWebhook.ForwardedWebhookResponse
	{
		forwardedWebhookResponse, err := NewForwardedWebhookResponse(fwdResp)
		if err != nil {
			return err
		}
		forwardedWebhook.ForwardedWebhookResponse = *forwardedWebhookResponse
	}

	return nil
}
