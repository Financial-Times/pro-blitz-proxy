package main

import (
	"net/http"
	"strings"
)

const cacheBustPrefix = "/__internal/cache/bust"

type CacheBustingProxy struct {
	proxy http.Handler
	store interface {
		DeleteById(id string) error
	}
}

func (p *CacheBustingProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if !strings.HasPrefix(req.URL.Path, cacheBustPrefix) {
		p.proxy.ServeHTTP(w, req)
		return
	}

	uris_str := strings.TrimLeft(req.URL.Path, cacheBustPrefix)
	uris_arr := strings.Split(uris_str, "/")
	if len(uris_arr) == 0 {
		http.Error(w, "Invalid cache busting url", http.StatusInternalServerError)
		return
	}

	if err := p.store.DeleteById(uris_arr[0]); err != nil {
		http.Error(w, "Unable to clear cache", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
