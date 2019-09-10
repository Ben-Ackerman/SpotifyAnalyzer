package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/app"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
)

const (
	redirectURL = "http://127.0.0.1:5000/spotify/callback"
	//TODO move clientID and secret into OS vars
	clientID     = "e6dc8124143d4e72ae40b06ea162c98d"
	clientSecret = "2f1255cfd38b481c9730fd3477cfa3d9"
	port         = 5000
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	fmt.Println("Server started on http://127.0.0.1:5000")
	s := &app.Server{
		Router:               http.DefaultServeMux,
		PortForLyricsService: ":7777",
		SpotifyAuth:          spotifyapi.NewAuthenticator(redirectURL, clientID, clientSecret, spotifyapi.ScopeUserTopRead),
	}
	s.Routes()
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s)
	return err
}
