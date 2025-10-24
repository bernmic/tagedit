package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dhowden/tag"
)

func (c *Config) parseID3(song SongMetadata) (SongMetadata, error) {
	f, err := os.Open(song.Path)
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("error loading file: %v", err))
		return song, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			l(SEVERITY_ERROR, fmt.Sprintf("error closing file in readData: %v", err))
		}
	}()

	id3tag, err := tag.ReadFrom(f)
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("Error opening mp3 file %s: %v", song.Path, err))
		return song, err
	}
	song.Format = string(id3tag.Format())
	song.FileType = string(id3tag.FileType())
	song.Title = id3tag.Title()
	song.Artist = id3tag.Artist()
	song.Album = id3tag.Album()
	song.Genre = id3tag.Genre()
	song.Track = formatTrack(id3tag)
	song.Year = strconv.Itoa(id3tag.Year())
	if p := id3tag.Picture(); p != nil {
		song.Cover = Cover{Data: p.Data, MIME: p.MIMEType}
	}
	song.Composer = id3tag.Composer()
	song.Comment = id3tag.Comment()
	song.Lyrics = id3tag.Lyrics()
	song.Disc = formatDisc(id3tag)
	return song, nil
}

func formatTrack(metadata tag.Metadata) string {
	track, total := metadata.Track()
	if total == 0 && track == 0 {
		return ""
	}
	if total == 0 {
		return strconv.Itoa(track)
	}
	return fmt.Sprintf("%d/%d", track, total)
}

func formatDisc(metadata tag.Metadata) string {
	disc, total := metadata.Disc()
	if total == 0 && disc == 0 {
		return ""
	}
	if total == 0 {
		return strconv.Itoa(disc)
	}
	return fmt.Sprintf("%d/%d", disc, total)
}

// GetCoverFromID3 reads the cover image from the ID3 tags.
func GetCoverFromID3(filename string) ([]byte, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("error opening file: %v", err))
		return nil, "", err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			l(SEVERITY_ERROR, fmt.Sprintf("error closing file for cover: %v", err))
		}
	}()

	id3tag, err := tag.ReadFrom(f)
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("ERROR Error reading mp3 file: %v", err))
		return nil, "", err
	}
	if p := id3tag.Picture(); p != nil {
		return p.Data, p.MIMEType, nil
	}
	l(SEVERITY_WARN, fmt.Sprintf("No cover found in ID3: "+filename))
	return nil, "", errors.New("no cover found")
}

func getRating(id3tag tag.Metadata) int {
	ratingsBunch := id3tag.Raw()["POPM"]
	if ratingsBunch != nil {
		us := ratingsBunch.([]uint8)
		for i, u := range us {
			if u == 0 {
				return int(us[i+1])
			}
		}
	}
	return 0
}

func (c *Config) updateID3(song SongMetadata) error {
	l(SEVERITY_INFO, fmt.Sprintf("Updating song: %s", song.Path))
	return nil
}
