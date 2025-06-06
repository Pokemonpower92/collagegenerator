package authapi

import (
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	r := router.NewRouter()

	r.RegisterRoute("POST /authenticate", handler.Authenticate)
	r.RegisterRoute("POST /authorize", handler.Authorize)

	serverConfig := config.NewServerConfig()
	s := server.NewAuthServer(r, serverConfig)
	s.Start()
}
