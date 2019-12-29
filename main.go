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
	"time"
)

var redirectUrl = "http://localhost:3000/callback"
var state = "abcd"

type App struct {
	Auth         spotify.Authenticator
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
	log.Print("WORKING")

	token, err := app.Auth.Token(state, r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}

	musiclib.GetTracks(app.MusicLibPath, &app.Tracks)

	client := app.Auth.NewClient(token)

	// @TODO: error handling.
	user, _ := client.CurrentUser()
	playlistName := "Spotification " + time.Now().Format(time.RFC3339)
	playlist, _ := client.CreatePlaylistForUser(user.ID, playlistName, "", false)

	// @TODO: Add preventions of hitting rate limit.
	for _, track := range app.Tracks {
		// @NOTE: Flags like "artist:" don't work unless the artist is the _exact_ name.
		query := track.Name + " " + track.Artist
		searchResult, err := client.Search(query, spotify.SearchTypeTrack)

		if err != nil {
			log.Print(err)
		}

		if len(searchResult.Tracks.Tracks) == 0 {
			log.Print("Cannot find song for track", track)
		} else {
			// @TODO: Make array of ids and add in bulk of 100 (limit). This will help with api rate limit.
			client.AddTracksToPlaylist(playlist.ID, searchResult.Tracks.Tracks[0].ID)
		}
	}
}
