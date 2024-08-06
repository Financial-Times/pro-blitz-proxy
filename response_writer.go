package main

import "net/http"

type CacheStore interface {
	Save(string, []byte, map[string][]string) error
}

type CachingResponseWriter struct {
	id         string
	writer     http.ResponseWriter
	cacheStore CacheStore
}

func (w *CachingResponseWriter) Write(p []byte) (n int, err error) {
	go w.cacheStore.Save(w.id, p, w.writer.Header())
	return w.writer.Write(p)
}

func (w *CachingResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *CachingResponseWriter) WriteHeader(statusCode int) {
	w.writer.WriteHeader(statusCode)
}
