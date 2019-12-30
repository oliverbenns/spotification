package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/oliverbenns/spotification/browser"
	"github.com/oliverbenns/spotification/musiclib"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var redirectUrl = "http://localhost:3000/callback"
var state = ""

type App struct {
	Auth           spotify.Authenticator
	Client         spotify.Client
	MusicLibPath   string
	Playlist       *spotify.FullPlaylist
	RemoteTrackIds []spotify.ID
	Tracks         []musiclib.Track
}

func main() {
	relativePath := os.Args[1]
	path, _ := filepath.Abs(relativePath)

	app := App{
		Auth:         spotify.NewAuthenticator(redirectUrl, spotify.ScopeUserLibraryModify, spotify.ScopeUserLibraryRead, spotify.ScopePlaylistModifyPrivate),
		MusicLibPath: path,
	}

	url := app.Auth.AuthURL(state)

	browser.Open(url)

	mux := http.NewServeMux()

	mux.HandleFunc("/callback", app.CallbackHandler)
	err := http.ListenAndServe("localhost:3000", mux)

	if err != nil {
		panic(err)
	}
}

func (app *App) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := app.Auth.Token(state, r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	app.Client = app.Auth.NewClient(token)
	app.Client.AutoRetry = true

	musiclib.GetTracks(app.MusicLibPath, &app.Tracks)

	app.CreateSpotifyPlaylist()
	app.FindSpotifyTracks()
	app.AddSpotifyTracks()

	log.Print("Total mp3 tracks: " + strconv.Itoa(len(app.Tracks)))
	log.Print("Added Spotify tracks: " + strconv.Itoa(len(app.RemoteTrackIds)))
}

func (app *App) CreateSpotifyPlaylist() {
	user, _ := app.Client.CurrentUser()
	playlistName := "Spotification " + time.Now().Format(time.RFC3339)
	app.Playlist, _ = app.Client.CreatePlaylistForUser(user.ID, playlistName, "Playlist created by https://github.com/oliverbenns/spotification", false)
}

func (app *App) FindSpotifyTracks() {
	for _, track := range app.Tracks {
		// @NOTE: Flags like "artist:" don't work unless the artist is the _exact_ name.
		query := track.Name + " " + track.Artist
		searchResult, err := app.Client.Search(query, spotify.SearchTypeTrack)

		if err != nil {
			log.Print(err)
		}

		trackPrettyPrint := track.Artist + " - " + track.Name

		if len(searchResult.Tracks.Tracks) == 0 {
			log.Println("❌ " + trackPrettyPrint)
		} else {
			log.Println("✅ " + trackPrettyPrint)
			app.RemoteTrackIds = append(app.RemoteTrackIds, searchResult.Tracks.Tracks[0].ID)
		}
	}
}

// @NOTE: Add the tracks in bulk to reduce # of queries which helps reduce api rate limit hits.
func (app *App) AddSpotifyTracks() {
	totalCount := len(app.RemoteTrackIds)
	index := 0

	for index < totalCount {
		var increment int
		if totalCount-index >= 100 {
			increment = 100
		} else {
			increment = totalCount % 100
		}

		ids := app.RemoteTrackIds[index : index+increment]
		_, err := app.Client.AddTracksToPlaylist(app.Playlist.ID, ids...)

		if err != nil {
			log.Panic("Error adding ids to playlist", err)
		}

		index += increment
	}
}
