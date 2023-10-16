package main

func fetchRulesFromDB() []*Rule {
	// For now, just return a hardcoded list of rules
	return []*Rule{
		{
			JmespathConditions: map[RuleConditionName]RuleConditionExpression{
				"condition1-false": `contains('return false', 'true')`,
				"condition2-true":  `Method == 'POST'`,
				// "condition2-true":  `contains('return true', 'true')`,
			},
			ForwardDestination: RuleForwardDestionation{
				Url: "http://localhost:8081/post",
			},
		},
	}
}
