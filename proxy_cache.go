package main

import (
	"net/http"
)

type CachingProxy struct {
	store interface {
		Exists(string) bool
		Get(string) ([]byte, map[string][]string, error)
		Save(string, []byte, map[string][]string) error
	}
	proxy http.Handler
}

func (p *CachingProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	skipCache, ok := req.Context().Value(ContextCacheSkipKey).(bool)
	if ok && skipCache {
		w.Header().Add("X-Cache", "miss")
		p.proxy.ServeHTTP(w, req)
		return
	}

	cacheReq := CacheableRequest{req}
	id, err := cacheReq.GetID()
	if err != nil {
		http.Error(w, "failed to read request id", http.StatusInternalServerError)
		return
	}

	if p.store.Exists(id) {
		w.Header().Add("X-Cache", "hit")
		data, _, err := p.store.Get(id)
		if err != nil {
			http.Error(w, "cache hit failed", http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	w.Header().Add("X-Cache", "miss")
	w = &CachingResponseWriter{
		id:         id,
		writer:     w,
		cacheStore: p.store,
	}
	p.proxy.ServeHTTP(w, req)
}
