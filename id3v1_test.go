package main

import (
	"fmt"
	"testing"
)

const (
	file1 = "F:\\tmp\\Audio\\AC-DC - Back In Black.MP3"
	file2 = "F:\\tmp\\Audio\\Armageddon - Buzzard.mp3"
)

func TestOpen(t *testing.T) {
	v1, err := Open(file1)

	if err != nil {
		t.Fatalf("error opening file: %v\n", err)
	}
	fmt.Println(v1)
}

func TestHasID3V1(t *testing.T) {
	b := HasID3V1(file1)
	if !b {
		t.Fatalf("expected true for file1")
	}
	b = HasID3V1(file2)
	if b {
		t.Fatalf("expected false for file1")
	}
}

func TestRemoveID3V1(t *testing.T) {
	err := RemoveID3V1(file1, false)
	if err != nil {
		t.Fatalf("Error %v", err)
	}
}
