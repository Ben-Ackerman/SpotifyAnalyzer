package app

import (
	"context"
	"html/template"
	"log"
	"net/http"

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
	Router               *http.ServeMux
	PortForLyricsService string
	SpotifyAuth          spotifyapi.Authenticator
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

		tracks, err := client.GetUserTopTracks(50, spotifyapi.SpotifyTimeRangeShort)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if tracks != nil {
			temp, err := template.ParseFiles("web/results.html")
			if err != nil {
				log.Println(err.Error())
			}
			if err := temp.Execute(w, nil); err != nil {
				log.Fatalf(err.Error())
			}
			s.callLyricsService(tracks)
		}
	}
}

func (s *Server) callLyricsService(p *spotifyapi.PagingTrack) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(s.PortForLyricsService, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewLyricsClient(conn)

	tracks, err := pagingToTracks(p)
	if err != nil {
		log.Fatalf("Error when calling pagingToTracks: %s", err)
	}

	response, err := c.GetLyrics(context.Background(), tracks)
	if err != nil {
		log.Fatalf("Error when calling GetLyrics: %s", err)
	}

	for i := 0; i < len(response.TrackInfo); i++ {
		log.Printf("Response from server: %s\n\n", response.TrackInfo[i].GetLyrics())
	}
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
		temp, err := template.ParseFiles("web/login.html")
		if err != nil {
			log.Println(err.Error())
		}
		temp.Execute(w, nil)
	}
}
