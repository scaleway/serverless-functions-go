package main

import (
	scw "github.com/scaleway/serverless-functions-go/examples/handler"
	"github.com/scaleway/serverless-functions-go/local"
)

func main() {
	// Replace "Handle" with your function handler name if necessary
	local.ServeHandler(scw.Handle, local.WithPort(8080))
}
