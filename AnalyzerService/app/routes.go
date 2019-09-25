package app

// Routes sets up the handlerFuncs for the server router
func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleRoot())
	s.Router.HandleFunc("/v1/my/wordcount", s.handleGetLyricsWordCount())
}
