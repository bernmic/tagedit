package main

import (
	"net/http"
)

type Cors struct {
	handler http.Handler
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-type")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	c.handler.ServeHTTP(w, r)
}

func NewCorsMiddleware(next http.Handler) *Cors {
	return &Cors{next}
}
