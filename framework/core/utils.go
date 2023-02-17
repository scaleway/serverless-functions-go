package core

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

var (
	// ErrInvalidHTTPResponseFormat Error type for mal-formatted responses from handler.
	ErrInvalidHTTPResponseFormat = errors.New("handler result for HTTP response is mal-formatted")

	// ErrRespEmpty is returned if we try to read a nil response.
	ErrRespEmpty = errors.New("http response is empty")
)

// ResponseHTTP Type for HTTP triggers response emitted by function handlers.
type ResponseHTTP struct {
	StatusCode      int                 `json:"statusCode"`
	Body            json.RawMessage     `json:"body"`
	Headers         map[string][]string `json:"headers"`
	IsBase64Encoded bool                `json:"isBase64Encoded"`
}

// IsJSON returns true if the input is a valid JSON message.
func IsJSON(bytes []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bytes, &js) == nil
}

// GetResponse Transform a response string into an HTTP Response structure.
func GetResponse(response *http.Response) (*ResponseHTTP, error) {
	if response == nil {
		return nil, ErrRespEmpty
	}

	handlerResponse := ResponseHTTP{
		Headers:    response.Header,
		StatusCode: response.StatusCode,
	}

	// Read body content
	if response.Body != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("error on body read : %s", err.Error())

			return nil, ErrInvalidHTTPResponseFormat
		}

		// Unmarshalling into handler response can modify original body so we keep it in new variable
		bbCopy := make([]byte, len(bodyBytes))
		copy(bbCopy, bodyBytes)

		handlerResponse.Body = bodyBytes

		if IsJSON(bbCopy) {
			_ = json.Unmarshal(bbCopy, &handlerResponse)

			// first we try to find headers as map[string]string
			type headerFlat struct {
				Headers map[string]string `json:"headers"`
			}

			hflat := headerFlat{}

			_ = json.Unmarshal(bbCopy, &hflat)

			for key, val := range hflat.Headers {
				handlerResponse.Headers[key] = []string{val}
			}

			// then we try to find headers as map[string][]string
			type headerSlice struct {
				Headers map[string][]string `json:"headers"`
			}

			hslice := headerSlice{}
			_ = json.Unmarshal(bbCopy, &hslice)

			// if headers are both definied as "flat" or "slice", the "slice" version will override "flat" version
			for key, val := range hslice.Headers {
				if len(val) == 0 {
					// avoid overriding key with empty value.
					continue
				}
				handlerResponse.Headers[key] = val
			}
		}
	}

	if response.StatusCode == 0 {
		handlerResponse.StatusCode = http.StatusOK
	}

	return &handlerResponse, nil
}
