package app

import (
	"context"
	"log"
	"net/http"
	"sync"

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

	var waitgroup sync.WaitGroup
	for i := 0; i < len(tracks); i++ {
		waitgroup.Add(1)
		go func(j int) {
			if len(tracks[j].GetGeniusURI()) == 0 {
				uri, err := s.GeniusClient.GetSongURL(tracks[j].GetArtist(), tracks[j].GetName())
				if err != nil {
					tracks[j].GeniusURI = ""
					log.Println(err.Error())
				} else {
					tracks[j].GeniusURI = uri
				}
			}

			lyrics, err := s.GeniusClient.GetSongLyrics(tracks[j].GetGeniusURI())
			if err != nil {
				tracks[j].Lyrics = ""
				log.Println(err.Error())
			} else {
				tracks[j].Lyrics = lyrics
			}

			waitgroup.Done()
		}(i)
	}

	waitgroup.Wait()
	in.TrackInfo = tracks

	return in, nil
}
