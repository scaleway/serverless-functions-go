package framework

import "net/http"

func InfraProcessing(httpResp http.ResponseWriter, httpReq *http.Request) {
	httpResp.Header().Add("X-Forwarded-For", httpReq.RemoteAddr)
	httpResp.Header().Add("X-Forwarded-Proto", "https")
	httpResp.Header().Add("X-Envoy-External-Address", httpReq.RemoteAddr)
}
