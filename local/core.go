package local

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/scaleway/serverless-functions-go/framework/core"
	"github.com/scaleway/serverless-functions-go/framework/function"
)

const (
	// payloadSizeWarn if the body length in bytes is higher than this value it will trigger warning.
	// In production this will stop the request and return 500.
	payloadSizeWarn = 6291456
)

// CoreProcessing processes the main core
func CoreProcessing(httpResp http.ResponseWriter, httpReq *http.Request, handler function.ScwFuncV1) {
	if core.IsRejectedRequest(httpReq) {
		log.Default().Println("request will be rejected for calling favico or robots.txt")
	}

	if httpReq.ContentLength > payloadSizeWarn {
		log.Default().Println("request can be rejected because it's too big")
	}

	bodyBytes, err := io.ReadAll(httpReq.Body)
	if err != nil {
		panic(err)
	}

	formattedRequest := core.FormatEventHTTP(httpReq, bodyBytes)

	invoker := core.FunctionInvoker{}

	reqForFaaS, err := invoker.Execute(formattedRequest, core.GetExecutionContext())
	if err != nil {
		panic(err)
	}

	reqForFaaS.Host = httpReq.Host
	_ = SubProcessing(httpResp, reqForFaaS)

	InjectIngressHeaders(reqForFaaS)

	writerRecorder := httptest.NewRecorder()
	handler(writerRecorder, reqForFaaS)

	// Body is closed but linter reports it.
	//nolint:bodyclose
	recorderResp := writerRecorder.Result()
	defer func() {
		if recorderResp.Body != nil {
			recorderResp.Body.Close()
		}
	}()

	coreResp, err := core.GetResponse(recorderResp)
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

	core.SetHeaders(reqForFaaS.Header, httpResp.Header())

	httpResp.Header().Set("Access-Control-Allow-Origin", "*")
	httpResp.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	InjectEgressHeaders(httpResp)

	core.SetHeaders(recorderResp.Header, httpResp.Header())

	core.HydrateHTTPResponse(httpResp, responseBody, coreResp.StatusCode)
}
