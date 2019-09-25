package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/AnalyzerService/app"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	sessionStoreKey := os.Getenv("SessionsStoreKey")
	s := &app.Server{
		Router:                 http.DefaultServeMux,
		TargetForLyricsService: os.Getenv("LyricsServiceName") + ":" + os.Getenv("LyricsServicePort"),
		SessionStore:           sessions.NewCookieStore([]byte(sessionStoreKey)),
		CookieName:             "cookie-store-spotifyAnalzer",
	}
	err := s.Init()
	if err != nil {
		log.Fatalln(err)
	}

	servicePort := os.Getenv("AnalyzerServicePort")
	err = http.ListenAndServe(fmt.Sprintf(":%s", servicePort), s)
	return err
}
