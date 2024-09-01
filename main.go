package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2/clientcredentials"
)

var albums = make(map[string]spotify.SimpleAlbum)
var songTotal = 0
var ctx = context.Background()
var client *spotify.Client = nil

const userId = "31v3df7z4vnup4glne65eubncbbi"

func visitPLaylist(playlistId spotify.ID) {
	for i := 0; i < 1000; i++ {
		tracks, _ := client.GetPlaylistItems(ctx, playlistId, spotify.Limit(100), spotify.Offset(i*100))
		for j, track := range tracks.Items {
			simpletrack := track.Track.Track.SimpleTrack
			album := track.Track.Track.Album
			fmt.Println(i*100+j, "\t", simpletrack.Name, "\t", simpletrack.Artists[0].Name, "\t", album.Name)
			albums[album.Name] = album
		}
		if len(tracks.Items) == 0 {
			break
		}
	}
}

func visitUserPlaylists(userId string) {
	playlists, err := client.GetPlaylistsForUser(ctx, userId, spotify.Limit(100))
	if err != nil {
		panic(err)
	}
	if playlists.Total >= 100 {
		panic("too many playlists")
	}
	for i, playlist := range playlists.Playlists {
		fmt.Println(i, "/", playlists.Total, playlist.Name)
		visitPLaylist(playlist.ID)
	}
}

func main() {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)

	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}
	httpClient := spotifyauth.New().Client(ctx, token)
	client = spotify.New(httpClient)
	visitUserPlaylists(userId)

	fmt.Println(len(albums))
	f, err := os.OpenFile("albums", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, album := range albums {
		album_url := "https://open.spotify.com/album/" + album.ID
		fmt.Println(songTotal, album.Name, album_url)
		if _, err = f.WriteString(string(album_url) + "\n"); err != nil {
			panic(err)
		}
	}

}
