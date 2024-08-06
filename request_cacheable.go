package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type CacheableRequest struct {
	*http.Request
}

func (r *CacheableRequest) GetID() (string, error) {
	cacheID, _ := r.Context().Value(ContextCacheIDKey).(string)
	if cacheID != "" {
		return cacheID, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	h := md5.New()
	h.Write([]byte(r.URL.Path))
	h.Write(body)

	hashInBytes := h.Sum(nil)
	return hex.EncodeToString(hashInBytes), nil
}
