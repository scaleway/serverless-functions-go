package testing

import (
	"net/http"

	"github.com/google/uuid"
)

// InjectIngressHeaders simulates the infrastructure input layer where your FaaS will be deployed.
func InjectIngressHeaders(httpReq *http.Request) {
	reqID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	// headers for function
	httpReq.Header.Set("Forwarded", "for="+httpReq.RemoteAddr+";proto=http")
	httpReq.Header.Set("K-Proxy-Request", "activator")
	httpReq.Header.Set("X-Forwarded-For", httpReq.RemoteAddr)
	httpReq.Header.Set("X-Forwarded-For", "127.0.0.1")
	httpReq.Header.Add("X-Forwarded-For", "127.0.0.2")
	httpReq.Header.Set("X-Forwarded-Proto", "http")
	httpReq.Header.Set("X-Request-Id", reqID.String())
	httpReq.Header.Set("X-Envoy-External-Address", httpReq.RemoteAddr)
}

// InjectEgressHeaders simulates the infrastructure output layer where your FaaS will be deployed.
func InjectEgressHeaders(httpResp http.ResponseWriter) {
	httpResp.Header().Set("server", "envoy")
}
