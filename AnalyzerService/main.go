package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/AnalyzerService/app"
)

func main() {
	// This delay is here to give the database time to start up
	time.Sleep(10 * time.Second)

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	sessionStoreKey := os.Getenv("SessionsStoreKey")
	db, err := app.InitPostgresDB()
	if err != nil {
		log.Fatalf("Error setting up database: %s\n", err)
	}
	s := &app.Server{
		Router:                 http.DefaultServeMux,
		TargetForLyricsService: os.Getenv("LyricsServiceName") + ":" + os.Getenv("LyricsServicePort"),
		SessionStore:           sessions.NewCookieStore([]byte(sessionStoreKey)),
		CookieName:             "cookie-store-spotifyAnalzer",
		Database:               db,
	}
	err = s.Init()
	if err != nil {
		log.Fatalln(err)
	}

	servicePort := os.Getenv("AnalyzerServicePort")
	err = http.ListenAndServe(fmt.Sprintf(":%s", servicePort), s)
	return err
}
