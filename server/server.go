package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/go-kom/model"
	"github.com/opchaves/go-kom/services"
)

type Server struct {
	*chi.Mux

	name string
	host string
	port string

	Logger   *httplog.Logger
	DB       *pgxpool.Pool
	Q        *model.Queries
	Services *services.Services
}

func (s *Server) Run() error {
	host := fmt.Sprintf("%s:%s", s.host, s.port)
	s.Logger.Debug(fmt.Sprintf("Server is running on %s", host))

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
