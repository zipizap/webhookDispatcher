package common

import (
	"encoding/json"
	"io"
	"net/http"
)

// ForwardedWebhookResponse is very similar to http.Response with following peculiarities:
// - Body is set to []byte  (overrides/"hides" Response.Body in Response.Body, same as IncommingWebhook)
// - multiple fields unused: TLS, Close, Uncompressed,...
type ForwardedWebhookResponse struct {
	http.Response

	Body      interface{}
	bodyBytes []byte

	// These unused fields are overriden so struct can be represented as json
	GetBody bool
	TLS     bool
}

func NewForwardedWebhookResponse(fwdResp *http.Response) (forwardedWebhookResponse *ForwardedWebhookResponse, err error) {
	// Read the response body
	var fwdRespBodyIfc interface{}
	var fwdResBodyBytes []byte
	{
		fwdResBodyBytes, err = io.ReadAll(fwdResp.Body)
		if err != nil {
			return nil, err
		}
		if len(fwdResBodyBytes) > 0 {
			err = json.Unmarshal(fwdResBodyBytes, &fwdRespBodyIfc)
			if err != nil {
				fwdRespBodyIfc = map[string]string{"bodyraw": string(fwdResBodyBytes)}
				// return nil, err
			}
		}
	}
	forwardedWebhookResponse = &ForwardedWebhookResponse{
		Response: http.Response{
			Status:     fwdResp.Status,
			StatusCode: fwdResp.StatusCode,
			Proto:      fwdResp.Proto,
			ProtoMajor: fwdResp.ProtoMajor,
			ProtoMinor: fwdResp.ProtoMinor,
			Header:     fwdResp.Header,
			// Body:             nil,
			ContentLength:    fwdResp.ContentLength,
			TransferEncoding: fwdResp.TransferEncoding,
			// Close:            false,
			// Uncompressed:     false,
			Trailer: fwdResp.Trailer,
			// Request: fwdResp.Request,
			// TLS:              &tls.ConnectionState{},
		},
		Body:      fwdRespBodyIfc,
		bodyBytes: fwdResBodyBytes,
	}
	return forwardedWebhookResponse, nil
}

// AsJson returns json representation of ForwardedWebhookResponse
//
// Note that the json will sbow the value of ForwardedWebhookResponse.Body but not
// of ForwardedWebhookResponse.Response.Body (which gets "hidden")
func (o *ForwardedWebhookResponse) AsJson() (jsonString string, err error) {
	return _asJson(o)
}
