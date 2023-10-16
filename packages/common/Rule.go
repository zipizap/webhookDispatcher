package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jmespath/go-jmespath"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type RuleConditionName string

type RuleConditionExpression string

type RuleForwardDestinationName string
type RuleForwardDestination struct {
	Url string `yaml:"url"`
}

type Rule struct {
	RuleName            string                                                `yaml:"rule_name"`
	JmespathConditions  map[RuleConditionName]RuleConditionExpression         `yaml:"jmespath_conditions"`
	ForwardDestinations map[RuleForwardDestinationName]RuleForwardDestination `yaml:"forward_destinations"`
}

type EvalutatedRule struct {
	Rule
	EvaluationResult string // "true", "false", "error-message"
}

func loadRulesFromYamlFile(filepath string) (rules []*Rule, err error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &rules)
	if err != nil {
		return nil, err
	}

	return rules, nil
}

func NewEvaluatedRule(rule Rule) (erule *EvalutatedRule) {
	return &EvalutatedRule{
		Rule:             rule,
		EvaluationResult: "non-evaluated",
	}
}

// MatchRule will compare each erule.JmespathConditions, and return true if none of the conditions is false
//   - if there are no conditions => return true
//   - if all conditions are true => return true
//   - if one-or-more conditions are false => return false
//
// The IncommingWebhook json is shown in the logs, the jmespath is evaluated against it
func (erule *EvalutatedRule) MatchRule(incommingWebhook *IncommingWebhookRequest) (allConditionsAreTrue bool) {
	var err error
	var iAsJson string
	var iAsMapIfc map[string]interface{}
	{
		iAsJson, err = incommingWebhook.AsJson()
		if err != nil {
			log.Errorf("Unexpected error from incommingWebhook.asJson: '%s' '%s'", incommingWebhook.Timestamp, err)
			return false
		}
		err = json.Unmarshal([]byte(iAsJson), &iAsMapIfc)
		if err != nil {
			log.Errorf("Unexpected error from json.Unmarshal: %s", err)
			return false
		}
	}

	for conditionName, conditionExpression := range erule.JmespathConditions {
		var result bool
		{
			resultIfc, err := jmespath.Search(string(conditionExpression), iAsMapIfc)
			if err != nil {
				errMsg := fmt.Sprintf("Error matching condition with IncommingWebhook:\n  Error: '%s'\n  IncommingWebhook.Timestamp: '%s'\n  RuleName: '%s'\n  ConditionName: '%s'\n  ConditionExpression: '%s'\n  IncommingWebhook.JsonData: \n%s\n", err, incommingWebhook.Timestamp, erule.RuleName, conditionName, conditionExpression, iAsJson)
				erule.EvaluationResult = errMsg
				log.Errorf(errMsg)
				return false
			}
			var ok bool

			// Assure result is boolean
			result, ok = resultIfc.(bool)
			if !ok {
				errMsg := fmt.Sprintf("Error typecasting result to bool\n  IncommingWebhookRequest:  '%s'\n RuleName: '%s'\n  ConditionName: '%s'\n  ConditionExpression: '%s'\n  ConditionResult(Ifc): '%v'\n  IncommingWebhook: \n%s\n", incommingWebhook.Timestamp, erule.RuleName, conditionName, conditionExpression, resultIfc, iAsJson)
				erule.EvaluationResult = errMsg
				log.Errorf(errMsg)
				return false
			}
			erule.EvaluationResult = fmt.Sprintf("%v", resultIfc)
		}
		log.Infof("Evaluated IncommingWebhook '%s' with Rule/Condition\t[%v]\t'%s/%s' '%v'", incommingWebhook.Timestamp, erule.EvaluationResult, erule.RuleName, conditionName, conditionExpression)
		if !result {
			return false
		}
	}
	return true
}
