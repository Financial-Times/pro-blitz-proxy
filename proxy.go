package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ContextCacheHeaderKey string

const (
	ContextCacheIDKey   ContextCacheHeaderKey = "id"
	ContextCacheSkipKey ContextCacheHeaderKey = "skip"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Proxy struct {
	BackendAddr string
	HTTPClient  HTTPClient
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	u, err := url.Parse(p.BackendAddr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid backend address: '%s'", p.BackendAddr), http.StatusInternalServerError)
		return
	}
	u.Path = req.URL.Path
	req.URL = u

	// http: Request.RequestURI can't be set in client requests.
	// http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
