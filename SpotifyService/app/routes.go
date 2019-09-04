package app

//Routes sets up the handlerFuncs for the server router
func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.handleRoot())
	s.Router.HandleFunc("/login", s.handleSpotifyLogin())
	s.Router.HandleFunc("/spotify/callback", s.handleSpotifyCallback())
}
