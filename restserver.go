package main

import (
	"fmt"
	"log"
	"net/http"
)

func (c *Config) startHttpListener() {
	c.mux = http.NewServeMux()
	c.InitRouting()
	c.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	l("info", fmt.Sprintf("Starting http server on port %d", c.Port))
	cors := NewCorsMiddleware(c.mux)
	assets := NewAssetsMiddleware(cors)
	al := NewAccesslogMiddleware(assets)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%04d", c.Port), al))
}

func (c *Config) InitRouting() {
	c.mux.HandleFunc("GET /api/directories", c.directoryList)
	c.mux.HandleFunc("GET /api/songs", c.songList)
}
