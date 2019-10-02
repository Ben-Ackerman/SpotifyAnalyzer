package app

import (
	"context"
	"log"
	"net/http"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService/geniusapi"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
)

// Server represention of gRPC Server
type Server struct {
	GeniusClient *geniusapi.GeniusClient
	Router       *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// GetLyrics gets the lyrics for the input tracks
func (s *Server) GetLyrics(ctx context.Context, track *api.TracksInfo) (*api.LyricsInfo, error) {
	uri, err := s.GeniusClient.GetSongURL(ctx, track.GetArtist(), track.GetName())
	if err != nil {
		log.Printf("Error geting uri from genius.com: %s", err)
		uri = ""
	}

	lyrics, err := s.GeniusClient.GetSongLyrics(ctx, uri)
	if err != nil {
		log.Printf("Error geting lyrics from genius.com: %s", err)
		lyrics = ""
	}

	result := &api.LyricsInfo{
		GeniusURI: uri,
		Lyrics:    lyrics,
	}
	return result, nil
}
