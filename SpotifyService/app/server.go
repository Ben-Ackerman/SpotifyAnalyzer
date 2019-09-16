package app

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
	"google.golang.org/grpc"
)

var (
	//TODO implement psuedo-random state generator and remove from global name space
	oauthStateString = "psuedo-random"
)

// Server is a struct used to represent a server and implement http.Handler
type Server struct {
	Router                 *http.ServeMux
	TargetForLyricsService string
	SpotifyAuth            spotifyapi.Authenticator
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

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

			var sb strings.Builder
			sb.WriteString("<p>")
			sb.WriteString(fmt.Sprintf("%d tracks returned\n", len(tracks)))
			for _, val := range tracks {
				sb.WriteString(fmt.Sprintf("artist = %s; track = %s\nlyrics:\n%s\n\n",
					val.Artist, val.Name, val.Lyrics))
			}

			var lyricsBuilder strings.Builder
			for _, val := range tracks {
				lyr := cleanLyrics(val.Lyrics, removeSectionHeaders, trimWhiteSpace)
				lyricsBuilder.WriteString(lyr)
			}
			wordCounts := getWordCounts(lyricsBuilder.String())
			for key, val := range wordCounts {
				sb.WriteString(fmt.Sprintf("%s = %d\n", key, val))
			}

			sb.WriteString("</p>")
			temp, err := template.ParseFiles("web/results.html")
			if err != nil {
				log.Fatalf(err.Error())
			}
			input := sb.String()
			input = strings.ReplaceAll(input, "\n", "<br>")
			if err := temp.Execute(w, template.HTML(input)); err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}

func (s *Server) callLyricsService(p *spotifyapi.PagingTrack) ([]Track, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(s.TargetForLyricsService, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewLyricsClient(conn)

	apitracks, err := pagingToTracks(p)
	if err != nil {
		return nil, fmt.Errorf("Error when calling pagingToTracks: %s", err)
	}

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

func pagingToTracks(p *spotifyapi.PagingTrack) (*api.Tracks, error) {
	length := len(p.Tracks)

	trackInfo := make([]*api.Tracks_TrackInfo, length)
	for i := 0; i < length; i++ {
		trackInfo[i] = &api.Tracks_TrackInfo{}
		trackInfo[i].Name = p.Tracks[i].Name
		trackInfo[i].Artist = p.Tracks[i].Artists[0].Name
	}
	tracks := &api.Tracks{}
	tracks.TrackInfo = trackInfo

	return tracks, nil
}

func (s *Server) handleSpotifyLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := s.SpotifyAuth.AuthCodeURL(oauthStateString, true)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		temp, err := template.ParseFiles("web/localTesting.html")
		if err != nil {
			log.Println(err.Error())
		}
		temp.Execute(w, nil)
	}
}
