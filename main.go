package main

import (
	"log/slog"
	"net/http"
)

func main() {
	conf := NewConfigFromEnv()

	registerSlogDefaultLogger(conf.SystemCode, GetLogLevel())

	var proxy http.Handler
	proxy = &Proxy{BackendAddr: conf.BackendAddr}
	proxy = &LoggingProxy{proxy}

	addr := conf.GetAddress()
	slog.Info("Starting proxy server on",
		slog.String("address", addr),
		slog.String("backend", conf.BackendAddr),
	)
	if err := http.ListenAndServe(addr, proxy); err != nil {
		slog.Error("Listen and serve", slog.String("err", err.Error()))
	}
}
