package app

import (
	"testing"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
)

func TestPagingToTracks(t *testing.T) {
	p := &spotifyapi.PagingTrack{
		Tracks: []spotifyapi.Track{
			spotifyapi.Track{
				SpotifyID: "id1",
			},
			spotifyapi.Track{
				SpotifyID: "id2",
			},
		},
	}
	r := spotifyPagingToTracks(p)
	if r == nil {
		t.Errorf("FAILED, result should not be nil")
	}
	if r[0].SpotifyID != "id1" {
		t.Errorf("FAILED, track 1 did not convert correctly")
	}
	if r[1].SpotifyID != "id2" {
		t.Errorf("FAILED, track 2 did not convert correctly")
	}
}
