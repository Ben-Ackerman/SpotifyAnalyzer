package spotifyapi

type paging struct {
	Endpoint    string `json:"href"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
	Total       int    `json:"total"`
	NextPageURL string `json:"next"`
	PrevPageURL string `json:"previous"`
}

// Artist stores spotify info about artists
// for more info on given fields see https://developer.spotify.com/documentation/web-api/reference/object-model/#artist-object-simplified
type Artist struct {
	Endpoint   string `json:"href"`
	SpotifyID  string `json:"id"`
	Name       string `json:"name"`
	SpotifyURI string `json:"uri"`
}

// Track stores spotify object info about tracks
// for more info on given fields see https://developer.spotify.com/documentation/web-api/reference/object-model/#track-object-simplified
type Track struct {
	Artists []Artist `json:"artists"`
	//Song duration in ms
	Duration    int    `json:"duration_ms"`
	Explicit    bool   `json:"explicit"`
	Endpoint    string `json:"href"`
	SpotifyID   string `json:"id"`
	Name        string `json:"name"`
	TrackNumver int    `json:"track_number"`
	SpotifyURI  string `json:"uri"`
}

// PagingArtist stores a spotify paginated object of artists objects
type PagingArtist struct {
	paging
	Artists []Artist `json:"items"`
}

// PagingTrack stores a spotify paginated object of track objects
type PagingTrack struct {
	paging
	Tracks []Track `json:"items"`
}
