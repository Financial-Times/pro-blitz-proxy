package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gotha/blitz-proxy/storage"
	"github.com/gotha/blitz-proxy/storage/dynamodb"
	"github.com/gotha/blitz-proxy/storage/file"
)

func getCacheStore(storeType StoreType) storage.IStorage {
	var cacheStore storage.IStorage
	switch storeType {
	case StoreTypeFile:
		cacheStore = &file.CacheStore{}
	case StoreTypeDynamodb:
		store, err := dynamodb.NewCacheStoreWithEnvConfig()
		if err != nil {
			slog.Error("could not create dynamodb store",
				slog.String("err", err.Error()),
			)
			os.Exit(1)
		}
		err = store.Init()
		if err != nil {
			slog.Error("could not initialize dynamodb store",
				slog.String("err", err.Error()),
			)
			os.Exit(1)
		}
		cacheStore = store
	default:
		slog.Error("unsupported store")
		os.Exit(1)
	}
	return cacheStore
}

func main() {
	conf := NewConfigFromEnv()
	registerSlogDefaultLogger(conf.SystemCode, GetLogLevel())

	cacheStore := getCacheStore(conf.StoreType)

	var proxy http.Handler
	proxy = &Proxy{
		BackendAddr: conf.BackendAddr,
		HTTPClient:  &http.Client{},
	}
	proxy = &CachingProxy{
		proxy: proxy,
		store: cacheStore,
	}
	proxy = &HeaderParsingProxy{proxy}

	if conf.NoCacheList != nil && len(conf.NoCacheList) > 0 {
		proxy = &NocacheProxy{
			proxy: proxy,
			list:  conf.NoCacheList,
		}
	}

	proxy = &CacheBustingProxy{
		proxy: proxy,
		store: cacheStore,
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
