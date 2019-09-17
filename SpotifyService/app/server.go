package app

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
	"google.golang.org/grpc"
)

var (
	//TODO implement psuedo-random state generator and remove from global name space
	oauthStateString = "psuedo-random"
)

// Server is a struct used to represent a server while storing its dependenies along with implementing http.Handler
type Server struct {
	Router                 *http.ServeMux
	TargetForLyricsService string
	SpotifyAuth            spotifyapi.Authenticator
	stopWordList           map[string]bool
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// handleSpotifyCallback contains the logic on what needs to be done when the spotify api redirects back to our service
func (s *Server) handleSpotifyCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := s.SpotifyAuth.Token(oauthStateString, r)
		if err != nil {
			log.Fatalf(err.Error())
		}
		client := s.SpotifyAuth.NewClient(token)

		pagingTracks, err := client.GetUserTopTracks(50, spotifyapi.SpotifyTimeRangeLong)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if pagingTracks != nil {
			tracks, err := s.callLyricsService(pagingTracks)
			if err != nil {
				log.Fatal(err.Error())
			}

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

// callLyricsService makes the grpc call to our lyricsservice
func (s *Server) callLyricsService(p *spotifyapi.PagingTrack) ([]Track, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(s.TargetForLyricsService, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewLyricsClient(conn)

	apitracks := pagingToTracks(p)

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

// handleLogin contains the logic on what to perform when the user enters login
func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := s.SpotifyAuth.AuthCodeURL(oauthStateString, true)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
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
