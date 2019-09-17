package geniusapi

// Response corresponds to genius.com's restful API response type.  This type is used to unmarshall the json response from genius.com.
type Response struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Response struct {
		Hits []*Hit `json:"hits"`
	} `json:"response"`
}

// Hit corresponds to genius.com's restful API hit type.  This type is used to unmarshall the json response from genius.com.
type Hit struct {
	Index  string `json:"index"`
	Type   string `json:"type"`
	Result *Song  `json:"result"`
}

// Song corresponds to genius.com's restful API track type.  This type is used to unmarshall the json response from genius.com.
type Song struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	URL           string  `json:"url"`
	PrimaryArtist *Artist `json:"primary_artist"`
}

// Artist corresponds to genius.com's restful API artist type.  This type is used to unmarshall the json response from genius.com.
type Artist struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
