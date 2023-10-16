package common

import (
	"encoding/json"
	"io"
	"net/http"
)

// IncommingWebhookRequest is very similar to http.Request with following peculiarities:
//   - Body is "overriden" to []byte
//     The overriden is done by hidding the embedded-struct Request.Body with the base Body field
//     , wich effectively overrides the Request.Body in json and jmespath (see https://goplay.tools/snippet/Qso-gzOZNKz)
//   - multiple fields unused: GetBody, Close,  TLS, Cancel, Response, ctx, ...
type IncommingWebhookRequest struct {
	http.Request

	Body      interface{}
	bodyBytes []byte

	// These unused fields are overriden so struct can be represented as json
	GetBody bool
	Cancel  bool
	TLS     bool
}

func NewIncommingWebhookRequest(r *http.Request) (incommingWebhook *IncommingWebhookRequest, err error) {
	// Read the request body
	var bodyBytes []byte
	var bodyIfc interface{}
	{
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bodyBytes, &bodyIfc)
		if err != nil {
			return nil, err
		}
	}

	incommingWebhook = &IncommingWebhookRequest{
		Request: http.Request{
			Method:           r.Method,
			URL:              r.URL,
			Proto:            r.Proto,
			ProtoMajor:       r.ProtoMajor,
			ProtoMinor:       r.ProtoMinor,
			Header:           r.Header,
			ContentLength:    r.ContentLength,
			TransferEncoding: r.TransferEncoding,
			Host:             r.Host,
			Form:             r.Form,
			PostForm:         r.PostForm,
			MultipartForm:    r.MultipartForm,
			Trailer:          r.Trailer,
			RemoteAddr:       r.RemoteAddr,
			RequestURI:       r.RequestURI,
		},
		Body:      bodyIfc,
		bodyBytes: bodyBytes,
	}
	return incommingWebhook, nil
}

// asJson returns json representation of IncommingWebhook
//
// Note that the json will sbow the value of IncommingWebhook.Body but not
// of IncommingWebhook.Request.Body (which gets "hidden")
func (o *IncommingWebhookRequest) asJson() (jsonString string, err error) {
	return _asJson(o)
}

// _asJson returns jsonString representation from ifc
func _asJson(ifc interface{}) (jsonString string, err error) {
	jsonBytes, err := json.MarshalIndent(ifc, "", "    ")
	if err != nil {
		return "", err
	}
	jsonString = string(jsonBytes)
	return jsonString, nil
}
