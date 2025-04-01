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
	imageSetHandler := handler.NewImageSetHandler(isRepo, isRepo)
	r.RegisterRoute("POST /imagesets", imageSetHandler.CreateImageSet)
	r.RegisterRoute("GET /imagesets", imageSetHandler.GetImageSets)
	r.RegisterRoute("GET /imagesets/{id}", imageSetHandler.GetImageSetById)

	tiRepo, err := repository.NewTargetImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer tiRepo.Close()
	targetImageHandler := handler.NewTargetImageHandler(tiRepo, tiRepo)
	r.RegisterRoute("POST /targets", targetImageHandler.CreateTargetImage)
	r.RegisterRoute("GET /targets", targetImageHandler.GetTargetImages)
	r.RegisterRoute("GET /targets/{id}", targetImageHandler.GetTargetImageById)

	acRepo, err := repository.NewAverageColorRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer acRepo.Close()
	averageColorHandler := handler.NewAverageColorHandler(acRepo, acRepo)
	r.RegisterRoute("POST /averagecolors", averageColorHandler.CreateAverageColor)
	r.RegisterRoute("GET /averagecolors", averageColorHandler.GetAverageColors)
	r.RegisterRoute("GET /averagecolors/{id}", averageColorHandler.GetAverageColorById)
	r.RegisterRoute("GET /imagesetimages/{id}", averageColorHandler.GetByImageSetId)

	cRepo, err := repository.NewCollageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer cRepo.Close()
	collageHandler := handler.NewCollageHandler(cRepo, cRepo)
	r.RegisterRoute("POST /collages", collageHandler.CreateCollage)
	r.RegisterRoute("GET /collages", collageHandler.GetCollages)
	r.RegisterRoute("GET /collages/{id}", collageHandler.GetCollageById)

	ciRepo, err := repository.NewCollageImgageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer ciRepo.Close()
	collageImageHandler := handler.NewCollageImageHandler(ciRepo, ciRepo)
	r.RegisterRoute("POST /collageimages", collageImageHandler.CreateCollageImage)
	r.RegisterRoute("GET /collageimages", collageImageHandler.GetCollageImages)
	r.RegisterRoute("GET /collageimages/{id}", collageImageHandler.GetCollageImageByCollageId)

	r.RegisterRoute("POST /files", handler.StoreFile)
	r.RegisterRoute("GET /files/{id}", handler.GetFileById)

	s := server.NewCollageServer(r)
	s.Start()
}
