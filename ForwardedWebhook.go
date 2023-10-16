package main

type ForwardedWebhook struct {
	IncommingWebhookRequest *IncommingWebhookRequest
	Rule                    *Rule

	ForwardedWebhookResponse ForwardedWebhookResponse
}

// asJson returns json representation of ForwardedWebhook
func (o *ForwardedWebhook) asJson() (jsonString string, err error) {
	return _asJson(o)
}
