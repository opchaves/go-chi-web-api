package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	*chi.Mux

	name string
	host string
	port string
}

func (s *Server) Run() error {
	host := fmt.Sprintf("%s:%s", s.host, s.port)
	fmt.Printf("Server %s is running on %s\n", s.name, host)

	return http.ListenAndServe(host, s)
}

func New(name string, options ...Option) *Server {
	s := &Server{
		Mux:  chi.NewRouter(),
		name: name,
		host: "0.0.0.0",
		port: "8080",
	}

	for _, fn := range options {
		fn(s)
	}

	return s
}
