package testing

import "net/http"

// InfraProcessing simulates the infrastructure layer where your FaaS will be deployed.
func InfraProcessing(httpResp http.ResponseWriter, httpReq *http.Request) {
	httpResp.Header().Add("X-Forwarded-For", httpReq.RemoteAddr)
	httpResp.Header().Add("X-Forwarded-Proto", "https")
	httpResp.Header().Add("X-Envoy-External-Address", httpReq.RemoteAddr)
	httpResp.Header().Add("x-envoy-upstream-service-time", "0")
}
