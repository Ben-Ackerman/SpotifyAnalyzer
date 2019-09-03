package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Ben-Ackerman/api"
	"github.com/Ben-Ackerman/genius"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	geniusAccessToken = "XSBawdJT3kZ0-0xZESIPVQf1weWj3mY53EYwPguSYlxUa3RysWHPb-9gJeyrCG3z"
	port              = 7777
)

// Server representation of gRPC Server
type Server struct {
	geniusClient *genius.GeniusClient
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := Server{}
	s.geniusClient = genius.NewGeniusClient(nil, geniusAccessToken)
	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Lyrics service to the server
	api.RegisterLyricsServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// GetLyrics gets the lyrics for the input tracks
func (s *Server) GetLyrics(ctx context.Context, in *api.Tracks) (*api.Tracks, error) {
	tracks := in.GetTrackInfo()
	for i := 0; i < len(tracks); i++ {
		//TODO do not search for urls if we already have it
		lyrics, err := getSongLyrics(s.geniusClient, tracks[i].Artist, tracks[i].Name)
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

func getSongLyrics(client *genius.GeniusClient, artist string, song string) (string, error) {
	url, err := client.GetSongURL(artist, song)
	if err != nil {
		return "", err
	}

	lyrics, err := client.GetSongLyrics(url)
	if err != nil {
		return "", err
	}

	return lyrics, nil
}

func getSongLyricsFromURL(client *genius.GeniusClient, url string) (string, error) {
	lyrics, err := client.GetSongLyrics(url)
	if err != nil {
		return "", err
	}

	return lyrics, nil
}
