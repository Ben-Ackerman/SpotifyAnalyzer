package app

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
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

// callLyricsService makes the grpc call to our lyricsservice
func (s *Server) callLyricsService(t []Track) ([]Track, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(s.TargetForLyricsService, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewLyricsClient(conn)

	apitracks := tracksToAPITracks(t)

	response, err := c.GetLyrics(context.Background(), apitracks)
	if err != nil {
		return nil, fmt.Errorf("Error when calling GetLyrics: %s", err)
	}

	tracks := make([]Track, len(response.TrackInfo))
	for i := 0; i < len(response.TrackInfo); i++ {
		tracks[i].Lyrics = response.GetTrackInfo()[i].GetLyrics()
		tracks[i].Artist = response.GetTrackInfo()[i].GetArtist()
		tracks[i].Name = response.GetTrackInfo()[i].GetName()
	}

	return tracks, nil
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

		type wordToCount struct {
			Word  string `json:"word"`
			Count int    `json:"count"`
		}
		type results struct {
			Result []wordToCount `json:"result"`
		}
		rv := &results{}

		tracks, ok := session.Values["usersTopTracks"].([]Track)
		if ok {
			if tracks != nil {
				tracks, err = s.callLyricsService(tracks)
				if err != nil {
					log.Printf("Error calling lyric service %s", err.Error())
					return
				}
			}

			var lyricsBuilder strings.Builder
			for _, val := range tracks {
				lyr := cleanLyrics(val.Lyrics, removeSectionHeaders, trimWhiteSpace)
				lyricsBuilder.WriteString(lyr)
			}
			wordCounts := getWordCounts(lyricsBuilder.String())
			wordCounts = s.removeStopWords(wordCounts)
			top20Words := getTopNWords(wordCounts, 20)

			rv.Result = make([]wordToCount, len(top20Words))
			for i, key := range top20Words {
				rv.Result[i].Word = key
				rv.Result[i].Count = wordCounts[key]
			}
		}
		json.NewEncoder(w).Encode(rv)
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
