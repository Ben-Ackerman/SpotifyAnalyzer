package geniusapi_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/LyricsService/geniusapi"
)

func TestNewClient(t *testing.T) {
	accessToken := "XSBawdJT3kZ0-0xZESIPVQf1weWj3mY53EYwPguSYlxUa3RysWHPb-9gJeyrCG3z"
	c := geniusapi.NewGeniusClient(nil, accessToken)
	url, err := c.GetSongURL(context.Background(), "Kendrick Lamar", "DNA")
	if err != nil {
		t.Error(err)
	}

	lyrics, err := c.GetSongLyrics(context.Background(), url)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(lyrics, "got royalty inside my DNA") {
		t.Errorf("Incorrect lyrics")
	}
}
