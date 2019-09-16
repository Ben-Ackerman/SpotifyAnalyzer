package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService/app"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService/geniusapi"
	"github.com/Ben-Ackerman/SpotifyAnalyzer/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	servicePort := os.Getenv("LyricsServicePort")
	geniusAccessToken := os.Getenv("GeniusAccessToken")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", servicePort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// create a server instance
	s := app.Server{
		GeniusClient: geniusapi.NewGeniusClient(nil, geniusAccessToken),
	}
	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the Lyrics service to the server
	api.RegisterLyricsServer(grpcServer, &s)

	reflection.Register(grpcServer)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}
