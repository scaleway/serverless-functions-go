package core

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const headerContentLen = "Content-Length"

// APIGatewayProxyRequest contains data coming from the API Gateway proxy.
type APIGatewayProxyRequest struct {
	Resource                        string                        `json:"resource"` // The resource path defined in API Gateway
	Path                            string                        `json:"path"`     // The url path for the caller
	HTTPMethod                      string                        `json:"httpMethod"`
	Headers                         map[string]string             `json:"headers"`
	MultiValueHeaders               map[string][]string           `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string             `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string           `json:"multiValueQueryStringParameters"`
	PathParameters                  map[string]string             `json:"pathParameters"`
	StageVariables                  map[string]string             `json:"stageVariables"`
	RequestContext                  APIGatewayProxyRequestContext `json:"requestContext"`
	Body                            string                        `json:"body"`
	IsBase64Encoded                 bool                          `json:"isBase64Encoded,omitempty"`
}

// APIGatewayProxyRequestContext contains the information to identify the AWS account and resources invoking the
// Lambda function. It also includes Cognito identity information for the caller.
type APIGatewayProxyRequestContext struct {
	AccountID    string                 `json:"accountId"`
	ResourceID   string                 `json:"resourceId"`
	Stage        string                 `json:"stage"`
	RequestID    string                 `json:"requestId"`
	ResourcePath string                 `json:"resourcePath"`
	Authorizer   map[string]interface{} `json:"authorizer"`
	HTTPMethod   string                 `json:"httpMethod"`
	APIID        string                 `json:"apiId"` // The API Gateway rest API Id
}

// FormatEventHTTP converts a http.Request to internal APIGatewayProxyRequest object.
func FormatEventHTTP(req *http.Request, bodyBytes []byte) APIGatewayProxyRequest {
	queryParameters := map[string]string{}
	for key, value := range req.URL.Query() {
		queryParameters[key] = value[len(value)-1]
	}

	input := string(bodyBytes)
	isBase64Encoded := true

	_, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		isBase64Encoded = false
	}

	flatHeader := make(map[string]string, len(req.Header))

	for key, val := range req.Header {
		if len(val) > 0 {
			flatHeader[key] = strings.Join(val, ",")
		}
	}

	return APIGatewayProxyRequest{
		Path:                  req.URL.Path,
		HTTPMethod:            req.Method,
		Headers:               flatHeader,
		QueryStringParameters: queryParameters,
		StageVariables:        map[string]string{},
		Body:                  input,
		IsBase64Encoded:       isBase64Encoded,
		RequestContext: APIGatewayProxyRequestContext{
			Stage:      "",
			HTTPMethod: req.Method,
		},
	}
}

// HydrateHttpResponse will try to fill the response writer with content of body. Addtionaly it adds
// Content-Length header, this has to by done always before any call to Write on the body.
func HydrateHTTPResponse(resp http.ResponseWriter, body json.RawMessage, statusCode int) {
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

func IsRejectedRequest(request *http.Request) bool {
	return request.URL.Path == "/favicon.ico" || request.URL.Path == "/robots.txt"
}

func SetHeaders(input, output http.Header) {
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
