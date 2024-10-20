package collageapi

import (
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/handler"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/router"
	"github.com/pokemonpower92/collagegenerator/internal/server"
)

func Start() {
	config.LoadEnvironmentVariables()
	r := router.NewRouter()
	c := config.NewPostgresConfig()
	isRepo, err := repository.NewImageSetRepository(c)
	if err != nil {
		panic(err)
	}
	imageSetHandler := handler.NewImageSetHandler(isRepo)
	r.RegisterHandler("GET /images/sets", imageSetHandler.GetImageSets)
	r.RegisterHandler("GET /images/sets/{id}", imageSetHandler.GetImageSetById)

	tiRepo, err := repository.NewTagrgetImageRepository(c)
	targetImageHandler := handler.NewTargetImageHandler(tiRepo)
	r.RegisterHandler("GET /images/targets", targetImageHandler.GetTargetImages)
	r.RegisterHandler("GET /images/targets/{id}", targetImageHandler.GetTargetImageById)

	s := server.NewImageSetServer(r)
	s.Start()
}
