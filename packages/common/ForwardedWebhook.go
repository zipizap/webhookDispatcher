package common

var AllForwardedWebhooks TAllForwardedWebhooks = []*ForwardedWebhook{}

type TAllForwardedWebhooks []*ForwardedWebhook

type ForwardedWebhook struct {
	IncommingWebhookRequest *IncommingWebhookRequest
	Rule                    *EvalutatedRule

	ForwardedWebhookResponses map[RuleForwardDestinationName]*ForwardedWebhookResponse
}

func NewForwardedWebhook(incommingWebhookRequest *IncommingWebhookRequest, rule *Rule) (forwardedWebhook *ForwardedWebhook) {
	erule := NewEvaluatedRule(*rule)
	forwardedWebhook = &ForwardedWebhook{
		IncommingWebhookRequest:   incommingWebhookRequest,
		Rule:                      erule,
		ForwardedWebhookResponses: make(map[RuleForwardDestinationName]*ForwardedWebhookResponse),
	}
	AllForwardedWebhooks.Append(forwardedWebhook)
	return forwardedWebhook
}

// AsJson returns json representation of ForwardedWebhook
func (o *ForwardedWebhook) AsJson() (jsonString string, err error) {
	return _asJson(o)
}

func (a *TAllForwardedWebhooks) Append(forwardedWebhook *ForwardedWebhook) {
	// #AllForwardedWebhooks = #IncommingWebhookRequest x #Rules
	AllForwardedWebhooksMaxLen := 100000

	if len(*a) == AllForwardedWebhooksMaxLen {
		*a = (*a)[1:]
	}
	*a = append(*a, forwardedWebhook)
}

func (o *TAllForwardedWebhooks) AsJson() (jsonString string, err error) {
	return _asJson(o)
}
