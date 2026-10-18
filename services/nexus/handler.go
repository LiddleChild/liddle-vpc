package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"strings"

	"github.com/LiddleChild/liddle-vpc/gristguard/grist"
)

type Proxy struct {
	repo  Repository
	proxy *httputil.ReverseProxy
}

func NewProxy(target *url.URL, repo Repository, apiKey string) Proxy {
	return Proxy{
		repo: repo,
		proxy: &httputil.ReverseProxy{
			Rewrite: func(r *httputil.ProxyRequest) {
				r.SetURL(target)
				r.Out.Host = r.In.Host
				r.Out.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
			},
		},
	}
}

func (p Proxy) Handle(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	key, ok := strings.CutPrefix(token, "Bearer ")
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	endpoints, err := p.repo.ListEndpointsByKey(r.Context(), key)
	if err != nil {
		slog.Error(err.Error())
		writeError(w, http.StatusInternalServerError, "gristguard: internal server error")
		return
	}

	if !slices.Contains(endpoints, fmt.Sprintf("%s %s", r.Method, r.Pattern)) {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	p.proxy.ServeHTTP(w, r)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(grist.Error{
		ErrorMessage: message,
	})
}
