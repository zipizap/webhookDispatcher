package common

type ForwardedWebhook struct {
	IncommingWebhookRequest *IncommingWebhookRequest
	Rule                    *Rule

	ForwardedWebhookResponses map[RuleForwardDestinationName]*ForwardedWebhookResponse
}

// asJson returns json representation of ForwardedWebhook
func (o *ForwardedWebhook) asJson() (jsonString string, err error) {
	return _asJson(o)
}
