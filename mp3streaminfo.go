package main

import (
	"fmt"
	"io"
	"os"
)

const (
	frameHeaderLength = 4
	id3v2Headerlength = 10
)

var m1l1Bitrates = []uint{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0}
var m1l2Bitrates = []uint{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0}
var m1l3Bitrates = []uint{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0}
var m2l1Bitrates = []uint{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0}
var m2l2Bitrates = []uint{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0}

var m1samplerates = []uint{44100, 48000, 32000, 0}
var m2samplerates = []uint{22050, 24000, 16000, 0}
var m25samplerates = []uint{11025, 12000, 8000, 0}

type Mp3StreamInfo struct {
	Bitrate     uint  `json:"bitrate"`
	Samplerate  uint  `json:"samplerate"`
	Duration    uint  `json:"duration"`
	Vbr         bool  `json:"vbr"`
	ChannelMode uint8 `json:"channel_mode"`
}

type MP3Frame struct {
	Version    uint8
	Layer      uint8
	Bitrate    uint
	Samplerate uint
	Framesize  uint
	Channel    uint8
}

// StreamInfo gets stream information from a mp3 file
func StreamInfo(filename string) (*Mp3StreamInfo, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			l(SEVERITY_ERROR, fmt.Sprintf("error closing file for streamInfo: %v", err))
		}
	}()
	return ReadStreamInfo(f)
}

func ReadStreamInfo(f *os.File) (*Mp3StreamInfo, error) {
	var pos int64
	id3header := make([]byte, id3v2Headerlength)
	c, err := f.Read(id3header)
	if err != nil || c != id3v2Headerlength {
		return nil, fmt.Errorf("could not read ID3V2 header: %v", err)
	}
	var ihlen int64
	if string(id3header[:3]) == "ID3" {
		//skip ID3 header
		ihlen = id3headerlen(id3header) + id3v2Headerlength
	}

	pos = ihlen
	frames := make([]MP3Frame, 0)
	for {
		p, err := f.Seek(pos, 0)
		if err != nil || p != pos {
			return nil, fmt.Errorf("cannot position to first frame: %v", err)
		}
		frameheader := make([]byte, frameHeaderLength)
		n, err := f.Read(frameheader)
		if err == io.EOF || n != frameHeaderLength {
			break
		}
		if frameheader[0] != 0xFF || (frameheader[1]&0xE0 != 0xE0) {
			break
		}
		fr := MP3Frame{}
		fr.Version = uint8((frameheader[1] >> 3) & 0x03)
		fr.Layer = uint8((frameheader[1] >> 1) & 0x03)
		//prot := (frameheader[1] >> 0) & 0x01
		br := (frameheader[2] >> 4) & 0x0F
		sr := (frameheader[2]) >> 2 & 0x03
		pad := (frameheader[2] >> 1) & 0x01
		//priv := (frameheader[2] >> 0) & 0x01
		fr.Channel = uint8((frameheader[3] >> 6) & 0x03)
		//me := (frameheader[3] >> 4) & 0x03
		//cr := (frameheader[3] >> 3) & 0x01
		//orig := (frameheader[3] >> 2) & 0x01
		//emph := (frameheader[3] >> 0) & 0x03
		if fr.Version == 0 || fr.Version == 2 {
			if fr.Layer == 3 {
				fr.Bitrate = m2l1Bitrates[br]
			} else {
				fr.Bitrate = m2l2Bitrates[br]
			}

		} else {
			if fr.Layer == 1 {
				fr.Bitrate = m1l3Bitrates[br]
			} else if fr.Layer == 2 {
				fr.Bitrate = m1l2Bitrates[br]
			} else if fr.Layer == 3 {
				fr.Bitrate = m1l1Bitrates[br]
			}
		}

		if fr.Version == 0 {
			fr.Samplerate = m25samplerates[sr]
		} else if fr.Version == 2 {
			fr.Samplerate = m2samplerates[sr]
		} else if fr.Version == 3 {
			fr.Samplerate = m1samplerates[sr]
		}
		if fr.Samplerate == 0 {
			//return nil, fmt.Errorf("could not determine samplerate for version %d, sr %d", fr.Version, sr)
			fr.Samplerate = 44100
		}
		frames = append(frames, fr)
		if fr.Layer == 3 {
			pos += int64((12*fr.Bitrate*1000/fr.Samplerate + uint(pad)) * 4)
		} else {
			pos += int64(144*fr.Bitrate*1000/fr.Samplerate + uint(pad))
		}
	}
	m := Mp3StreamInfo{}

	if len(frames) > 0 {
		br1 := frames[0].Bitrate
		var su uint64
		for _, fra := range frames {
			su += uint64(fra.Bitrate)
			if !m.Vbr && fra.Bitrate != br1 {
				m.Vbr = true
			}
		}
		if m.Vbr {
			m.Bitrate = uint(su / uint64(len(frames)))
		} else {
			m.Bitrate = frames[0].Bitrate
		}
		m.Samplerate = frames[0].Samplerate
		bytespersecond := (m.Bitrate * 1000) / 8
		st, _ := f.Stat()
		m.Duration = uint((st.Size() - ihlen) / int64(bytespersecond))
		m.ChannelMode = frames[0].Channel
	}
	return &m, nil
}

func id3headerlen(b []byte) int64 {
	if len(b) != 10 {
		return 0
	}
	b0 := b[6] & 0x7F
	b1 := b[7] & 0x7F
	b2 := b[8] & 0x7F
	b3 := b[9] & 0x7F

	lb := (b2 & 0x01) << 7
	b3 = b3 | lb

	b2 = b2 >> 1

	lb = (b1 & 0x03) << 6
	b2 = b2 | lb

	b1 = b1 >> 2

	lb = (b0 & 0x07) << 5
	b1 = b1 | lb

	b0 = b0 >> 3

	return int64(uint(b3) | uint(b2)<<8 | uint(b1)<<16 | uint(b0)<<24)
}
