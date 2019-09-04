package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService/app"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/genius"
	"google.golang.org/grpc"
)

var (
	geniusAccessToken = "XSBawdJT3kZ0-0xZESIPVQf1weWj3mY53EYwPguSYlxUa3RysWHPb-9gJeyrCG3z"
	port              = 7777
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// create a server instance
	s := app.Server{
		geniusClient = genius.NewGeniusClient(nil, geniusAccessToken),
	}
	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Lyrics service to the server
	api.RegisterLyricsServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}
}
