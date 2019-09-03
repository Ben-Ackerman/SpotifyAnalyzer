package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/Ben-Ackerman/api"
	"github.com/Ben-Ackerman/spotifyapi"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	redirectURL = "http://127.0.0.1:5000/spotify/callback"
	//TODO move clientID and secret into OS vars
	clientID          = "e6dc8124143d4e72ae40b06ea162c98d"
	clientSecret      = "2f1255cfd38b481c9730fd3477cfa3d9"
	port              = ":5000"
	lyricsServicePort = ":7777"
)

var (
	//TODO implement psuedo-random state generator
	oauthStateString = "psuedo-random"
	a                = spotifyapi.NewAuthenticator(redirectURL, clientID, clientSecret, spotifyapi.ScopeUserTopRead)
)

func main() {
	fmt.Println("Server started on http://127.0.0.1:5000")
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleSpotifyLogin)
	http.HandleFunc("/spotify/callback", handleSpotifyCallback)
	log.Fatal(http.ListenAndServe(port, nil))
}

func callLyricsService(p *spotifyapi.PagingTrack) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(lyricsServicePort, grpc.WithInsecure())
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

func handleMain(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("web/login.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	temp.Execute(w, nil)
}

func handleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := a.AuthCodeURL(oauthStateString, true)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
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

func handleSpotifyCallback(w http.ResponseWriter, r *http.Request) {
	token, err := a.Token(oauthStateString, r)
	if err != nil {
		log.Fatalf(err.Error())
	}
	client := a.NewClient(token)

	tracks, err := client.GetUserTopTracks(50, spotifyapi.SpotifyTimeRangeShort)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if tracks != nil {
		fmt.Fprintf(w, "<html><body><p>Results found</p></body></html>")
		callLyricsService(tracks)
	}
}
