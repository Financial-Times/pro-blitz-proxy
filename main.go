package main

import (
	"log/slog"
	"net/http"

	"github.com/gotha/blitz-proxy/storage/file"
)

func main() {
	conf := NewConfigFromEnv()

	registerSlogDefaultLogger(conf.SystemCode, GetLogLevel())

	fileCacheStore := &file.CacheStore{}

	var proxy http.Handler
	proxy = &Proxy{BackendAddr: conf.BackendAddr}
	proxy = &CachingProxy{
		proxy: proxy,
		store: fileCacheStore,
	}
	proxy = &HeaderParsignProxy{proxy}
	proxy = &CacheBustingProxy{
		proxy: proxy,
		store: fileCacheStore,
	}
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
