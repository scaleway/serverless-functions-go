package functest

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
	httpReq.Header.Add("Forwarded", "for="+httpReq.Host+";proto=http")
	httpReq.Header.Add("K-Proxy-Request", "activator")
	httpReq.Header.Add("X-Forwarded-For", httpReq.Host)
	httpReq.Header.Add("X-Forwarded-For", "127.0.0.1")
	httpReq.Header.Add("X-Forwarded-For", "127.0.0.2")
	httpReq.Header.Add("X-Forwarded-Proto", "http")
	httpReq.Header.Add("X-Request-Id", reqID.String())
	httpReq.Header.Add("X-Envoy-External-Address", httpReq.Host)
}

// InjectEgressHeaders simulates the infrastructure output layer where your FaaS will be deployed.
func InjectEgressHeaders(httpResp http.ResponseWriter) {
	httpResp.Header().Set("server", "envoy")
}
