package testing_test

import (
	"testing"

	scw "github.com/scaleway/serverless-functions-go/testing"
)

func TestOptions(t *testing.T) {
	options := [...]scw.Option{scw.WithPort(123)}

	server := scw.Server{}

	for idx := range options {
		options[idx](&server)
	}
}
