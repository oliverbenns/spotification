package musiclib

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

type Track struct {
	Artist string
	Name   string
}

func GetTracks(path string, tracks *[]Track) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			GetTracks(path+"/"+file.Name(), tracks)
			continue
		}

		fileName := file.Name()

		if isMusicTrack(fileName) {
			track, err := parseMusicTrack(fileName)

			if err != nil {
				log.Print(err, fileName)
				continue
			}

			*tracks = append(*tracks, track)
		}

	}
}

func isMusicTrack(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))

	return strings.HasSuffix(ext, ".mp3") ||
		strings.HasSuffix(ext, ".flac") ||
		strings.HasSuffix(ext, ".wav")

}

// @NOTE: Song must be "[artist] - [song]" e.g. "2Pac - Changes"
// If any other delimiters are found, they are ignored and will be part of the track name.
// This might mean that those without the standard naming convention might have additional
// info (e.g. album) as part of the song name.
func parseMusicTrack(fileName string) (Track, error) {
	out := strings.SplitN(fileName, " - ", 2)
	if len(out) != 2 {
		return Track{}, errors.New("Cannot parse track")
	}

	return Track{
		Artist: out[0],
		Name:   strings.TrimSuffix(out[1], filepath.Ext(out[1])),
	}, nil
}
