package geniusapi

import (
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

type GeniusClient struct {
	AccessToken string
	client      *http.Client
}

func NewGeniusClient(geniusClient *http.Client, token string) *GeniusClient {
	if geniusClient == nil {
		geniusClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	c := &GeniusClient{AccessToken: token, client: geniusClient}
	return c
}

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

func (c *GeniusClient) GetSongURL(artist string, song string) (string, error) {
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

		if strings.EqualFold(re.ReplaceAllString(hitArtist, ""), re.ReplaceAllString(artist, "")) && strings.Contains(hit.Result.Url, "lyrics") {
			url = hit.Result.Url
			break
		}
	}
	if url == "" {
		return "", fmt.Errorf("No match found for artist: %s and song: %s\n search term = %s", artist, song, searchQuery)
	}

	return url, nil
}

func (c *GeniusClient) GetSongLyrics(songURL string) (string, error) {
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
