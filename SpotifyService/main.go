package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

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
	sessionStoreKey := os.Getenv("SessionsStoreKey")
	s := &app.Server{
		Router:       http.DefaultServeMux,
		SpotifyAuth:  spotifyapi.NewAuthenticator(redirectURL, clientID, clientSecret, spotifyapi.ScopeUserTopRead),
		SessionStore: sessions.NewCookieStore([]byte(sessionStoreKey)),
		CookieName:   "cookie-store-spotifyAnalzer",
	}
	s.Init()

	servicePort := os.Getenv("SpotifyServicePort")
	err := http.ListenAndServe(fmt.Sprintf(":%s", servicePort), s)
	return err
}
