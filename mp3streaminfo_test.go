package main

import (
	"fmt"
	"testing"
)

const (
	fileV = "F:\\tmp\\Audio\\AC-DC - Back In Black.MP3"
	fileC = "F:\\tmp\\Audio\\Armageddon - Buzzard.mp3"
)

func TestStreamInfo(t *testing.T) {
	m, err := StreamInfo(fileC)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(m)
	m, err = StreamInfo(fileV)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m)
}
