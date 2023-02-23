package main

import (
	scw "github.com/scaleway/serverless-functions-go/examples/handler"
	server "github.com/scaleway/serverless-functions-go/testing"
)

func main() {
	// Replace "Handle" with your function handler name if necessary
	server.ServeHandlerLocally(scw.Handle, server.WithPort(8080))
}
