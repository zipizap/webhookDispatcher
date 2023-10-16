package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// IncommingWebhookRequest is very similar to http.Request with following peculiarities:
//   - Body is "overriden" to []byte
//     The overriden is done by hidding the embedded-struct Request.Body with the base Body field
//     , wich effectively overrides the Request.Body in json and jmespath (see https://goplay.tools/snippet/Qso-gzOZNKz)
//   - multiple fields unused: GetBody, Close,  TLS, Cancel, Response, ctx, ...
type IncommingWebhookRequest struct {
	Timestamp string
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
		if len(bodyBytes) > 0 {
			err = json.Unmarshal(bodyBytes, &bodyIfc)
			if err != nil {
				return nil, err
			}
		}
	}

	incommingWebhook = &IncommingWebhookRequest{
		Timestamp: _timestampNow(),
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

// AsJson returns json representation of IncommingWebhook
//
// Note that the json will sbow the value of IncommingWebhook.Body but not
// of IncommingWebhook.Request.Body (which gets "hidden")
func (o *IncommingWebhookRequest) AsJson() (jsonString string, err error) {
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

// Ex: "20231018-140117.565"
func _timestampNow() (ts string) {
	t := time.Now()
	ts = fmt.Sprintf("%d%02d%02d-%02d%02d%02d.%03d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		t.Nanosecond()/1000000) // convert nanoseconds to milliseconds
	return ts
}
