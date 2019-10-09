package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"github.com/go-redis/redis"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
)

var (
	//TODO implement psuedo-random state generator and remove from global name space
	oauthStateString = "psuedo-random"
)

// Track is a stuct used to store meta data about a given track
type Track struct {
	SpotifyID string `json:"spotifyID"`
	//Rank      int    `json:"rank"`
}

// Tracks is a type to represent an array of tracks
type Tracks []Track

// Server is a struct used to represent a server while storing its dependenies along with implementing http.Handler
type Server struct {
	Router         *http.ServeMux
	SpotifyAuth    spotifyapi.Authenticator
	cookieName     string
	cookieDuration time.Duration
	RedisClient    *redis.Client
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// handleSpotifyCallback contains the logic on what needs to be done when the spotify api redirects back to our service
func (s *Server) handleSpotifyLoginCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Check cookies for a sessionID
		sessionID := ""
		for _, cookie := range r.Cookies() {
			if cookie.Name == s.cookieName && time.Now().Before(cookie.Expires) {
				sessionID = cookie.Value
			}
		}

		// If we already have a sessionID check to see if that ID is still stored in redis.
		// If it has evicted/expired from redis we no longer have the user data and need to
		// recalculate it.
		userDataFound := false
		if sessionID != "" {
			_, err := s.RedisClient.Get(sessionID).Result()
			if err == redis.Nil {
				log.Println("sessionID does not exist")
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("Error in spotify callback: %s", err.Error())
				return
			} else {
				userDataFound = true
			}
		}

		// If we already have a session for this user no need to recalcute their user data just send them back
		// to the login page.
		if userDataFound {
			temp, err := template.ParseFiles("src/spotifyRedirect.html")
			if err != nil {
				log.Println(err.Error())
			}
			temp.Execute(w, nil)
			return
		}

		// If we didn't find the users data in redis we need to calculate it.
		token, err := s.SpotifyAuth.Token(oauthStateString, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error in spotify callback: %s", err.Error())
			return
		}

		sessionID, err = generateSessionID()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error generating sessionID in handleSpotifyLoginCallback: %s", err)
			return
		}

		client := s.SpotifyAuth.NewClient(token)
		keyNames := []string{"tracksShort", "tracksMedium", "tracksLong"}
		timeFrames := []string{spotifyapi.SpotifyTimeRangeShort,
			spotifyapi.SpotifyTimeRangeMedium,
			spotifyapi.SpotifyTimeRangeLong}

		var waitgroup sync.WaitGroup
		for i := 0; i < len(keyNames); i++ {
			waitgroup.Add(1)
			go func(j int) {
				errFound := false
				pagingTracks, err := client.GetUserTopTracks(50, timeFrames[j])
				if err != nil {
					log.Printf("Error in spotify callback getting top 50 tracks: %s", err.Error())
					errFound = true
				}
				if !errFound {
					var tracks Tracks
					if pagingTracks != nil {
						tracks = spotifyPagingToTracks(pagingTracks)
					}

					var jsonData []byte
					jsonData, err := json.Marshal(tracks)
					if err != nil {
						log.Printf("Error marshalling data %s", err)
					}

					err = s.RedisClient.Set(fmt.Sprintf("%s-%s", sessionID, keyNames[j]), string(jsonData), s.cookieDuration).Err()
					if err != nil {
						log.Printf("Error setting redis value: %s", err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
				waitgroup.Done()
			}(i)
		}
		waitgroup.Wait()

		expiration := time.Now().Add(s.cookieDuration)
		cookie := http.Cookie{Name: s.cookieName, Value: sessionID, Expires: expiration}
		http.SetCookie(w, &cookie)

		temp, err := template.ParseFiles("src/spotifyRedirect.html")
		if err != nil {
			log.Printf("Error writing template: %s", err.Error())
		}
		temp.Execute(w, nil)
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

}

// addCookie dds a cookie to a http.response
func (s *Server) addSessionID(w http.ResponseWriter, sessionID string) {
	expire := time.Now().Add(s.cookieDuration)
	cookie := http.Cookie{
		Name:     s.cookieName,
		Value:    sessionID,
		Expires:  expire,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

// spotifyPagingToTracks takes in a spotifyapi PagingTrack struct and creates the corresponding Tracks structs.
func spotifyPagingToTracks(p *spotifyapi.PagingTrack) Tracks {
	length := len(p.Tracks)

	tracks := make([]Track, length)
	for i := 0; i < length; i++ {
		tracks[i].SpotifyID = p.Tracks[i].SpotifyID
		//tracks[i].Rank = i
	}
	return tracks
}

// Generates a random sessionID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
