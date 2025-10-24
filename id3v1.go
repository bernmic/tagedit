package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/text/encoding/charmap"
)

type V1Tag struct {
	Title   string
	Artist  string
	Album   string
	Year    string
	Comment string
	Genre   string
	Track   string
}

const (
	tagLength = 128
)

var _genres = map[byte]string{
	0:   "blues",
	1:   "classic_rock",
	2:   "country",
	3:   "dance",
	4:   "disco",
	5:   "funk",
	6:   "grunge",
	7:   "hip_hop",
	8:   "jazz",
	9:   "metal",
	10:  "new_age",
	11:  "oldies",
	12:  "other",
	13:  "pop",
	14:  "rnb",
	15:  "rap",
	16:  "reggae",
	17:  "rock",
	18:  "techno",
	19:  "industrial",
	20:  "alternative",
	21:  "ska",
	22:  "death_metal",
	23:  "pranks",
	24:  "soundtrack",
	25:  "euro_techno",
	26:  "ambient",
	27:  "trip_hop",
	28:  "vocal",
	29:  "jazz_funk",
	30:  "fusion",
	31:  "trance",
	32:  "classical",
	33:  "instrumental",
	34:  "acid",
	35:  "house",
	36:  "game",
	37:  "sound_clip",
	38:  "gospel",
	39:  "noise",
	40:  "alternrock",
	41:  "bass",
	42:  "soul",
	43:  "punk",
	44:  "space",
	45:  "meditative",
	46:  "instrumental_pop",
	47:  "instrumental_rock",
	48:  "ethnic",
	49:  "gothic",
	50:  "darkwave",
	51:  "techno_industrial",
	52:  "electronic",
	53:  "pop_folk",
	54:  "eurodance",
	55:  "dream",
	56:  "southern_rock",
	57:  "comedy",
	58:  "cult",
	59:  "gangsta",
	60:  "top_40",
	61:  "christian_rap",
	62:  "pop_funk",
	63:  "jungle",
	64:  "native_american",
	65:  "cabaret",
	66:  "new_wave",
	67:  "psychadelic",
	68:  "rave",
	69:  "showtunes",
	70:  "trailer",
	71:  "lo_fi",
	72:  "tribal",
	73:  "acid_punk",
	74:  "acid_jazz",
	75:  "polka",
	76:  "retro",
	77:  "musical",
	78:  "rock_n_roll",
	79:  "hard_rock",
	80:  "folk",
	81:  "folk_rock",
	82:  "national_folk",
	83:  "swing",
	84:  "fast_fusion",
	85:  "bebob",
	86:  "latin",
	87:  "revival",
	88:  "celtic",
	89:  "bluegrass",
	90:  "avantgarde",
	91:  "gothic_rock",
	92:  "progressive_rock",
	93:  "psychedelic_rock",
	94:  "symphonic_rock",
	95:  "slow_rock",
	96:  "big_band",
	97:  "chorus",
	98:  "easy_listening",
	99:  "acoustic",
	100: "humour",
	101: "speech",
	102: "chanson",
	103: "opera",
	104: "chamber_music",
	105: "sonata",
	106: "symphony",
	107: "booty_bass",
	108: "primus",
	109: "porn_groove",
	110: "satire",
	111: "slow_jam",
	112: "club",
	113: "tango",
	114: "samba",
	115: "folklore",
	116: "ballad",
	117: "power_ballad",
	118: "rhythmic_soul",
	119: "freestyle",
	120: "duet",
	121: "punk_rock",
	122: "drum_solo",
	123: "a_capella",
	124: "euro_house",
	125: "dance_hall",
	126: "Goa",
	127: "Drum & Bass",
	128: "Club-House",
	129: "Hardcore	   	",
	130: "Terror",
	131: "Indie",
	132: "BritPop",
	133: "Afro-Punk",
	134: "Polsk Punk",
	135: "Beat",
	136: "Christian Gangsta Rap",
	137: "Heavy Metal",
	138: "Black Metal",
	139: "Crossover",
	140: "Contemporary Christian",
	141: "Christian Rock",
	142: "Merengue",
	143: "Salsa",
	144: "Thrash Metal",
	145: "Anime",
	146: "JPop",
	147: "Synthpop",
	148: "Abstract",
	149: "Art Rock",
	150: "Baroque",
	151: "Bhangra",
	152: "Big Beat",
	153: "Breakbeat",
	154: "Chillout",
	155: "Downtempo",
	156: "Dub",
	157: "EBM",
	158: "Eclectic",
	159: "Electro",
	160: "Electroclash",
	161: "Emo",
	162: "Experimental",
	163: "Garage",
	164: "Global",
	165: "IDM",
	166: "Illbient",
	167: "Industro-Goth",
	168: "Jam Band",
	169: "Krautrock",
	170: "Leftfield",
	171: "Lounge",
	172: "Math Rock",
	173: "New Romantic",
	174: "Nu-Breakz",
	175: "Post-Punk",
	176: "Post-Rock",
	177: "Psytrance",
	178: "Shoegaze",
	179: "Space Rock",
	180: "Trop Rock",
	181: "World Music",
	182: "Neoclassical",
	183: "Audiobook",
	184: "Audio Theatre",
	185: "Neue Deutsche Welle",
	186: "Podcast",
	187: "Indie Rock",
	188: "G-Funk",
	189: "Dubstep",
	190: "Garage Rock",
	191: "Psybient",
	255: "None ",
}

// Open reads ID3V1 tag from file and returns all found data or nil and an error
func Open(file string) (*V1Tag, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			l(SEVERITY_ERROR, fmt.Sprintf("error closing file for id2v1: %v", err))
		}
	}()
	return ReadID3V1(f)
}

func ReadID3V1(f *os.File) (*V1Tag, error) {
	tag := V1Tag{}

	stats, err := f.Stat()
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("error stat file for hasId2v1: %v", err))
	}
	_, err = f.Seek(stats.Size()-tagLength, 0)
	if err != nil {
		return nil, err
	}
	data := make([]byte, tagLength)
	n, err := f.Read(data)
	if err != nil {
		return nil, err
	}
	if n != tagLength {
		return nil, fmt.Errorf("not enough data for a id3v1 tag: %d", n)
	}
	if string(data[:3]) == "TAG" {
		tag.Title = zstring(data[3:33])
		tag.Artist = zstring(data[33:63])
		tag.Album = zstring(data[63:93])
		tag.Year = zstring(data[93:97])
		tag.Comment = zstring(data[97:127])
		tag.Genre = _genres[data[127]]
		if data[125] == 0 && data[126] != 0 {
			//id3v1.1
			tag.Track = strconv.Itoa(int(data[126]))
		}
	}
	return &tag, nil
}

// HasID3V1 checks if a file has an ID3V1 header
func HasID3V1(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	defer func() {
		err := f.Close()
		if err != nil {
			l(SEVERITY_ERROR, fmt.Sprintf("error closing file for hasId2v1: %v", err))
		}
	}()

	stats, err := f.Stat()
	if err != nil {
		l(SEVERITY_ERROR, fmt.Sprintf("error stat file for hasId2v1: %v", err))
	}
	_, err = f.Seek(stats.Size()-tagLength, 0)
	if err != nil {
		return false
	}
	data := make([]byte, tagLength)
	n, err := f.Read(data)
	if err != nil {
		return false
	}
	if n != tagLength {
		return false
	}
	return string(data[:3]) == "TAG"
}

// RemoveID3V1 removes an ID3V1 header if exists
func RemoveID3V1(file string, removeTmp bool) error {
	if HasID3V1(file) {
		s, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		tmp := findFreeTmp(file)
		if err := os.Rename(file, tmp); err != nil {
			return err
		}
		err = os.WriteFile(file, s[0:len(s)-tagLength], 0644)
		if err != nil {
			return err
		}
		if removeTmp {
			err = os.Remove(tmp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func decodeISO8859_1ToUTF8(bytes []byte) string {
	encoded, _ := charmap.ISO8859_1.NewDecoder().Bytes(bytes)
	return string(encoded[:])
}

func zstring(n []byte) string {
	return decodeISO8859_1ToUTF8(n[:clen(n)])
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func findFreeTmp(file string) string {
	for i := 1; i < 1000; i++ {
		tmp := fmt.Sprintf("%s.tmp.(%d)", file, i)
		if _, err := os.Stat(tmp); errors.Is(err, os.ErrNotExist) {
			return tmp
		}
	}
	return ""
}
