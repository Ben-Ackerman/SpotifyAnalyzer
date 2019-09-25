package app

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
)

type Server struct {
	Router                 *http.ServeMux
	TargetForLyricsService string
	stopWordList           map[string]bool
	SessionStore           *sessions.CookieStore
	CookieName             string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// handles the logic on what to do when the user first enters the site.
// not our NGINX server does not redirect for root so this function is only called
// when we are testing this service locally without a proxy
func (s *Server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("src/localTesting.html")
		if err != nil {
			log.Println(err.Error())
		}
		temp.Execute(w, nil)
	}
}

func (s *Server) handleGetLyricsWordCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.SessionStore.Get(r, s.CookieName)
		if err != nil {
			log.Printf("Error in handleSpotifyLoginCallback: %s\n", err)
		}

		if spotifyLogin, ok := session.Values["loggedInWithSpotify"].(bool); !ok || !spotifyLogin {
			http.Error(w, "User must first log in with spotify", http.StatusBadRequest)
		}

		tracks, ok := session.Values["usersTopTracks"].([]Track)
		if ok {
			var lyricsBuilder strings.Builder
			for _, val := range tracks {
				lyr := cleanLyrics(val.Lyrics, removeSectionHeaders, trimWhiteSpace)
				lyricsBuilder.WriteString(lyr)
			}
			wordCounts := getWordCounts(lyricsBuilder.String())
			wordCounts = s.removeStopWords(wordCounts)
			top20Words := getTopNWords(wordCounts, 20)

			var sb strings.Builder
			sb.WriteString("[")
			for i, key := range top20Words {
				if i != 0 {
					sb.WriteString(",")
				}
				sb.WriteString(fmt.Sprintf(`{"word":"%s", "count":%d}`, key, wordCounts[key]))
			}
			sb.WriteString("]")
			temp, err := template.ParseFiles("src/results.html")
			if err != nil {
				log.Fatalf(err.Error())
			}
			input := sb.String()
			if err := temp.Execute(w, template.JS(input)); err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

// Init calls necessary initialization for the server
func (s *Server) Init() error {
	err := s.InitStopWords()
	if err != nil {
		return err
	}

	s.Routes()
	gob.Register([]Track{})
	return nil
}

// InitStopWords initializes the stop words for the server
func (s *Server) InitStopWords() error {
	s.stopWordList = make(map[string]bool)
	file, err := os.Open("src/stopwords.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(text), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		s.stopWordList[line] = true
	}

	return nil
}

// removeStopWords filters out stop words from wordCounts and returns a new map
func (s *Server) removeStopWords(wordCounts map[string]int) map[string]int {
	newMap := make(map[string]int)
	for key := range wordCounts {
		_, ok := s.stopWordList[key]
		if !ok {
			newMap[key] = wordCounts[key]
		}
	}
	return newMap
}
