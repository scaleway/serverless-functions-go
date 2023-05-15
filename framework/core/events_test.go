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

	formattedEvent, err := FormatEvent(&req)

	require.NoError(t, err)
	assert.NotNil(t, formattedEvent)

	assert.Equal(t, testBody, formattedEvent.Body)
	assert.Equal(t, testURLPath, formattedEvent.Path)
	assert.Equal(t, http.MethodPatch, formattedEvent.HTTPMethod)
	assert.Equal(t, http.MethodPatch, formattedEvent.RequestContext.HTTPMethod)

	// headers are flattened so expected result changed
	assert.Equal(t, map[string]string{"multi": "val1,val2"}, formattedEvent.Headers)
	assert.False(t, formattedEvent.IsBase64Encoded)
}
