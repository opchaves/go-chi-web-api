package server

// Allows to specify options to the server.
type Option func(*Server)

func UsePort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

func UseHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}
