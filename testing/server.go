package testing

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/scaleway/serverless-functions-go/framework/function"
)

// Keep the handler to the function .
var localHandler function.ScwFuncV1

// ServeHandlerLocally is the entry point for offline testing. It will serve the handler to a local webserver.
// Read options.go to check advanced paramenter and documentation.
//
// Note that if handler function panics in real life it would make your function return error 500 but
// in order to keep error trace panic will occurs anywhen while using this testing server.
func ServeHandlerLocally(handler function.ScwFuncV1, options ...Option) {
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

	srv := &http.Server{
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 7 * time.Second,
		Handler:           http.HandlerFunc(mainHandlerFnc),
	}

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}

// mainHandlerFnc is the function served by the local server, it will call the localHandler function.
func mainHandlerFnc(httpResp http.ResponseWriter, httpReq *http.Request) {
	defer httpReq.Body.Close()

	CoreProcessing(httpResp, httpReq, localHandler)
}
