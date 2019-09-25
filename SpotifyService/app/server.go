package app

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
	"github.com/gorilla/sessions"
)

var (
	//TODO implement psuedo-random state generator and remove from global name space
	oauthStateString = "psuedo-random"
)

// Track is a stuct used to store meta data about a given track
type Track struct {
	Artist    string
	Name      string
	GeniusURI string
	Lyrics    string
}

// Server is a struct used to represent a server while storing its dependenies along with implementing http.Handler
type Server struct {
	Router       *http.ServeMux
	SpotifyAuth  spotifyapi.Authenticator
	SessionStore *sessions.CookieStore
	CookieName   string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// handleSpotifyCallback contains the logic on what needs to be done when the spotify api redirects back to our service
func (s *Server) handleSpotifyLoginCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.SessionStore.Get(r, s.CookieName)
		if err != nil {
			log.Printf("Error in handleSpotifyLoginCallback: %s\n", err)
		}

		if spotifyLogin, ok := session.Values["loggedInWithSpotify"].(bool); ok && spotifyLogin {
			// No need to login twice
			// Prevents errors when user back-clicks back to this page
			return
		}

		token, err := s.SpotifyAuth.Token(oauthStateString, r)
		if err != nil {
			log.Printf("Error in spotify callback: %s", err.Error())
			return
		}
		client := s.SpotifyAuth.NewClient(token)

		pagingTracks, err := client.GetUserTopTracks(50, spotifyapi.SpotifyTimeRangeLong)
		if err != nil {
			log.Printf("Error in spotify callback getting top 50 tracks: %s", err.Error())
			return
		}

		var tracks []Track
		if pagingTracks != nil {
			tracks = spotifyPagingToTracks(pagingTracks)
		}

		if tracks != nil {
			//tracks, err = s.callLyricsService(tracks)
			for i := 0; i < len(tracks); i++ {
				tracks[i].Lyrics = "test"
			}
			if err != nil {
				log.Printf("Error calling lyric service %s", err.Error())
				return
			}
		}

		session.Values["loggedInWithSpotify"] = true
		session.Values["usersTopTracks"] = tracks
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//TODO remove
		fmt.Println(w, "loggin success")
	}
}

// handleLogin contains the logic on what to perform when the user enters login
func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := s.SpotifyAuth.AuthCodeURL(oauthStateString, true)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// Init calls necessary initialization for the server
func (s *Server) Init() {
	s.Routes()
	gob.Register([]Track{})
}

// spotifyPagingToTracks takes in a spotifyapi PagingTrack struct and creates the corresponding slice of Track structs.
func spotifyPagingToTracks(p *spotifyapi.PagingTrack) []Track {
	length := len(p.Tracks)

	tracks := make([]Track, length)
	for i := 0; i < length; i++ {
		tracks[i].Name = p.Tracks[i].Name
		tracks[i].Artist = p.Tracks[i].Artists[0].Name
	}
	return tracks
}
