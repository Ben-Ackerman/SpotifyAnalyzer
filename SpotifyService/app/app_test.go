package app

import (
	"testing"

	"github.com/Ben-Ackerman/SpotifyAnalyzer/SpotifyService/spotifyapi"
)

func TestPagingToTracks(t *testing.T) {
	p := &spotifyapi.PagingTrack{
		Tracks: []spotifyapi.Track{
			spotifyapi.Track{
				Artists: []spotifyapi.Artist{
					spotifyapi.Artist{
						Name: "name1",
					},
				},
				Name: "name1",
			},
			spotifyapi.Track{
				Artists: []spotifyapi.Artist{
					spotifyapi.Artist{
						Name: "name2",
					},
				},
				Name: "name2",
			},
		},
	}
	r := spotifyPagingToTracks(p)
	if r == nil {
		t.Errorf("FAILED, result should not be nil")
	}
	if r[0].Artist != "name1" || r[0].Name != "name1" {
		t.Errorf("FAILED, track 1 did not convert correctly")
	}
	if r[1].Artist != "name2" || r[1].Name != "name2" {
		t.Errorf("FAILED, track 2 did not convert correctly")
	}
}
