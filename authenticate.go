// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("log in to Spotify your browser:", url)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
//https://open.spotify.com/track/
	search, err := client.Search("queen", spotify.SearchTypeTrack)
	fmt.Print(search)
	var newSong = spotify.ID("68FhagAoZr9Ld8oCp9JoYP")
	fmt.Println("You are logged in as:", user.ID)
	tracks, err := client.GetPlaylistTracks("5oMgA72pUi1qtyqSZ0KVKX")
	playlist, err := client.GetPlaylist("5oMgA72pUi1qtyqSZ0KVKX")
	fmt.Print(tracks)
	fmt.Print(playlist)
	addPlaylist := true
	for _, track := range tracks.Tracks {
		if track.Track.SimpleTrack.ID == newSong {
			addPlaylist = false
		}
	}
	if addPlaylist {
		playlist, _ := client.AddTracksToPlaylist("5oMgA72pUi1qtyqSZ0KVKX", newSong)
		fmt.Print(playlist)
	}
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}
