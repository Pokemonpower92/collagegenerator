package database

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/service"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/utils"
)

func Seed() {
	config.LoadEnvironmentVariables()
	c := config.NewPostgresConfig()
	ctx := context.Background()
	isRepo, err := repository.NewImageSetRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer isRepo.Close()
	imSet, err := isRepo.Create(sqlc.CreateImageSetParams{
		Name:        uuid.NewString(),
		Description: "A seeded imageset",
	})
	if err != nil {
		panic(err)
	}
	store := datastore.NewStore()
	images := []struct {
		id    uuid.UUID
		color color.RGBA
	}{
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 200, G: 200, B: 200, A: 255},
		},
	}
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	for _, im := range images {
		draw.Draw(
			img,
			img.Bounds(),
			&image.Uniform{im.color},
			image.Point{},
			draw.Src,
		)
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			panic(err)
		}
		store.PutFile(im.id, &buf)
	}
	acRepo, err := repository.NewAverageColorRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer acRepo.Close()
	for _, im := range images {
		imageSetImage, err := store.GetRGBA(im.id)
		if err != nil {
			panic(err)
		}
		average := utils.CalculateAverageColor(imageSetImage)
		_, err = acRepo.Create(sqlc.CreateAverageColorParams{
			ID:         im.id,
			ImagesetID: imSet.ID,
			FileName:   im.id.String(),
			R:          int32(average.R),
			G:          int32(average.G),
			B:          int32(average.B),
			A:          int32(average.A),
		})
		if err != nil {
			panic(err)
		}
	}
	tiRepo, err := repository.NewTagrgetImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer tiRepo.Close()
	targetImage, err := tiRepo.Create(sqlc.CreateTargetImageParams{
		ID:          images[len(images)-1].id,
		Name:        "Grey",
		Description: "A seeded target image",
	})
	if err != nil {
		panic(err)
	}
	cRepo, err := repository.NewCollageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer cRepo.Close()
	col, err := cRepo.Create(sqlc.CreateCollageParams{
		Name:          "Seed Collage",
		Description:   "A seeded collage",
		ImageSetID:    imSet.ID,
		TargetImageID: targetImage.ID,
	})
	service.CreateCollage(col)
	if err != nil {
		panic(err)
	}
}
