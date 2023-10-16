package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func WebhookMatcher(incommingWebhookRequest *IncommingWebhookRequest) {
	// Fetch the rules
	var rules []*Rule
	{
		rulesYamlFile := "./rules.yaml"
		var err error
		rules, err = loadRulesFromYamlFile(rulesYamlFile)
		if err != nil {
			log.Panicf("Error unexpected from LoadRulesFromYamlFile, '%s'", err)
			return
		}
	}

	// Iterate over the rules
	log.Infof("Received  IncommingWebhook '%s' - going to verify if any rule matches", incommingWebhookRequest.Timestamp)

	for _, rule := range rules {
		// Create a ForwardedWebhook instance
		forwardedWebhook := NewForwardedWebhook(incommingWebhookRequest, rule)
		erule := forwardedWebhook.Rule

		// When  incommingWebhook matches rule, call WebhookLauncher(...)
		if erule.MatchRule(incommingWebhookRequest) {
			log.Infof("IncommingWebhookRequest '%s' matched rule '%s'", incommingWebhookRequest.Timestamp, erule.RuleName)
			// Call WebhookLauncher
			//go WebhookLauncher(forwardedWebhook)
			WebhookLauncher(forwardedWebhook)
		}

		// Debug allForwardedWebhooksAsJson
		{
			allForwardedWebhooksAsJson, err := AllForwardedWebhooks.AsJson()
			if err != nil {
				log.Errorf("Error unexpected from AllForwardedWebhooks.AsJson(), '%s'", err)
				return
			}
			log.Debugf("AllForwardedWebhooks is now:\n%s\n", allForwardedWebhooksAsJson)
			// tmp, remove me
			_ = os.WriteFile("allForwardedWebhooksAsJson.json", []byte(allForwardedWebhooksAsJson), 0644)
		}
	}
}
