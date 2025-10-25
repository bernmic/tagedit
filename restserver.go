package main

import (
	"encoding/json"
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
	c.mux.HandleFunc("GET /api", c.sysInfo)
	c.mux.HandleFunc("GET /api/directories", c.directoryList)
	c.mux.HandleFunc("GET /api/songs", c.songList)
	c.mux.HandleFunc("PATCH /api/songs", c.updateSongs)
}

func (c *Config) updateSongs(w http.ResponseWriter, r *http.Request) {
	sl := SongList{}
	err := json.NewDecoder(r.Body).Decode(&sl)
	if err != nil {
		l(SEVERITY_WARN, fmt.Sprintf("Error decoding songs: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, song := range sl.Songs {
		if song.Changed {
			err = c.updateID3(song)
			if err != nil {
				l(SEVERITY_ERROR, fmt.Sprintf("Error updating song %s: %v", song.Path, err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if song.NewName != "" {
				l(SEVERITY_INFO, fmt.Sprintf("rename file %s to %s", song.Path, song.NewName))
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}
