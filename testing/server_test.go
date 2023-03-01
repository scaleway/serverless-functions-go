package testing_test

import (
	"fmt"
	"io"
	"net/http"
	stdtest "testing"
	"time"

	"github.com/google/uuid"
	"github.com/scaleway/serverless-functions-go/testing"
	"github.com/stretchr/testify/assert"
)

func TestServ(t *stdtest.T) {
	var handler func(http.ResponseWriter, *http.Request)

	const testingMessage = "simple test"
	handler = func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testingMessage))
	}

	go testing.ServeHandlerLocally(handler, testing.WithPort(49860))

	// req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", nil)

	time.Sleep(2 * time.Second)
	resp, err := http.Get("http://localhost:49860")
	assert.NoError(t, err)

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
	assert.Equal(t, "for=;proto=http", resp.Header.Get("forwarded"))
	assert.Equal(t, "", resp.Header.Get("X-Envoy-External-Address"))
	assert.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, fmt.Sprintf("%d", (len(testingMessage))), resp.Header.Get("Content-Length"))

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, testingMessage, string(bodyBytes))
}
