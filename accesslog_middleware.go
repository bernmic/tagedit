package main

import (
	"fmt"
	"net/http"
	"time"
)

type AccessLogMiddleware struct {
	handler http.Handler
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func NewAccesslogMiddleware(next http.Handler) *AccessLogMiddleware {
	return &AccessLogMiddleware{handler: next}
}

func (p *AccessLogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	rw := StatusRecorder{ResponseWriter: w}
	p.handler.ServeHTTP(&rw, r)
	accessLog(r, rw.Status, payload(rw.Status), time.Now().Sub(startTime))
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func payload(status int) string {
	switch status {
	case 200:
		return "ok"
	case 201:
		return "created"
	case 404:
		return "not found"
	case 500:
		return "internal server error"
	default:
		return fmt.Sprintf("%d", status)
	}
}
