package geniusapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	geniusBaseURL = "https://api.genius.com"
)

// GeniusClient struct contains the accessToken and client for calling genius.com's developer api
type GeniusClient struct {
	AccessToken string
	client      *http.Client
}

// NewGeniusClient is a constructor for initiating a GeniusClient
func NewGeniusClient(geniusClient *http.Client, token string) *GeniusClient {
	if geniusClient == nil {
		geniusClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	c := &GeniusClient{AccessToken: token, client: geniusClient}
	return c
}

// executeRequest sets the necessary headers and executes req.  It then validates the http.Request
func (c *GeniusClient) executeRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// SearchSong takes in a query which repsents the query to execute against genius.com's search restful API
func (c *GeniusClient) SearchSong(query string) (*Response, error) {
	url := geniusBaseURL + "/search"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Add("q", query)
	req.URL.RawQuery = queryParams.Encode()

	bytes, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	var searchResponse Response
	err = json.Unmarshal(bytes, &searchResponse)
	if err != nil {
		return nil, err
	}

	return &searchResponse, nil
}

// GetSongURL takes in an artist and a song and returns the URL of the corresponding track on genius.com.
// Not if no matching url is found using genius.com's search API then an error is thrown
func (c *GeniusClient) GetSongURL(ctx context.Context, artist string, song string) (string, error) {
	artist = strings.TrimSpace(artist)
	song = strings.TrimSpace(song)
	searchQuery := fmt.Sprintf("%s %s", artist, song)
	response, err := c.SearchSong(searchQuery)
	if err != nil {
		return "", err
	}

	url := ""
	for i := 0; i < len(response.Response.Hits); i++ {
		hit := response.Response.Hits[i]
		hitArtist := strings.TrimSpace(hit.Result.PrimaryArtist.Name)

		//Genius autocapitalizes every artist name some some people add a zero-width space to allow lower case
		//Need to remove zero witdth space
		//https://genius.com/discussions/295632-Some-sort-of-command-to-allow-lowercase-letters-stylization

		re, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			return "", err
		}

		if strings.EqualFold(re.ReplaceAllString(hitArtist, ""), re.ReplaceAllString(artist, "")) && strings.Contains(hit.Result.URL, "lyrics") {
			url = hit.Result.URL
			break
		}
	}
	if url == "" {
		return "", fmt.Errorf("No match found for artist: %s and song: %s\n search term = %s", artist, song, searchQuery)
	}

	return url, nil
}

// GetSongLyrics scapes the provided URL for lyrics
func (c *GeniusClient) GetSongLyrics(ctx context.Context, songURL string) (string, error) {
	// Request the HTML page.
	if len(songURL) < 0 {
		return "", fmt.Errorf("No url provided to GeniusClient.GetSongLyrics")
	}
	res, err := c.client.Get(songURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error getting lyrics: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	lyrics := doc.Find(".lyrics").Text()

	return lyrics, nil
}
