package core

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResponseCustomCode(t *testing.T) {
	t.Parallel()

	// testing with common status codes from different categories (1XX, 2XX, etc...)
	statusCodesToTest := [...]int{
		http.StatusEarlyHints,
		http.StatusOK,
		http.StatusNonAuthoritativeInfo,
		http.StatusFound,
		http.StatusBadRequest,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}

	for _, code := range statusCodesToTest {
		httpResp := http.Response{
			StatusCode: code,
		}

		resp, err := GetResponse(&httpResp)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, code, resp.StatusCode)
	}
}

func TestGetResponseDefault(t *testing.T) {
	t.Parallel()

	httpResp := http.Response{}

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetResponseBadInput(t *testing.T) {
	t.Parallel()

	httpResp := http.Response{
		Body: nil,
	}

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "", string(resp.Body))
}

func TestGetResponseAllFields(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader("body_test")
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	httpResp.Header = make(http.Header, 1)
	httpResp.Header.Set("Key", "value")

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "body_test", string(resp.Body))
	assert.Equal(t, "value", resp.Headers["Key"][0])
	assert.False(t, resp.IsBase64Encoded)
}

func TestGetResponseB64Encoded(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`{"isBase64Encoded": true, "body": "body_test"}`)
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	httpResp.Header = make(http.Header, 1)
	httpResp.Header.Set("Key", "value")

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	// Using assert.Contains instead of Equal because as b64 encoded is enabled, clean unmarshall is not performed here.
	assert.Contains(t, string(resp.Body), "body_test")
	assert.Equal(t, "value", resp.Headers["Key"][0])
	assert.True(t, resp.IsBase64Encoded)
}

func TestGetResponseJsonInJsonBody(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`{"body": {"foo":"bar"}}`)
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	httpResp.Header = make(http.Header, 1)
	httpResp.Header.Set("Key", "value")

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.JSONEq(t, string(resp.Body), `{"foo":"bar"}`)
	assert.Equal(t, "value", resp.Headers["Key"][0])
	assert.False(t, resp.IsBase64Encoded)
}

func TestGetResponseJsonInJson(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`{"foo":"bar"}`)
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	httpResp.Header = make(http.Header, 1)
	httpResp.Header.Set("Key", "value")

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.JSONEq(t, string(resp.Body), `{"foo":"bar"}`)
	assert.Equal(t, "value", resp.Headers["Key"][0])
	assert.False(t, resp.IsBase64Encoded)
}

func TestGetResponseHeaderArrays(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`{"foo":"bar"}`)
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	httpResp.Header = make(http.Header, 3)
	httpResp.Header.Add("array", "value1")
	httpResp.Header.Add("array", "value2")
	httpResp.Header.Add("array", "value3")

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.JSONEq(t, string(resp.Body), `{"foo":"bar"}`)
	assert.Len(t, resp.Headers["Array"], 3)
	assert.False(t, resp.IsBase64Encoded)
}

func TestGetResponseNotJson(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`notjson.1`)
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "notjson.1", string(resp.Body))
	assert.False(t, resp.IsBase64Encoded)
}

func TestGetResponsFullJson(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`
	{
		"body": "bodyTest",
        "headers": {
			"flat": "flatval",
			"array": ["val1", "val2"]
		},
        "statusCode": 302
	}`)

	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, `"bodyTest"`, string(resp.Body))
	assert.False(t, resp.IsBase64Encoded)
	assert.Equal(t, []string{"flatval"}, resp.Headers["flat"])
	assert.Equal(t, []string{"val1", "val2"}, resp.Headers["array"])
}

func TestGetResponseHeaderFlat(t *testing.T) {
	t.Parallel()

	bodyContent := strings.NewReader(`{"foo":"bar"}`)
	bodyContentCloser := io.NopCloser(bodyContent)

	httpResp := http.Response{
		Body:       bodyContentCloser,
		StatusCode: http.StatusFound,
	}

	httpResp.Header = make(http.Header, 1)
	httpResp.Header.Add("array", "value1")

	resp, err := GetResponse(&httpResp)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.JSONEq(t, string(resp.Body), `{"foo":"bar"}`)
	assert.Len(t, resp.Headers["Array"], 1)
	assert.False(t, resp.IsBase64Encoded)
}

func TestGetResponseNil(t *testing.T) {
	t.Parallel()

	resp, err := GetResponse(nil)
	assert.ErrorIs(t, err, ErrRespEmpty)
	assert.Nil(t, resp)
}
