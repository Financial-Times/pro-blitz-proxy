package main

import (
	"net/http"
	"regexp"
)

type NocacheProxy struct {
	proxy http.Handler
	list  []*regexp.Regexp
}

func (p *NocacheProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, re := range p.list {
		match := re.MatchString(req.URL.String())
		if match {
			req.Header.Set("Cache-Control", "no-cache")
			break
		}
	}
	p.proxy.ServeHTTP(w, req)
}
