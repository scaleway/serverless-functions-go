package core

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTriggerDefault(t *testing.T) {
	t.Parallel()

	tt, err := GetTriggerType("")
	require.NoError(t, err)
	assert.Equal(t, TriggerTypeHTTP, tt)
}

func TestGetTriggerInvalid(t *testing.T) {
	t.Parallel()

	tt, err := GetTriggerType("invalid")
	assert.Error(t, ErrNotSupportedTrigger, err)
	assert.Equal(t, TriggerType(""), tt)
}

func TestGetTriggerValidHTTP(t *testing.T) {
	t.Parallel()

	tt, err := GetTriggerType(string(TriggerTypeHTTP))
	require.NoError(t, err)
	assert.Equal(t, TriggerTypeHTTP, tt)
}

func TestFormatEventGood(t *testing.T) {
	t.Parallel()

	const (
		testBody    = "sampleBody"
		testURLPath = "127.sample.path"
	)

	testURL := url.URL{Path: testURLPath}

	sampleHeaders := map[string][]string{
		"multi": {"val1", "val2"},
	}

	bodyReader := io.NopCloser(strings.NewReader(testBody))
	req := http.Request{
		Method: http.MethodPatch,
		Body:   bodyReader,
		Header: sampleHeaders,
		URL:    &testURL,
	}

	formattedEvent, err := FormatEvent(&req, TriggerTypeHTTP)

	require.NoError(t, err)
	assert.NotNil(t, formattedEvent)

	// try cast event

	castedEvt, castSucceeds := formattedEvent.(APIGatewayProxyRequest)
	require.True(t, castSucceeds)
	assert.NotNil(t, castedEvt)

	assert.Equal(t, testBody, castedEvt.Body)
	assert.Equal(t, testURLPath, castedEvt.Path)
	assert.Equal(t, http.MethodPatch, castedEvt.HTTPMethod)
	assert.Equal(t, http.MethodPatch, castedEvt.RequestContext.HTTPMethod)

	// headers are flattened so expected result changed
	assert.Equal(t, map[string]string{"multi": "val1,val2"}, castedEvt.Headers)
	assert.False(t, castedEvt.IsBase64Encoded)
}
