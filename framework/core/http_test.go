package core

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatEventHttp(t *testing.T) {
	t.Parallel()

	testHeaders := http.Header{
		"array":  []string{"val1", "val2"},
		"comma":  []string{"val1,val2"},
		"single": []string{"val1"},
		"empty":  []string{},
	}

	req := http.Request{
		Header: testHeaders,
		// to avoid error in url generation
		URL:    &url.URL{Path: "127.0.0.1"},
		Method: http.MethodPost,
	}

	testBody := []byte("testing")

	ret := FormatEventHTTP(&req, testBody)
	assert.Equal(t, "val1,val2", ret.Headers["array"])
	assert.Equal(t, "val1,val2", ret.Headers["comma"])
	assert.Equal(t, "val1", ret.Headers["single"])
	_, ok := ret.Headers["empty"]
	assert.False(t, ok)

	assert.Equal(t, "127.0.0.1", req.URL.Path)

	assert.Equal(t, http.MethodPost, req.Method)

	assert.False(t, ret.IsBase64Encoded)

	assert.Equal(t, testBody, []byte(ret.Body))
}

func TestFormatEventHttpBase64(t *testing.T) {
	t.Parallel()

	testHeaders := http.Header{
		"array":  []string{"val1", "val2"},
		"comma":  []string{"val1,val2"},
		"single": []string{"val1"},
		"empty":  []string{},
	}

	req := http.Request{
		Header: testHeaders,
		// to avoid error in url generation
		URL:    &url.URL{Path: "127.0.0.1"},
		Method: http.MethodPost,
	}

	rawBodyValue := []byte("base64sample")
	// create a base64 string :
	b64content := base64.StdEncoding.EncodeToString(rawBodyValue)

	testBody := []byte(b64content)

	ret := FormatEventHTTP(&req, testBody)

	assert.True(t, ret.IsBase64Encoded)

	assert.Equal(t, b64content, ret.Body)

	b64decoded, err := base64.StdEncoding.DecodeString(ret.Body)
	require.NoError(t, err)
	assert.Equal(t, b64decoded, rawBodyValue)
}
