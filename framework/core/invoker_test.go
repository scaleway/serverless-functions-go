package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestFunctionInvoker_streamRequest_userAgent(t *testing.T) {
// 	t.Parallel()

// 	// write User-Agent to body to compare more easily
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, r.Header.Get(userAgentHeaderKey))
// 	}))
// 	defer server.Close()

// 	fn := FunctionInvoker{
// 		client:   http.DefaultClient,
// 	}

// 	type args struct {
// 		reqBody CoreRuntimeRequest
// 	}

// 	tests := []struct {
// 		name     string
// 		args     args
// 		wantBody string
// 	}{
// 		{
// 			name: "user agent is set and forwarded",
// 			args: args{
// 				reqBody: CoreRuntimeRequest{
// 					Event: events.APIGatewayProxyRequest{
// 						Headers: map[string]string{
// 							"User-Agent": "my user agent",
// 						},
// 						MultiValueHeaders: map[string][]string{
// 							"multi-header": {"h1", "h2"},
// 						},
// 					},
// 					TriggerType: events.TriggerTypeHTTP,
// 				},
// 			},
// 			wantBody: "my user agent",
// 		},
// 		{
// 			name: "user agent is not set, defaulting to http package user agent",
// 			args: args{
// 				reqBody: CoreRuntimeRequest{
// 					Event: events.APIGatewayProxyRequest{
// 						Headers: map[string]string{},
// 					},
// 					TriggerType: events.TriggerTypeHTTP,
// 				},
// 			},
// 			// default user agent when calling from a Go program
// 			// https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/http/request.go;l=512
// 			wantBody: "Go-http-client/1.1",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotRes, err := fn.streamRequest(tt.args.reqBody)
// 			require.NoError(t, err)

// 			body, err := io.ReadAll(gotRes.Body)
// 			require.NoError(t, err)

// 			bodyString := string(body)
// 			assert.Equal(t, tt.wantBody, bodyString)
// 		})
// 	}
// }

// func TestFunctionInvoker_streamRequest_contentType(t *testing.T) {
// 	t.Parallel()

// 	// write Content-Type to body to compare more easily
// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, r.Header.Get(contentTypeHeaderKey))
// 	}))
// 	defer server.Close()

// 	fn := FunctionInvoker{
// 		client:   http.DefaultClient,
// 		subrtURL: server.URL,
// 	}

// 	type args struct {
// 		reqBody CoreRuntimeRequest
// 	}

// 	tests := []struct {
// 		name     string
// 		args     args
// 		wantBody string
// 	}{
// 		{
// 			name: "content type is set and forwarded",
// 			args: args{
// 				reqBody: CoreRuntimeRequest{
// 					Event: events.APIGatewayProxyRequest{
// 						Headers: map[string]string{
// 							"Content-Type": "text/plain",
// 						},
// 						MultiValueHeaders: map[string][]string{
// 							"multi-header": {"h1", "h2"},
// 						},
// 					},
// 					TriggerType: events.TriggerTypeHTTP,
// 				},
// 			},
// 			wantBody: "text/plain",
// 		},
// 		{
// 			name: "content type is not set, defaulting to application/json",
// 			args: args{
// 				reqBody: CoreRuntimeRequest{
// 					Event: events.APIGatewayProxyRequest{
// 						Headers: map[string]string{},
// 					},
// 					TriggerType: events.TriggerTypeHTTP,
// 				},
// 			},
// 			wantBody: "application/json",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotRes, err := fn.streamRequest(tt.args.reqBody)
// 			require.NoError(t, err)

// 			body, err := io.ReadAll(gotRes.Body)
// 			require.NoError(t, err)

// 			bodyString := string(body)
// 			assert.Equal(t, tt.wantBody, bodyString)
// 		})
// 	}
// }

func TestStreamRequestBadInput(t *testing.T) {
	t.Parallel()

	fi, err := NewInvoker("", "", "", "", "", false)
	assert.NoError(t, err)
	assert.NotNil(t, fi)

	rtReq := CoreRuntimeRequest{}
	httpRep, err := fi.streamRequest(rtReq)
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

	// rtReq := CoreRuntimeRequest{
	// 	TriggerType: TriggerTypeHTTP,
	// }

	// httpRep, err := invoker.streamRequest(rtReq)
	// require.NoError(t, err)
	// assert.NotNil(t, httpRep)
	// assert.Equal(t, http.StatusOK, httpRep.StatusCode)

	// defer httpRep.Body.Close()

	// body, err := io.ReadAll(httpRep.Body)
	// require.NoError(t, err)
	// assert.Equal(t, returnString, string(body))
	// assert.Contains(t, httpRep.Header.Values(headerTestKey), headerTestVal)
	// assert.Contains(t, httpRep.Header.Values(headerTestKey), headerTestVal2)
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
	httpreq, err := http.NewRequest(http.MethodPost, server.URL, stringReadCloser)
	require.NoError(t, err)

	event, err := FormatEvent(httpreq, TriggerTypeHTTP)
	require.NoError(t, err)
	assert.NotNil(t, event)

	resp, err := invoker.Execute(event, GetExecutionContext(), TriggerTypeHTTP)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	// assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// defer resp.Body.Close()
	// bodyBytes, err := io.ReadAll(resp.Body)
	// require.NoError(t, err)
	// assert.Equal(t, []byte(returnString), bodyBytes)
	// assert.Equal(t, strconv.Itoa(len(returnString)), resp.Header.Get("Content-Length"))
}
