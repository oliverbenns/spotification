# Spotification

Migrate mp3 tracks stored on disk to Spotify. It does this by parsing the file name, searching the Spotify database and then adding them to a unique playlist.

This is different to _playing_ local mp3 files through the client which is what appears in most search results on this subject!!

I wanted to find the Spotify tracks because:

- It allows the tracks to be accessible on all devices
- Clears up my hard disk space
- Gives the ability to share tracks and playlists
- Have all the album artwork and meta data all set and good in the client

It also appears playing local files is quite a hidden feature so I wouldn't be surprised if support for this is eventually dropped.

## Getting Started

- Register a new app through [Spotify's developer dashboard](https://developer.spotify.com/dashboard/applications).
- Add the redirect URI `http://localhost:3000/callback` in the app settings.
- `cp .env.example .env` and fill in your app id and secret.
- Ensure Go 1.13 is installed.
- Run `go run main.go path/to/music`
