package testing

import (
	"net/http"

	"github.com/google/uuid"
)

// IngressProcessing simulates the infrastructure layer where your FaaS will be deployed.
func IngressProcessing(httpReq *http.Request) {
	reqId, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}

	// headers for function
	httpReq.Header.Add("Forwarded", "for="+httpReq.RemoteAddr+";proto=http")
	httpReq.Header.Add("K-Proxy-Request", "activator")
	httpReq.Header.Add("X-Forwarded-For", httpReq.RemoteAddr)
	httpReq.Header.Add("X-Forwarded-For", "127.0.0.1")
	httpReq.Header.Add("X-Forwarded-For", "127.0.0.2")
	httpReq.Header.Add("X-Forwarded-Proto", "http")
	httpReq.Header.Add("X-Request-Id", reqId.String())
	httpReq.Header.Add("X-Envoy-External-Address", httpReq.RemoteAddr)
}
