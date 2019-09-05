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
func (s *Server) GetLyrics(ctx context.Context, in *api.Tracks) (*api.Tracks, error) {
	tracks := in.GetTrackInfo()
	for i := 0; i < len(tracks); i++ {
		//TODO do not search for urls if we already have it
		if len(tracks[i].GetGeniusURI()) == 0 {
			uri, err := s.GeniusClient.GetSongURL(tracks[i].GetArtist(), tracks[i].GetName())
			if err != nil {
				return nil, nil
			}
			tracks[i].GeniusURI = uri
		}

		lyrics, err := s.GeniusClient.GetSongLyrics(tracks[i].GetGeniusURI())
		if err != nil {
			tracks[i].Lyrics = ""
			log.Println(err.Error())
		} else {
			tracks[i].Lyrics = lyrics
		}
	}

	in.TrackInfo = tracks
	return in, nil
}
