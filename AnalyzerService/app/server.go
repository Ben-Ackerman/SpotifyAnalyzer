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
	"sync"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
	"github.com/gorilla/sessions"
	"google.golang.org/grpc"
)

// Track is a stuct used to store meta data about a given track
type Track struct {
	ID        string
	Artist    string
	Name      string
	GeniusURI string
	Lyrics    string
	Rank      int
}

// Server represents an instance of a server
type Server struct {
	Router       *http.ServeMux
	stopWordList map[string]bool
	Database     Database

	// Session management dependences
	SessionStore *sessions.CookieStore
	CookieName   string

	// Lyric service dependences
	TargetForLyricsService string
	LyricServiceConn       *grpc.ClientConn
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
				var waitgroup sync.WaitGroup
				for i := 0; i < len(tracks); i++ {
					tracks[i].ID, err = s.Database.GetTrackID(tracks[i].Name, tracks[i].Artist)
					if err != nil {
						log.Printf("Error calling tracks database: %s", err)
						return
					}
					log.Printf("searching for Song name: %s found id of /%s/", tracks[i].Name, tracks[i].ID)
					if tracks[i].ID == "" {
						waitgroup.Add(1)
						go func(j int) {
							// Call services to populate info and then place in database
							lyrics, geniusURI, err := s.callLyricsService(r.Context(), tracks[j].Artist, tracks[j].Name)
							errFound := false
							if err != nil {
								log.Printf("Error calling lyric service: %s", err)
								errFound = true
							}

							if !errFound {
								tracks[j].Lyrics = lyrics
								tracks[j].GeniusURI = geniusURI

								// Insert populated track into database for use later if another user
								// need the same track info
								tracks[j].ID, err = s.Database.InsertTrack(&tracks[j])
								if err != nil {
									log.Printf("Error calling tracks database: %s", err)
									errFound = true
								}
								log.Printf("Inserting ID of %s\n", tracks[j].ID)
							}

							waitgroup.Done()
						}(i)
					} else {
						result, err := s.Database.GetTrack(tracks[i].ID)
						if err != nil {
							log.Printf("Error calling tracks database: %s", err)
							return
						}
						tracks[i].GeniusURI = result.GeniusURI
						tracks[i].Lyrics = result.Lyrics
					}
				}
				waitgroup.Wait()
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

// callLyricsService makes the grpc call to our lyricsservice and returns the lyrics and the geniusURI of the input track
func (s *Server) callLyricsService(ctx context.Context, artist string, name string) (string, string, error) {
	if artist == "" {
		return "", "", fmt.Errorf("Must provide an artist to callLyricsService")
	}
	if name == "" {
		return "", "", fmt.Errorf("Must provide an name to callLyricsService")
	}
	client := api.NewLyricsClient(s.LyricServiceConn)

	trackInfo := &api.TracksInfo{
		Name:   name,
		Artist: artist,
	}

	response, err := client.GetLyrics(ctx, trackInfo)
	if err != nil {
		return "", "", fmt.Errorf("Error when calling GetLyrics: %s", err)
	}

	return response.GetLyrics(), response.GetGeniusURI(), nil
}

// Close closes all necessary conncections when server stops
func (s *Server) Close() {
	s.Database.Close()
	s.LyricServiceConn.Close()
}

// Init calls necessary initialization for the server
func (s *Server) Init() error {
	if err := s.InitStopWords(); err != nil {
		return err
	}

	if err := s.InitLyricServiceConn(); err != nil {
		return err
	}

	s.Routes()
	gob.Register([]Track{})
	return nil
}

// InitLyricServiceConn initites a grpc client for the lyrics service
func (s *Server) InitLyricServiceConn() error {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(s.TargetForLyricsService, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %s", err)
	}
	s.LyricServiceConn = conn
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
