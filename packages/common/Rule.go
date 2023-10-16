package common

import (
	"os"

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

func LoadRulesFromYamlFile(filepath string) (rules []*Rule, err error) {
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
