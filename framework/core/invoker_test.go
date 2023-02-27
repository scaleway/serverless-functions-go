package core

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamRequestBadInput(t *testing.T) {
	t.Parallel()

	fi, err := NewInvoker("", "", "", "", "", false)
	assert.NoError(t, err)
	assert.NotNil(t, fi)

	rtReq := CoreRuntimeRequest{}

	genRequest, err := fi.StreamRequest(rtReq)
	assert.NoError(t, err)

	//nolint:bodyclose
	httpRep, err := fi.client.Do(genRequest)
	assert.Error(t, err)
	assert.Nil(t, httpRep)
}

func TestStreamRequestUserAgent(t *testing.T) {
	t.Parallel()

	const (
		returnString   = "return test from server"
		headerTestKey  = "HEADERTEST"
		headerTestVal  = "HEADERVAL"
		headerTestVal2 = "HEADERVAL2"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerTestKey, headerTestVal)
		w.Header().Add(headerTestKey, headerTestVal2)
		_, errW := w.Write([]byte(returnString))
		assert.NoError(t, errW)
	}))

	t.Cleanup(func() { server.Close() })

	invoker, err := NewInvoker("", "", "", "", server.URL, false)
	require.NoError(t, err)
	assert.NotNil(t, invoker)
}

func TestStreamExecute(t *testing.T) {
	t.Parallel()

	const (
		returnString   = "return test from server"
		headerTestKey  = "HEADERTEST"
		headerTestVal  = "HEADERVAL"
		headerTestVal2 = "HEADERVAL2"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerTestKey, headerTestVal)
		w.Header().Add(headerTestKey, headerTestVal2)
		w.WriteHeader(http.StatusNotFound)
		_, errW := w.Write([]byte(returnString))
		assert.NoError(t, errW)
	}))

	t.Cleanup(func() { server.Close() })

	invoker, err := NewInvoker("", "", "", "", server.URL, false)
	require.NoError(t, err)
	assert.NotNil(t, invoker)

	stringReader := strings.NewReader(returnString)
	stringReadCloser := io.NopCloser(stringReader)
	httpreq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, server.URL, stringReadCloser)
	require.NoError(t, err)

	event, err := FormatEvent(httpreq, TriggerTypeHTTP)
	require.NoError(t, err)
	assert.NotNil(t, event)

	resp, err := invoker.Execute(event, GetExecutionContext(), TriggerTypeHTTP)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
