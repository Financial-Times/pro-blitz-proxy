package main

import (
	"log/slog"
	"net/http"
)

type LoggingProxy struct {
	proxy http.Handler
}

func (p *LoggingProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	slog.Info("request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
	)
	p.proxy.ServeHTTP(w, req)
}
