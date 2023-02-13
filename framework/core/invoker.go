package core

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// CoreRuntimeRequest - Structure for a request from core runtime to sub-runtime with event,
// context, and handler informations to dynamically import.
type CoreRuntimeRequest struct {
	Event       interface{}      `json:"event"`
	Context     ExecutionContext `json:"context"`
	HandlerName string           `json:"handlerName"`
	HandlerPath string           `json:"handlerPath"`
	TriggerType TriggerType      `json:"-"`
}

// FunctionInvoker - In charge of running sub-runtime processes, and invoke it with all the necessary informations
// to bootstrap the language-specific wrapper to run function handlers.
type FunctionInvoker struct {
	RuntimeBridge   string
	RuntimeBinary   string
	HandlerFilePath string
	HandlerName     string
	IsBinary        bool
	client          *http.Client
	upstreamURL     string
}

const (
	userAgentHeaderKey   = "User-Agent"
	contentTypeHeaderKey = "Content-Type"
)

// NewInvoker - Initialize runtime configuration to execute function handler
// runtimeBinaryPath - Absolute Path to runtime binary (e.g. /usr/bin/python3, /usr/bin/node)
// runtimeBridgePath - Absolute Path to runtime bridge script to start sub-runtime (e.g. /home/app/index.js)
// handlerFilePath - Absolute Path to function handler file (e.g. /home/app/function/myFunction.js for JavaScript or /home/app/function/myHandler for a binary file)
// handlerName - Name of the exported function to use as a Handler (Only for non-compiled languages) to dynamically import function (e.g. handler)
// upstreamURL - URL to sub-runtime HTTP server (e.g. http://localhost:8081)
// isBinaryHandler - Wether function Handler is a binary (Compiled languages).
func NewInvoker(runtimeBinaryPath, runtimeBridgePath, handlerFilePath, handlerName, upstreamURL string, isBinaryHandler bool) (*FunctionInvoker, error) {
	// Need binary path => /usr/local/bin/python3
	// Need runtime bridgle file path => /home/app/runtimes/python3/index.py
	return &FunctionInvoker{
		RuntimeBridge:   runtimeBridgePath,
		RuntimeBinary:   runtimeBinaryPath,
		HandlerFilePath: handlerFilePath,
		HandlerName:     handlerName,
		IsBinary:        isBinaryHandler,
		client:          &http.Client{},
		upstreamURL:     upstreamURL,
	}, nil
}

// Execute - a given function handler, and handle response.
func (fn *FunctionInvoker) Execute(event interface{}, context ExecutionContext, triggerType TriggerType) (*http.Request, error) {
	reqBody := CoreRuntimeRequest{
		Event:       event,
		Context:     context,
		HandlerName: fn.HandlerName,
		HandlerPath: fn.HandlerFilePath,
		TriggerType: triggerType,
	}

	return fn.streamRequest(reqBody)
}

func (fn *FunctionInvoker) streamRequest(reqBody CoreRuntimeRequest) (*http.Request, error) {
	bodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	userAgent := ""
	event := APIGatewayProxyRequest{}

	if reqBody.TriggerType == TriggerTypeHTTP {
		request, castSucceeds := reqBody.Event.(APIGatewayProxyRequest)
		if castSucceeds {
			userAgent = request.Headers[userAgentHeaderKey]
			event = request
		}
	}

	body := bytes.NewReader(bodyJSON)

	request, err := http.NewRequest(http.MethodPost, fn.upstreamURL, body)
	if err != nil {
		return nil, err
	}

	for key, values := range event.Headers {
		request.Header.Set(key, values)
	}

	for key, values := range event.MultiValueHeaders {
		for idx := range values {
			request.Header.Add(key, values[idx])
		}
	}

	if request.Header.Get(contentTypeHeaderKey) == "" {
		request.Header.Set(contentTypeHeaderKey, "application/json")
	}

	if userAgent != "" {
		request.Header.Set(userAgentHeaderKey, userAgent)
	}

	return request, nil
}
