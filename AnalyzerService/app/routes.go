package app

// Routes sets up the handlerFuncs for the server router
func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleRoot())
	s.Router.HandleFunc("/v1/my/wordcount/short", s.handleGetLyricsWordCount("tracksShort", "usersTopTracksShort"))
	s.Router.HandleFunc("/v1/my/wordcount/medium", s.handleGetLyricsWordCount("tracksMedium", "usersTopTracksMed"))
	s.Router.HandleFunc("/v1/my/wordcount/long", s.handleGetLyricsWordCount("tracksLong", "usersTopTracksLong"))
}
