package core

import (
	"errors"
	"io"
	"net/http"
)

var (
	// ErrNotSupportedTrigger Error when event is assigned to not supported trigger types.
	ErrNotSupportedTrigger = errors.New("trigger Type is not supported by Scaleway Functions Runtime")

	// ErrReadBody returned if body read returns error
	ErrReadBody = errors.New("unable to read request body")
)

// FormatEvent Format event
func FormatEvent(req *http.Request) (APIGatewayProxyRequest, error) {
	// request body is the event
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return APIGatewayProxyRequest{}, ErrReadBody
	}

	return FormatEventHTTP(req, bodyBytes), nil
}
