package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Severity string

const (
	SEVERITY_DEBUG Severity = "debug"
	SEVERITY_INFO  Severity = "info"
	SEVERITY_WARN  Severity = "warn"
	SEVERITY_ERROR Severity = "error"
	SEVERITY_FATAL Severity = "fatal"
)

func l(severity Severity, payload string) {
	log.Printf("%s %s\n", severity, payload)
}

func accessLog(r *http.Request, httpCode int, payload string, duration time.Duration) {
	switch httpCode {
	case http.StatusInternalServerError, http.StatusBadRequest:
		l(SEVERITY_ERROR, fmt.Sprintf("error %s, %s %s, %d, %s, %dms", r.RemoteAddr, r.Method, r.RequestURI, httpCode, payload, duration.Milliseconds()))
	case http.StatusNotFound, http.StatusUnauthorized:
		l(SEVERITY_WARN, fmt.Sprintf("warning %s, %s %s, %d, %s, %dms", r.RemoteAddr, r.Method, r.RequestURI, httpCode, payload, duration.Milliseconds()))
	default:
		l(SEVERITY_INFO, fmt.Sprintf("info %s, %s %s, %d, %s, %dms", r.RemoteAddr, r.Method, r.RequestURI, httpCode, payload, duration.Milliseconds()))
	}
}
