package common

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func WebhookLauncher(forwardedWebhook *ForwardedWebhook) {
	log.Infof("Launching ForwardedWebhook(s) for IncommingWebhookRequest+Rule '%s'", forwardedWebhook.Rule.RuleName)
	err := forwardRequest(forwardedWebhook)
	if err != nil {
		log.Errorf("Error launching ForwardedWebhook(s): '%s'", err)
	}
	fAsJson, err := forwardedWebhook.asJson()
	if err != nil {
		log.Errorf("Error from forwardedWebhook.asJson(): '%s'", err)
	} else {
		log.Debugf("Launch complete of ForwardedWebhook(s):\n%s\n", fAsJson)
	}
}

// forwardRequest creates a new POST request to the given URL with the same headers and body as the original request
func forwardRequest(forwardedWebhook *ForwardedWebhook) (err error) {
	incReq := forwardedWebhook.IncommingWebhookRequest
	rule := forwardedWebhook.Rule
	fwdMethod := incReq.Method
	for fwdDestinationName, fwdDestination := range rule.ForwardDestinations {
		log.Infof("Launching ForwardedWebhook for rule '%s' forward_destination '%s'", rule.RuleName, fwdDestinationName)
		fwdUrl := fwdDestination.Url
		fwdBodyBytes := incReq.bodyBytes
		fwdHeaders := incReq.Header

		// Create a new request
		var fwdReq *http.Request
		{
			fwdReq, err = http.NewRequest(fwdMethod, fwdUrl, bytes.NewBuffer(fwdBodyBytes))
			if err != nil {
				log.Errorf("Unexpected error from http.NewRequest: %s", err)
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
				log.Errorf("Unexpected error from client.Do: %s", err)
				continue
			}
			defer fwdResp.Body.Close()
		}

		// Append to forwardedWebhook.ForwardedWebhookResponses[fwdDestinationName]
		{
			forwardedWebhookResponse, err := NewForwardedWebhookResponse(fwdResp)
			if err != nil {
				log.Errorf("Unexpected error from NewForwardedWebhookResponse: %s", err)
			}
			if forwardedWebhook.ForwardedWebhookResponses == nil {
				forwardedWebhook.ForwardedWebhookResponses = make(map[RuleForwardDestinationName]*ForwardedWebhookResponse)
			}
			forwardedWebhook.ForwardedWebhookResponses[fwdDestinationName] = forwardedWebhookResponse
		}
	}
	return nil
}
