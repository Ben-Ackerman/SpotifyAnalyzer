package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/app"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	clientID := os.Getenv("SpotifyClientID")
	clientSecret := os.Getenv("SpotifyClientSecret")
	redirectURL := os.Getenv("SpotifyRedirectURL")
	s := &app.Server{
		Router:                 http.DefaultServeMux,
		TargetForLyricsService: os.Getenv("LyricsServiceName") + ":" + os.Getenv("LyricsServicePort"),
		SpotifyAuth:            spotifyapi.NewAuthenticator(redirectURL, clientID, clientSecret, spotifyapi.ScopeUserTopRead),
	}
	s.InitStopWords()
	s.Routes()

	servicePort := os.Getenv("SpotifyServicePort")
	err := http.ListenAndServe(fmt.Sprintf(":%s", servicePort), s)
	return err
}
