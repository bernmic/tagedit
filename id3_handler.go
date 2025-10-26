package main

import (
	"fmt"

	"github.com/bogem/id3v2/v2"
)

func (c *Config) parseID3New(song SongMetadata) (SongMetadata, error) {
	id3tag, err := id3v2.Open(song.Path, id3v2.Options{Parse: true})
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("error while opening mp3 file %s: %v", song.Path, err))
	}
	defer func() {
		err := id3tag.Close()
		if err != nil {
			l(SEVERITY_ERROR, fmt.Sprintf("error closing file in readData: %v", err))
		}
	}()
	song.Format = fmt.Sprintf("ID3V2.%d", id3tag.Version())
	song.FileType = id3tag.GetTextFrame(TFLT).Text
	song.HasID3V2 = id3tag.HasFrames()
	song.Title = id3tag.Title()
	song.Artist = id3tag.Artist()
	song.AlbumArtist = id3tag.GetTextFrame(TPE2).Text
	song.Album = id3tag.Album()
	song.Year = id3tag.Year()
	song.Genre = id3tag.Genre()
	song.Track = id3tag.GetTextFrame(TRCK).Text
	song.Composer = id3tag.GetTextFrame(TCOM).Text
	song.Lyrics = id3tag.GetTextFrame(TEXT).Text
	song.Disc = id3tag.GetTextFrame(TMED).Text
	frames := id3tag.GetFrames(COMM)
	if len(frames) > 0 {
		for _, frame := range frames {
			comment, ok := frame.(id3v2.CommentFrame)
			if ok {
				song.Comment = comment.Text
			}
		}
	}
	frames = id3tag.GetFrames(APIC)
	if len(frames) > 0 {
		covers := make([]Cover, 0)
		for _, frame := range frames {
			picture, ok := frame.(id3v2.PictureFrame)
			if !ok {
				l(SEVERITY_WARN, fmt.Sprintf("invalid frame type: %v", frame))
				continue
			}
			covers = append(covers, Cover{Data: picture.Picture, MIME: picture.MimeType, Description: picture.Description})
		}
		if len(covers) == 1 {
			song.Cover = covers[0]
		} else {
			for i, cover := range covers {
				l(SEVERITY_INFO, fmt.Sprintf("cover %d: %s", i, cover))
			}
			song.Cover = covers[0]
		}
	}
	return song, nil
}

func (c *Config) updateID3(song SongMetadata) error {
	l(SEVERITY_INFO, fmt.Sprintf("Updating song: %s", song.Path))
	return nil
}
