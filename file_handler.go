package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xhenner/mp3-go"
)

type DirectoryList struct {
	Directories []string `json:"directories"`
	Count       int      `json:"count"`
}

func (c *Config) directoryList(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("parent")
	if p == "" {
		p = c.CurrentPath
		if p == "" {
			p = c.LibraryPath
		}
	} else {
		p = c.LibraryPath + "/" + p
	}
	dl := DirectoryList{Directories: make([]string, 0), Count: 0}
	entries, err := os.ReadDir(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			dl.Directories = append(dl.Directories, entry.Name())
			dl.Count++
		}
	}
	b, err := json.Marshal(dl)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
	}
}

type Cover struct {
	Data []byte `json:"data"`
	MIME string `json:"mime"`
}

type SongMetadata struct {
	Path     string `json:"path,omitempty"`
	Title    string `json:"title,omitempty"`
	Artist   string `json:"artist,omitempty"`
	Album    string `json:"album,omitempty"`
	Genre    string `json:"genre,omitempty"`
	Track    string `json:"track,omitempty"`
	Year     string `json:"year,omitempty"`
	Cover    Cover  `json:"cover,omitempty"`
	Composer string `json:"composer,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Lyrics   string `json:"lyrics,omitempty"`
	Disc     string `json:"disc,omitempty"`
	// for changes
	Changed bool   `json:"changed,omitempty"`
	NewName string `json:"new_name,omitempty"`
	// readonly attributes
	FileType    string `json:"file_type,omitempty"`
	Format      string `json:"format,omitempty"`
	Bitrate     int    `json:"bitrate,omitempty"`
	Samplerate  int    `json:"samplerate,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	StereoMode  string `json:"stereo_mode,omitempty"`
	BitrateMode string `json:"bitrate_mode,omitempty"`
}

type SongList struct {
	Songs []SongMetadata `json:"songs"`
	Count int            `json:"count"`
}

func (c *Config) songList(w http.ResponseWriter, r *http.Request) {
	includeCover := r.URL.Query().Get("cover") == "true"
	p := r.URL.Query().Get("parent")
	if p == "" {
		p = c.CurrentPath
		if p == "" {
			p = c.LibraryPath
		}
	} else {
		p = c.LibraryPath + "/" + p
	}
	sl := SongList{Songs: make([]SongMetadata, 0), Count: 0}
	entries, err := os.ReadDir(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".mp3") {
			sm, err := c.songMetadata(p + "/" + entry.Name())
			if err == nil {
				if !includeCover {
					sm.Cover.Data = []byte{}
				}
				sl.Songs = append(sl.Songs, sm)
				sl.Count++
			} else {
				l(SEVERITY_WARN, fmt.Sprintf("Error parsing song metadata for song %s: %s", p+"/"+entry.Name(), err))
			}
		}
	}
	b, err := json.Marshal(sl)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(b)
	}
}

func (c *Config) songMetadata(path string) (SongMetadata, error) {
	song := SongMetadata{Path: path}
	mp3File, err := mp3.Examine(path, false)
	if err != nil {
		return song, err
	}
	song.Bitrate = mp3File.Bitrate
	song.Samplerate = mp3File.Sampling
	song.Duration = int(mp3File.Length)
	song.StereoMode = mp3File.Mode
	song.BitrateMode = mp3File.Type

	return c.parseID3(song)
}
