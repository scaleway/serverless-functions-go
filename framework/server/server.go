package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gitlab.infra.online.net/paas/faas-runtimes/core/events"
	"gitlab.infra.online.net/paas/faas-runtimes/core/handler"
)

const (
	defaultPort         = 8080
	defaultUpstreamHost = "http://127.0.0.1"
	defaultUpstreamPort = 8081
	headerTriggerType   = "SCW_TRIGGER_TYPE"
	payloadMaxSize      = 6291456
	headerContentLen    = "Content-Length"
)

// config keys.
const (
	handlerNameKey   = "SCW_HANDLER_NAME"
	handlerPathKey   = "SCW_HANDLER_PATH"
	runtimeBinaryKey = "SCW_RUNTIME_BINARY"
	runtimeBridgeKey = "SCW_RUNTIME_BRIDGE"
	isBinaryKey      = "SCW_HANDLER_IS_BINARY"
	upstreamPortKey  = "SCW_UPSTREAM_PORT"
	upstreamHostKey  = "SCW_UPSTREAM_HOST"
)

const (
	timeoutDuration = 5 * time.Minute
	maxHeaderSize   = 1 << 20
)

var upstreamURL = defaultUpstreamHost + ":" + strconv.Itoa(defaultUpstreamPort)

// Configure function Invoker from environment variables.
func setUpFunctionInvoker() (*handler.FunctionInvoker, error) {
	// Exported function to execute
	handlerName := os.Getenv(handlerNameKey)
	// Absolute path to handler file
	handlerPath := os.Getenv(handlerPathKey)
	// Absolute path to runtime binary (e.g. python, node)
	runtimeBinary := os.Getenv(runtimeBinaryKey)
	// Absolute path to sub-runtime file
	runtimeBridgeFile := os.Getenv(runtimeBridgeKey)
	// Whether handler is binary or not (mostly used for compiled languages)
	isBinaryHandler := os.Getenv(isBinaryKey)
	// Host/Port for sub-runtime HTTP server
	upstreamPort := os.Getenv(upstreamPortKey)
	upstreamHost := os.Getenv(upstreamHostKey)

	// Configure connection to upstream server (Function runtime's server)
	if upstreamHost == "" {
		upstreamHost = defaultUpstreamHost
	}

	if upstreamPort == "" {
		upstreamPort = strconv.Itoa(defaultUpstreamPort)
	}
	upstreamURL = fmt.Sprintf("%s:%s", upstreamHost, upstreamPort)

	fnInvoker, err := handler.NewInvoker(runtimeBinary, runtimeBridgeFile, handlerPath, handlerName, upstreamURL, isBinaryHandler == "true")
	if err != nil {
		return nil, fmt.Errorf("new invoker error %w", err)
	}

	return fnInvoker, nil
}

func waitUntilUpstreamIsReachable() {
	upstreamHostPort := strings.Replace(upstreamURL, "http://", "", 1)
	dialTimeout := time.Second

	timeout := time.After(timeoutDuration)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-timeout:
			log.Fatal().Msgf("Timeout reached (%s), function was never ready to serve requests", timeoutDuration)
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", upstreamHostPort, dialTimeout)
			if err != nil {
				log.Debug().Msg("Function is not ready to serve requests yet...")

				break
			}

			if conn != nil {
				log.Debug().Msg("Function is now ready to serve requests!")
				ticker.Stop()
				conn.Close()

				return
			}
		}
	}
}

// Start takes the function Handler, at the moment only supporting HTTP Triggers (Api Gateway Proxy events)
// It takes care of wrapping the handler with an HTTP server, which receives requests when functions are triggered
// And execute the handler after formatting the HTTP CoreRuntimeRequest to an API Gateway Proxy Event.
func Start() error {
	portEnv := os.Getenv("PORT")
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		port = defaultPort
	}

	requestHandler, err := buildRequestHandler()
	if err != nil {
		return fmt.Errorf("build request error %w", err)
	}

	s := &http.Server{
		ReadHeaderTimeout: time.Minute,
		Addr:              fmt.Sprintf(":%d", port),
		MaxHeaderBytes:    maxHeaderSize, // Max header of 1MB
		Handler:           http.HandlerFunc(requestHandler),
		// NOTE: we should either set timeouts or make explicit we don't need them
		// see https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/
	}

	// core runtime will not listen until upstream is reachable
	waitUntilUpstreamIsReachable()

	return s.ListenAndServe()
}

func buildRequestHandler() (func(http.ResponseWriter, *http.Request), error) {
	fnInvoker, err := setUpFunctionInvoker()
	if err != nil {
		return nil, err
	}

	// Start function server
	if err := fnInvoker.Start(); err != nil {
		return nil, fmt.Errorf("invorker start error %w", err)
	}

	allowFavicon, _ := strconv.ParseBool(os.Getenv("SCW_ALLOW_FAVICON"))
	return buildHandler(fnInvoker, allowFavicon), nil
}

// hydrateHttpResponse will try to fill the response writer with content of body. Addtionaly it adds
// Content-Length header, this has to by done always before any call to Write on the body.
func hydrateHTTPResponse(resp http.ResponseWriter, body json.RawMessage, statusCode int) {
	// when lambda returns a string as body it expects to return it without json encoding
	var bodyString string
	err := json.Unmarshal(body, &bodyString)
	if err == nil {
		resp.Header().Set(headerContentLen, strconv.Itoa(len(bodyString)))
		resp.WriteHeader(statusCode)

		_, _ = resp.Write([]byte(bodyString))

		return
	}

	resp.Header().Set(headerContentLen, strconv.Itoa(len(body)))
	resp.WriteHeader(statusCode)

	_, _ = resp.Write(body)
}

func buildHandler(functionInvoker *handler.FunctionInvoker, allowFavicon bool) func(_ http.ResponseWriter, _ *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		// Allow CORS
		response.Header().Set("Access-Control-Allow-Origin", "*")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// drop unwanted requests (ex: favicon.ico and robots.txt)
		if dropRequest(request, allowFavicon) {
			response.WriteHeader(http.StatusNotFound)
			return
		}

		// 2: check payload size
		defaultPayloadMaxSizeEnv := os.Getenv("SCW_PAYLOAD_MAX_SIZE")
		defaultPayloadMaxSize, err := strconv.ParseInt(defaultPayloadMaxSizeEnv, 10, 64)
		if err != nil {
			defaultPayloadMaxSize = int64(payloadMaxSize)
		}

		// Access log
		triggerLog := "Function Triggered"

		if request.ContentLength > defaultPayloadMaxSize {
			log.Debug().Msg(triggerLog)
			log.Error().Err(ErrPayloadTooLarge)
			http.Error(response, ErrPayloadTooLarge.Error(), http.StatusRequestEntityTooLarge)
			return
		}

		// Check event publisher
		triggerType, err := events.GetTriggerType(request.Header.Get(headerTriggerType))
		if err != nil {
			log.Debug().Msg(triggerLog)
			log.Err(err)
			http.Error(response, err.Error(), http.StatusBadRequest)
			return
		}

		// Format event and context
		event, err := events.FormatEvent(request, triggerType)
		if err != nil {
			log.Err(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		context := events.GetExecutionContext()

		triggerLog += ": " + request.URL.Path
		log.Debug().Msg(triggerLog)

		// Execute Handler Based on runtime
		handlerResponse, err := functionInvoker.Execute(event, context, triggerType)
		if err != nil {
			log.Err(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		defer handlerResponse.Body.Close()

		// Do not try to format HTTP response if trigger is NOT of type HTTP (would be pointless as nobody is waiting for the response)
		if triggerType != events.TriggerTypeHTTP {
			_, _ = io.WriteString(response, "executed properly") // for a trigger 201 Created might be better, so we default to 200
			return
		}

		// Get statusCode, response body, and headers
		handlerRes, err := handler.GetResponse(handlerResponse)
		if err != nil {
			log.Err(err)
			http.Error(response, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send HTTP response with Handler
		// Set Headers
		setHeaders(handlerRes.Headers, response.Header())

		responseBody := handlerRes.Body

		// If user's handler specifies the parameter isBase64Encoded, we need to transform base64 response to byte array
		if handlerRes.IsBase64Encoded && len(handlerRes.Body) > 0 {
			var bodyString string
			if err := json.Unmarshal(handlerRes.Body, &bodyString); err != nil {
				log.Err(err)
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}

			base64Binary, err := base64.StdEncoding.DecodeString(bodyString)
			if err != nil {
				log.Err(err)
				http.Error(response, err.Error(), http.StatusInternalServerError)
				return
			}

			responseBody = base64Binary
		}

		// 8 : Write HTTP response with content + add necessary headers.
		hydrateHTTPResponse(response, responseBody, handlerRes.StatusCode)
	}
}

func dropRequest(request *http.Request, allowFavicon bool) bool {
	return (!allowFavicon && request.URL.Path == "/favicon.ico") || request.URL.Path == "/robots.txt"
}

func setHeaders(input, output http.Header) {
	const (
		originCORS  = "access-control-allow-origin"
		headersCORS = "access-control-allow-headers"
	)
	// first loop to reset CORS values if necessary
	// This prevent this kind of header errors : access-control-allow-origin: *,*
	for key := range input {
		lowerKey := strings.ToLower(key)

		if lowerKey == originCORS {
			output.Del(originCORS)
		} else if lowerKey == headersCORS {
			output.Del(headersCORS)
		}
	}

	for key, values := range input {
		for idx := range values {
			output.Add(key, values[idx])
		}
	}
}
