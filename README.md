# Spotification

Migrate mp3 tracks stored on disk to Spotify. It does this by parsing the file name, searching the Spotify database and then adding them to a unique playlist.

This is different to _playing_ local mp3 files through the client which is what appears in most search results on this subject!!

This is not what I wanted because:

- I wanted to have my songs accessible on all devices
- Clear up my hard disk space
- Have the ability to share songs and playlists
- Have all the album artwork and meta data all set and good in the client.

It appears playing local files is also quite a hidden feature so I wouldn't be surprised if Spotify eventually drop support for this.

## Getting Started

- Register a new app through [Spotify's developer dashboard](https://developer.spotify.com/dashboard/applications).
- `cp .env.example .env` and fill in your app id and secret.
- Ensure Go 1.13 is installed.
- Run `go run main.go path/to/music`
