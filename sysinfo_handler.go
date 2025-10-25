package main

import (
	"encoding/json"
	"net/http"
	"runtime"
	"runtime/debug"
)

type SysInfo struct {
	Version     string `json:"version"`
	CommitTime  string `json:"commit_time"`
	Os          string `json:"os"`
	GoVersion   string `json:"go_version"`
	LibraryPath string `json:"library_path"`
}

func (c *Config) sysInfo(w http.ResponseWriter, r *http.Request) {
	si := SysInfo{Version: VERSION, Os: runtime.GOOS, LibraryPath: c.LibraryPath, GoVersion: runtime.Version()}
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.time" {
				si.CommitTime = s.Value
			}
		}
	}
	b, err := json.Marshal(si)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
	}
}
