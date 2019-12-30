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
var state = "abcd"

type App struct {
	Auth         spotify.Authenticator
	Client       spotify.Client
	MusicLibPath string
	Tracks       []musiclib.Track
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

	musiclib.GetTracks(app.MusicLibPath, &app.Tracks)

	app.Client = app.Auth.NewClient(token)
	app.Client.AutoRetry = true

	app.AddSpotifyTracks()
}

func (app *App) AddSpotifyTracks() {
	musiclib.GetTracks(app.MusicLibPath, &app.Tracks)

	// @TODO: error handling.
	user, _ := app.Client.CurrentUser()
	playlistName := "Spotification " + time.Now().Format(time.RFC3339)
	playlist, _ := app.Client.CreatePlaylistForUser(user.ID, playlistName, "", false)

	var trackIds []spotify.ID

	// @TODO: Add preventions of hitting rate limit.
	for _, track := range app.Tracks {
		// @NOTE: Flags like "artist:" don't work unless the artist is the _exact_ name.
		query := track.Name + " " + track.Artist
		searchResult, err := app.Client.Search(query, spotify.SearchTypeTrack)

		if err != nil {
			log.Print(err)
		}

		if len(searchResult.Tracks.Tracks) == 0 {
			log.Print("Cannot find song for track", track)
		} else {
			trackIds := append(trackIds, searchResult.Tracks.Tracks[0].ID)
		}
	}

	// @NOTE: Add the tracks in bulk to reduce # of queries and aid with hitting api rate limit.
	// app.Client.AddTracksToPlaylist(playlist.ID, searchResult.Tracks.Tracks[0].ID)

	log.Print("Total mp3 tracks: " + strconv.Itoa(len(app.Tracks)))
	log.Print("Added Spotify tracks: " + strconv.Itoa(len(trackIds)))
}
