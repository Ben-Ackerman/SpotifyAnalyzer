package geniusapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	GeniusBaseURL      = "https://api.genius.com"
	requiredURLKeyword = "lyrics"
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
	url := GeniusBaseURL + "/search"
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
	searchQuery := fmt.Sprintf("%s %s", artist, song)
	response, err := c.SearchSong(searchQuery)
	if err != nil {
		return "", err
	}

	url := ""
	for i := 0; i < len(response.Response.Hits); i++ {
		hit := response.Response.Hits[i]
		hitArtist := hit.Result.PrimaryArtist.Name
		if strings.EqualFold(artist, hitArtist) {
			url = hit.Result.Url
			break
		}
	}

	if url == "" {
		return "", fmt.Errorf("No match found for artist: %s and song: %s", artist, song)
	} else if !strings.Contains(url, "lyrics") {
		return "", fmt.Errorf("URL: %s does not contain the keyword: %s", url, requiredURLKeyword)
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
