package main

import (
	"embed"
	"net/http"
	"strings"
)

//go:embed assets
var assets embed.FS

type Asset struct {
	handler http.Handler
}

func (a *Asset) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
		r.URL.Path = path
	}
	if f, err := assets.Open(assetFile(path)); err == nil {
		f.Close()
		handleAssets(w, r)
		return
	}
	a.handler.ServeHTTP(w, r)
}

func NewAssetsMiddleware(next http.Handler) *Asset {
	return &Asset{next}
}

func assetFile(u string) string {
	return "assets" + u
}

func handleAssets(w http.ResponseWriter, r *http.Request) {
	data, err := assets.ReadFile(assetFile(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	lc := strings.ToLower(r.RequestURI)
	switch {
	case strings.HasSuffix(lc, ".htm"), strings.HasSuffix(lc, ".html"):
		w.Header().Add("Content-Type", "text/html")
	case strings.HasSuffix(lc, ".css"):
		w.Header().Add("Content-Type", "text/css")
	case strings.HasSuffix(lc, ".jpg"), strings.HasSuffix(lc, ".jpeg"):
		w.Header().Add("Content-Type", "image/jpeg")
	case strings.HasSuffix(lc, ".png"):
		w.Header().Add("Content-Type", "image/png")
	case strings.HasSuffix(lc, ".gif"):
		w.Header().Add("Content-Type", "image/gif")
	case strings.HasSuffix(lc, ".ico"):
		w.Header().Add("Content-Type", "image/x-icon")
	case strings.HasSuffix(lc, ".js"):
		w.Header().Add("Content-Type", "application/javascript")
	case strings.HasSuffix(lc, ".json"):
		w.Header().Add("Content-Type", "application/json")
	case strings.HasSuffix(lc, ".map"):
		w.Header().Add("Content-Type", "application/json")
	case strings.HasSuffix(lc, ".svg"):
		w.Header().Add("Content-Type", "image/svg+xml")
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
}
