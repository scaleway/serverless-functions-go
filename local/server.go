package local

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/scaleway/serverless-functions-go/framework/function"
)

// ServeHandlerLocally is the entry point for offline testing. It will serve the handler to a local webserver.
// Read options.go to check advanced paramenter and documentation.
//
// Note that if handler function panics in real life it would make your function return error 500 but
// in order to keep error trace panic will occurs anywhen while using this testing server.
func ServeHandlerLocally(handler function.ScwFuncV1, options ...Option) {
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

	decoratedHandler := func(httpResp http.ResponseWriter, httpReq *http.Request) {
		CoreProcessing(httpResp, httpReq, handler)

		if httpReq.Body != nil && httpReq.Body != http.NoBody {
			httpReq.Body.Close()
		}
	}

	srv := &http.Server{
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 7 * time.Second,
		Handler:           http.HandlerFunc(decoratedHandler),
	}

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
