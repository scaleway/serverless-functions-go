package local

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type subRuntimeRequest struct {
	Event struct {
		HTTPMethod            string            `json:"httpMethod"`
		Headers               map[string]string `json:"headers"`
		QueryStringParameters map[string]string `json:"queryStringParameters"`
		Body                  string            `json:"body"`
	} `json:"event"`
}

// SubProcessing simulates speicifc go workflow that can happens in the FaaS environment.
// The request of the body must complies with subRuntimeRequest type to be processed.
func SubProcessing(httpResp http.ResponseWriter, httpReq *http.Request) error {
	bodyBytes, err := io.ReadAll(httpReq.Body)
	if err != nil {
		httpResp.WriteHeader(http.StatusInternalServerError)
		_, _ = httpResp.Write([]byte(err.Error()))

		return err
	}

	var req subRuntimeRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		httpResp.WriteHeader(http.StatusInternalServerError)
		_, _ = httpResp.Write([]byte("Cannot unmarshal event from core runtime"))

		return err
	}

	httpReq.Method = req.Event.HTTPMethod

	httpReq.Header = make(map[string][]string, len(req.Event.Headers))
	for key, value := range req.Event.Headers {
		httpReq.Header[key] = []string{value}
	}

	params := httpReq.URL.Query()
	for key, value := range req.Event.QueryStringParameters {
		params.Set(key, value)
	}

	httpReq.URL.RawQuery = params.Encode()

	httpReq.Body = io.NopCloser(strings.NewReader(req.Event.Body))

	return nil
}
