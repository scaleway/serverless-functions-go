package framework

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/scaleway/serverless-functions-go/framework/core"
	"github.com/scaleway/serverless-functions-go/framework/function"
)

const (
	// payloadSizeWarn if the body length in bytes is higher than this value it will trigger warning.
	// In production this will stop the request and return 500.
	payloadSizeWarn = 6291456

	headerContentLen = "Content-Length"
)

// CoreProcessing processes the main core
func CoreProcessing(httpResp http.ResponseWriter, httpReq *http.Request, handler function.ScwFuncV1) {
	httpResp.Header().Set("Access-Control-Allow-Origin", "*")
	httpResp.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if isRejectedRequest(httpReq) {
		log.Default().Print("request can be rejected for calling favico or robots.txt\n")
	}

	if httpReq.ContentLength > payloadSizeWarn {
		log.Default().Print("request can be rejected because it's too big\n")
	}

	bodyBytes, err := io.ReadAll(httpReq.Body)
	if err != nil {
		panic(err)
	}

	formattedRequest := core.FormatEventHTTP(httpReq, bodyBytes)

	invoker := core.FunctionInvoker{}
	reqForFaaS, err := invoker.Execute(formattedRequest, core.GetExecutionContext(), core.TriggerTypeHTTP)
	if err != nil {
		panic(err)
	}

	SubProcessing(httpResp, reqForFaaS)

	reqId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	// headers for function
	reqForFaaS.Header.Add("Forwarded", "for="+httpReq.RemoteAddr+";proto=http")
	reqForFaaS.Header.Add("K-Proxy-Request", "activator")
	reqForFaaS.Header.Add("X-Forwarded-For", httpReq.RemoteAddr)
	reqForFaaS.Header.Add("X-Forwarded-For", "100.1.1.1")
	reqForFaaS.Header.Add("X-Forwarded-For", "100.1.1.2")
	reqForFaaS.Header.Add("X-Forwarded-Proto", "http")
	reqForFaaS.Header.Add("X-Request-Id", reqId.String())
	reqForFaaS.Header.Add("X-Envoy-External-Address", httpReq.RemoteAddr)

	writerRecorder := httptest.NewRecorder()
	handler(writerRecorder, reqForFaaS)

	coreResp, err := core.GetResponse(writerRecorder.Result())
	if err != nil {
		panic(err)
	}

	responseBody := coreResp.Body

	// If user's handler specifies the parameter isBase64Encoded, we need to transform base64 response to byte array
	if coreResp.IsBase64Encoded && len(coreResp.Body) > 0 {
		var bodyString string
		if err := json.Unmarshal(coreResp.Body, &bodyString); err != nil {
			panic(err)
		}

		base64Binary, err := base64.StdEncoding.DecodeString(bodyString)
		if err != nil {
			panic(err)
		}

		responseBody = base64Binary
	}

	setHeaders(reqForFaaS.Header, httpResp.Header())
	// Exeute generated by subrt layer

	httpResp.Header().Set("server", "envoy")

	hydrateHTTPResponse(httpResp, responseBody, coreResp.StatusCode)
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

func isRejectedRequest(request *http.Request) bool {
	return request.URL.Path == "/favicon.ico" || request.URL.Path == "/robots.txt"
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
