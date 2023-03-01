package main

import (
	scw "github.com/scaleway/serverless-functions-go/examples/handler"
	"github.com/scaleway/serverless-functions-go/functest"
)

func main() {
	// Replace "Handle" with your function handler name if necessary
	functest.ServeHandlerLocally(scw.Handle, functest.WithPort(8080))
}
