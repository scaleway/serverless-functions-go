package core

import (
	"errors"
	"io"
	"net/http"
)

var (
	// TriggerTypeHTTP - Event trigger of type HTTP.
	TriggerTypeHTTP TriggerType = "http"
	// ValidTriggerTypes - List of supported trigger types.
	ValidTriggerTypes = []TriggerType{TriggerTypeHTTP}
	// ErrNotSupportedTrigger - Error when event is assigned to not supported trigger types.
	ErrNotSupportedTrigger = errors.New("trigger Type is not supported by Scaleway Functions Runtime")
	// ErrReadBody - returned if body read returns error
	ErrReadBody = errors.New("unable to read request body")
)

// TriggerType - Enumeration of valid trigger types supported by runtime.
type TriggerType string

// GetTriggerType - check that a given trigger type is supported by runtime.
func GetTriggerType(triggerType string) (TriggerType, error) {
	if triggerType == "" {
		return TriggerTypeHTTP, nil
	}

	for _, validType := range ValidTriggerTypes {
		if string(validType) == triggerType {
			return validType, nil
		}
	}

	return "", ErrNotSupportedTrigger
}

// FormatEvent - Format event according to given trigger type, if trigger type if not HTTP, then we assume that event
// has already been formatted by event-source.
func FormatEvent(req *http.Request, triggerType TriggerType) (interface{}, error) {
	// request body is the event
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, ErrReadBody
	}

	if triggerType == TriggerTypeHTTP {
		return FormatEventHTTP(req, bodyBytes), nil
	}

	return string(bodyBytes), nil
}
