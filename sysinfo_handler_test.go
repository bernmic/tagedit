package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSysinfoHandler(t *testing.T) {
	t.Run("GET Sysinfo", func(t *testing.T) {
		config := NewConfig()
		request := httptest.NewRequest(http.MethodGet, "/api", nil)
		response := httptest.NewRecorder()
		config.sysInfo(response, request)
		if status := response.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		si := SysInfo{}
		err := json.Unmarshal(response.Body.Bytes(), &si)
		if err != nil {
			t.Errorf("could not unmarshal body: %v", err)
		}
		if si.Version != VERSION {
			t.Errorf("handler returned wrong version: got %v want %v", si.Version, VERSION)
		}
		if si.LibraryPath != config.LibraryPath {
			t.Errorf("handler returned wrong library path: got %v want %v", si.LibraryPath, config.LibraryPath)
		}
	})
}
