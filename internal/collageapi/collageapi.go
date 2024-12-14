package collageapi

import (
	"context"

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
	ctx := context.Background()

	isRepo, err := repository.NewImageSetRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer isRepo.Close()
	imageSetHandler := handler.NewImageSetHandler(isRepo)
	r.RegisterRoute("POST /imagesets", imageSetHandler.CreateImageSet)
	r.RegisterRoute("GET /imagesets", imageSetHandler.GetImageSets)
	r.RegisterRoute("GET /imagesets/{id}", imageSetHandler.GetImageSetById)

	tiRepo, err := repository.NewTagrgetImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer tiRepo.Close()
	targetImageHandler := handler.NewTargetImageHandler(tiRepo)
	r.RegisterRoute("POST /targets", targetImageHandler.CreateTargetImage)
	r.RegisterRoute("GET /targets", targetImageHandler.GetTargetImages)
	r.RegisterRoute("GET /targets/{id}", targetImageHandler.GetTargetImageById)

	acRepo, err := repository.NewAverageColorRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer acRepo.Close()
	averageColorHandler := handler.NewAverageColorHandler(acRepo)
	r.RegisterRoute("POST /averagecolors", averageColorHandler.CreateAverageColor)
	r.RegisterRoute("GET /averagecolors", averageColorHandler.GetAverageColors)
	r.RegisterRoute("GET /averagecolors/{id}", averageColorHandler.GetAverageColorById)

	cRepo, err := repository.NewCollageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer cRepo.Close()
	collageHandler := handler.NewCollageHandler(cRepo)
	r.RegisterRoute("POST /collages", collageHandler.CreateCollage)
	r.RegisterRoute("GET /collages", collageHandler.GetCollages)
	r.RegisterRoute("GET /collages/{id}", collageHandler.GetCollageById)

	imageHandler := handler.NewImageHandler()
	r.RegisterRoute("POST /images", imageHandler.StoreImage)
	r.RegisterRoute("GET /images/{id}", imageHandler.GetImageById)

	s := server.NewCollageServer(r)
	s.Start()
}
