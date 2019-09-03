package genius_test

import (
	"github.com/Ben-Ackerman/Genius_API"
	"testing"
	"strings"
)

func TestNewClient(t *testing.T) {
	accessToken := "XSBawdJT3kZ0-0xZESIPVQf1weWj3mY53EYwPguSYlxUa3RysWHPb-9gJeyrCG3z"
	c := genius.NewGeniusClient(nil, accessToken)
	url, err := c.GetSongURL("Kendrick Lamar", "DNA")
	if err != nil {
		t.Error(err)
	}

	lyrics, err := c.GetSongLyrics(url)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(lyrics, "got royalty inside my DNA") {
		t.Errorf("Incorrect lyrics")
	}
}

