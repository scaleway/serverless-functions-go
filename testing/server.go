package testing

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/scaleway/serverless-functions-go/framework/function"
)

var (
	localHandler function.ScwFuncV1
	latency      time.Duration
)

// ScalewayRouter is the entry point for offline testing. It will serve
func ScalewayRouter(handler function.ScwFuncV1, options ...Option) {
	localHandler = handler

	server := Server{
		port: "0",
	}

	for idx := range options {
		options[idx](&server)
	}

	listener, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		panic(err)
	}

	fmt.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)

	if err := http.Serve(listener, http.HandlerFunc(mainHandlerFnc)); err != nil {
		panic(err)
	}
}

func mainHandlerFnc(httpResp http.ResponseWriter, httpReq *http.Request) {
	defer httpReq.Body.Close()

	CoreProcessing(httpResp, httpReq, localHandler)

	httpResp.Header().Add("x-envoy-upstream-service-time", latency.String())

	time.Sleep(latency)
}
