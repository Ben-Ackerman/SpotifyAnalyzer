package app

// Routes sets up the handlerFuncs for the server router
func (s *Server) Routes() {
	s.Router.HandleFunc("/spotify/login", s.handleLogin())
	s.Router.HandleFunc("/spotify/callback", s.handleSpotifyLoginCallback())
}
