package genius

type Response struct {
	Meta struct {
		Status int `json:"status"`
		Message string `json:"message"`
		
	} `json:"meta"`
	Response struct {
		Hits []*Hit `json:"hits"`
	} `json:"response"`
}

type Hit struct {
	Index string `json:"index"`
	Type string `json:"type"`
	Result *Song `json:"result"`
}

type Song struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Url string `json:"url"`
	PrimaryArtist *Artist `json:"primary_artist"`
}

type Artist struct {
	ID int `json:"id"`
	Name string `json:"name"`
}