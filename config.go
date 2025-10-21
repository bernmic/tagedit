package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	LibraryPath string `json:"library_path,omitempty" yaml:"library_path,omitempty"`
	CurrentPath string `json:"current_path,omitempty" yaml:"current_path,omitempty"`
	Port        int    `json:"port,omitempty" yaml:"port,omitempty"`
	mux         *http.ServeMux
}

func NewConfig() *Config {
	// init default values
	c := new(Config)
	c.LibraryPath = "/"
	c.Port = 8080
	// load tagedit.yaml
	b, err := os.ReadFile("tagedit.yaml")
	if err != nil {
		b, err = os.ReadFile("tagedit.yml")
	}
	if b != nil {
		err = yaml.Unmarshal(b, c)
		if err != nil {
			fmt.Printf("error parsing config file: %v\n", err)
		}
	}

	// overwrite with env vars
	if l, ok := os.LookupEnv("TAGEDIT_LIBRARY_PATH"); ok && l != "" {
		c.LibraryPath = l
	}
	if l, ok := os.LookupEnv("TAGEDIT_CURRENT_PATH"); ok && l != "" {
		c.CurrentPath = l
	}
	if l, ok := os.LookupEnv("TAGEDIT_PORT"); ok && l != "" {
		c.Port, err = strconv.Atoi(l)
		if err != nil {
			fmt.Printf("error parsing environment variable TAGEDIT_PORT: %v\n", err)
			c.Port = 8080
		}
	}
	return c
}
