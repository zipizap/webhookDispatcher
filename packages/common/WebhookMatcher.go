package common

import (
	"github.com/jmespath/go-jmespath"
	log "github.com/sirupsen/logrus"
)

func WebhookMatcher(incommingWebhookRequest *IncommingWebhookRequest) {
	// Fetch the rules
	// var rules [](*Rule)
	rulesYamlFile := "./rules.yaml"
	rules, err := LoadRulesFromYamlFile(rulesYamlFile)
	if err != nil {
		log.Errorf("Error unexpected from LoadRulesFromYamlFile, '%s'", err)
		return
	}

	// Iterate over the rules
	log.Infof("Going to verify if any rule matches against the IncommingWebhookRequest contents")
	for _, rule := range rules {
		// If the incommingWebhook matches the rule
		if matchRule(incommingWebhookRequest, rule) {
			log.Infof("IncommingWebhookRequest matched rule '%s'", rule.RuleName)
			// Create a ForwardedWebhook instance
			forwardedWebhook := &ForwardedWebhook{
				IncommingWebhookRequest: incommingWebhookRequest,
				Rule:                    rule,
			}

			// Call WebhookLauncher
			//go WebhookLauncher(forwardedWebhook)
			WebhookLauncher(forwardedWebhook)
		}
	}
}

// matchRule will compare each rule.JmespathConditions, and return true if one of the conditions is true
// The IncommingWebhook json is shown in the logs, the jmespath is evaluated against it
func matchRule(incommingWebhook *IncommingWebhookRequest, rule *Rule) bool {
	var err error
	var iAsJson string
	{
		iAsJson, err = incommingWebhook.asJson()
		if err != nil {
			log.Errorf("Unexpected error from incommingWebhook,asJson: '%s'", err)
			return false
		}
	}

	for conditionName, conditionExpression := range rule.JmespathConditions {
		var result bool
		{
			resultIfc, err := jmespath.Search(string(conditionExpression), incommingWebhook)
			if err != nil {
				log.Errorf("Error matching condition with IncommingWebhook\n  RuleName: '%s'\n  ConditionName: '%s'\n  ConditionExpression: '%s'\n  IncommingWebhook.JsonData: \n%s\n", rule.RuleName, conditionName, conditionExpression, iAsJson)
				return false
			}
			var ok bool
			result, ok = resultIfc.(bool)
			if !ok {
				log.Errorf("Error typecasting resultIfc to result (bool)\n  RuleName: '%s'\n  ConditionName: '%s'\n  ConditionExpression: '%s'\n  ConditionResult(Ifc): '%v'\n  IncommingWebhook: \n%s\n", rule.RuleName, conditionName, conditionExpression, resultIfc, iAsJson)
				return false
			}
		}
		log.Infof("Evaluated IncommingWebhook with Rule/Condition\t[%t]\t'%s/%s' '%v'", result, rule.RuleName, conditionName, conditionExpression)
		if !result {
			return false
		}
	}
	return true
}
