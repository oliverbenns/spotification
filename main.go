package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/oliverbenns/spotification/browser"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

var redirectUrl = "http://localhost:3000/callback"
var state = "abcd"

func main() {
	auth := spotify.NewAuthenticator(redirectUrl, spotify.ScopeUserLibraryModify, spotify.ScopeUserLibraryRead)

	url := auth.AuthURL(state)

	browser.Open(url)

	log.Print(url)

	mux := http.NewServeMux()

	callbackHandler := createCallbackHandler(auth)
	mux.HandleFunc("/callback", callbackHandler)
	err := http.ListenAndServe("localhost:3000", mux)

	if err != nil {
		panic(err)
	}
}

func createCallbackHandler(auth spotify.Authenticator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("WORKING")

		token, err := auth.Token(state, r)
		if err != nil {
			log.Print(err)
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}

		log.Print("HERE", token)

		client := auth.NewClient(token)

		trackPage, _ := client.CurrentUsersTracks()

		for _, track := range trackPage.Tracks {
			log.Print(track.FullTrack)
		}
	}
}
