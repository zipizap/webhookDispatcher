package main

type RuleConditionName string

// RuleConditionExpression must return a bool-value from jmespath
// Ex: This always returns true: `contains('return true', 'true')`
type RuleConditionExpression string

type RuleForwardDestionation struct {
	Url string
}

type Rule struct {
	// JmespathConditions - one of the conditions should match (ie, conditions are OR'ed together)
	JmespathConditions map[RuleConditionName]RuleConditionExpression
	ForwardDestination RuleForwardDestionation
}
