package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type HTTPClientMock struct {
	DoF func(*http.Request) (*http.Response, error)
}

func (m *HTTPClientMock) Do(r *http.Request) (*http.Response, error) {
	if m.DoF != nil {
		return m.DoF(r)
	}
	return nil, nil
}

func TestProxy(t *testing.T) {
	tests := map[string]struct {
		backendAddr  string
		client       HTTPClient
		expectedResp *http.Response
	}{
		"test that error is returned when bad backend url is configured": {
			backendAddr: "http://%xyz",
			client: &HTTPClientMock{
				DoF: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       io.NopCloser(strings.NewReader("this is response")),
						Header:     map[string][]string{},
						StatusCode: http.StatusOK,
					}, nil
				},
			},
			expectedResp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("invalid backend address: 'http://%xyz'\n")),
				StatusCode: http.StatusInternalServerError,
			},
		},
		"test that error is returned when http client returns one": {
			backendAddr: "http://example.com",
			client: &HTTPClientMock{
				DoF: func(*http.Request) (*http.Response, error) {
					return &http.Response{}, fmt.Errorf("random error")
				},
			},
			expectedResp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("Bad Request\n")),
				StatusCode: http.StatusBadRequest,
			},
		},
		"test that response is returned from http client preserving headers": {
			backendAddr: "http://example.com",
			client: &HTTPClientMock{
				DoF: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       io.NopCloser(strings.NewReader("this is response")),
						Header:     map[string][]string{"X-Custom-Header": {"value1"}},
						StatusCode: http.StatusOK,
					}, nil
				},
			},
			expectedResp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("this is response")),
				Header:     map[string][]string{"X-Custom-Header": {"value1"}},
				StatusCode: http.StatusOK,
			},
		},
		"test that errors response from client is also returned as it is": {
			backendAddr: "http://example.com",
			client: &HTTPClientMock{
				DoF: func(*http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       io.NopCloser(strings.NewReader("this is error")),
						Header:     map[string][]string{"X-Custom-Header": {"value2"}},
						StatusCode: http.StatusInternalServerError,
					}, nil
				},
			},
			expectedResp: &http.Response{
				Body:       io.NopCloser(strings.NewReader("this is error")),
				Header:     map[string][]string{"X-Custom-Header": {"value2"}},
				StatusCode: http.StatusInternalServerError,
			},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			proxy := Proxy{
				BackendAddr: test.backendAddr,
				HTTPClient:  test.client,
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			proxy.ServeHTTP(w, req)
			res := w.Result()
			defer res.Body.Close()
			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			expectedData, err := io.ReadAll(test.expectedResp.Body)
			if err != nil {
				t.Fatalf("could not read expected response body: %s", err)
			}
			if string(data) != string(expectedData) {
				t.Errorf("expected res '%v' got '%v'", test.expectedResp, string(data))
			}
			if res.StatusCode != test.expectedResp.StatusCode {
				t.Errorf("expected status code '%v'; got '%v'", test.expectedResp.StatusCode, res.StatusCode)
			}
		})
	}
}
