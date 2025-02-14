package server

import (
	"log"
	"net/http"

	"github.com/pokemonpower92/collagegenerator/internal/router"
)

type Server struct {
	server *http.Server
	router *router.Router
}

func (s *Server) Start() {
	log.Printf("Starting server on %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func NewCollageServer(router *router.Router) *Server {
	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: router.Mux,
	}

	return &Server{
		server: server,
		router: router,
	}
}
