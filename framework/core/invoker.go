package core

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// CoreRuntimeRequest - Structure for a request from core runtime to sub-runtime with event,
// context, and handler informations to dynamically import.
type CoreRuntimeRequest struct {
	Event       APIGatewayProxyRequest `json:"event"`
	Context     ExecutionContext       `json:"context"`
	HandlerName string                 `json:"handlerName"`
	HandlerPath string                 `json:"handlerPath"`
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

// NewInvoker Initialize runtime configuration to execute function handler
func NewInvoker(
	runtimeBinaryPath,
	runtimeBridgePath,
	handlerFilePath, handlerName,
	upstreamURL string,
	isBinaryHandler bool,
) (*FunctionInvoker, error) {
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

// Execute a given function handler, and handle response.
func (fn *FunctionInvoker) Execute(event APIGatewayProxyRequest, ctx ExecutionContext) (*http.Request, error) {
	reqBody := CoreRuntimeRequest{
		Event:       event,
		Context:     ctx,
		HandlerName: fn.HandlerName,
		HandlerPath: fn.HandlerFilePath,
	}

	return fn.StreamRequest(reqBody)
}

//nolint:gocritic
func (fn *FunctionInvoker) StreamRequest(reqBody CoreRuntimeRequest) (*http.Request, error) {
	bodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	event := reqBody.Event

	body := bytes.NewReader(bodyJSON)

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fn.upstreamURL, body)
	if err != nil {
		return nil, err
	}

	request.URL.Path = event.Path

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

	if event.Headers[userAgentHeaderKey] != "" {
		request.Header.Set(userAgentHeaderKey, event.Headers[userAgentHeaderKey])
	}

	return request, nil
}
