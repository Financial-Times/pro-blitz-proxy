package main

import (
	"context"
	"net/http"
)

type HeaderParsingProxy struct {
	proxy http.Handler
}

func (p *HeaderParsingProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if items := req.Header["X-Blitz-Cache-Id"]; len(items) > 0 {
		ctx = context.WithValue(ctx, ContextCacheIDKey, items[0])
	}

	if items := req.Header["Cache-Control"]; len(items) > 0 && items[0] == "no-cache" {
		ctx = context.WithValue(ctx, ContextCacheSkipKey, true)
	}

	p.proxy.ServeHTTP(w, req.WithContext(ctx))
}
