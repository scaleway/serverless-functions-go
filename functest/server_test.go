package functest_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/scaleway/serverless-functions-go/functest"
	"github.com/stretchr/testify/assert"
)

func TestServSimpleResponse(t *testing.T) {
	t.Parallel()

	const testingMessage = "simple test"

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testingMessage))

		assert.NotEmpty(t, r)
		assert.Contains(t, r.Header.Get("Forwarded"), "proto=http")
		assert.Equal(t, "activator", r.Header.Get("K-Proxy-Request"))
		assert.Equal(t, "http", r.Header.Get("X-Forwarded-Proto"))
	}

	go functest.ServeHandlerLocally(handler, functest.WithPort(49860))

	time.Sleep(2 * time.Second)

	//nolint:noctx
	resp, err := http.Get("http://localhost:49860")
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "gzip", resp.Header.Get("Accept-Encoding"))
	assert.Equal(t, "activator", resp.Header.Get("K-Proxy-Request"))
	assert.Equal(t, "http", resp.Header.Get("X-Forwarded-Proto"))
	assert.Len(t, resp.Header.Get("X-Request-Id"), len(uuid.New().String()))
	assert.Len(t, resp.Header.Values("X-Forwarded-For"), 3)
	assert.NotEmpty(t, resp.Header.Get("Date"))
	assert.Equal(t, "envoy", resp.Header.Get("server"))
	assert.NotEmpty(t, resp.Header.Get("User-Agent"))
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("forwarded"), "proto=http")
	assert.NotEmpty(t, resp.Header.Get("X-Envoy-External-Address"))
	assert.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, fmt.Sprintf("%d", (len(testingMessage))), resp.Header.Get("Content-Length"))

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, testingMessage, string(bodyBytes))
}

func TestServDumpResponse(t *testing.T) {
	t.Parallel()

	handler := func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		assert.NoError(t, err)
		fmt.Fprintf(w, "%s\n", string(dump))
	}

	go functest.ServeHandlerLocally(handler, functest.WithPort(49861))

	time.Sleep(2 * time.Second)

	//nolint:noctx
	resp, err := http.Get("http://localhost:49861")
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.NotNil(t, resp)

	respBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NotEmpty(t, respBytes)
	respStr := string(respBytes)

	assert.Contains(t, respStr, "Host:")
	assert.Contains(t, respStr, "proto=http")
	assert.Contains(t, respStr, "K-Proxy-Request: activator")
	assert.Contains(t, respStr, "X-Envoy-External-Address:")
	assert.Contains(t, respStr, "X-Forwarded-For:")
	assert.Contains(t, respStr, "X-Forwarded-Proto: http")
	assert.Contains(t, respStr, "X-Request-Id")
}
