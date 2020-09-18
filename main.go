package main

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
)

func main() {
	const state = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	scopes := []string{spotify.ScopeUserLibraryRead, spotify.ScopeUserModifyPlaybackState,
		spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate}

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		Scopes:       scopes,
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	msg, page, err := client.FeaturedPlaylists()
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(msg)
	for _, playlist := range page.Playlists {
		fmt.Println("  ", playlist.Name)
	}

	user, err := client.GetUsersPublicProfile(spotify.ID("ihca9m4mwdfy1prbe83lirde7"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Println("User ID:", user.ID)
	fmt.Println("Display name:", user.DisplayName)
	fmt.Println("Spotify URI:", string(user.URI))
	fmt.Println("Endpoint:", user.Endpoint)
	fmt.Println("Followers:", user.Followers.Count)
	playlist, err := client.AddTracksToPlaylist("5oMgA72pUi1qtyqSZ0KVKX", "6ttsH99vfvkAPF3s1tIPqB")
	fmt.Print(playlist)
}
