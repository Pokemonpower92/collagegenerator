package database

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/client"
	"github.com/pokemonpower92/collagegenerator/internal/imageprocessing"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/service"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/store"
)

type RandomImage struct {
	id    uuid.UUID
	image *image.RGBA
}

// Generates a 50x50 pixel image of a random, uniform color
// and stores it in the provided store.
func generateRandomImage(store store.Store) *RandomImage {
	id := uuid.New()
	color := color.RGBA{
		R: uint8(rand.Intn(255)),
		G: uint8(rand.Intn(255)),
		B: uint8(rand.Intn(255)),
		A: 255,
	}
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	draw.Draw(
		img,
		img.Bounds(),
		&image.Uniform{color},
		image.Point{},
		draw.Src,
	)
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	err := store.PutFile(id, &buf)
	if err != nil {
		panic(err)
	}
	return &RandomImage{
		id:    id,
		image: img,
	}
}

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
		Name:        "SeedSet",
		Description: "A seeded imageset",
	})
	if err != nil {
		panic(err)
	}

	acRepo, err := repository.NewAverageColorRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer acRepo.Close()

	// Generates 100 random images for the seed image set
	// and creates the average color records for them.
	store := client.NewFileClient("http://filestore:8081/files", slog.Default())
	for range 100 {
		randomImage := generateRandomImage(store)
		average := imageprocessing.CalculateAverageColor(randomImage.image)
		_, err = acRepo.Create(sqlc.CreateAverageColorParams{
			ID:         randomImage.id,
			ImagesetID: imSet.ID,
			FileName:   randomImage.id.String(),
			R:          int32(average.R),
			G:          int32(average.G),
			B:          int32(average.B),
			A:          int32(average.A),
		})
		if err != nil {
			panic(err)
		}
	}

	tiRepo, err := repository.NewTargetImageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer tiRepo.Close()

	randomTargetImage := generateRandomImage(store)
	targetImage, err := tiRepo.Create(sqlc.CreateTargetImageParams{
		ID:          randomTargetImage.id,
		Name:        "SeedTargetImage",
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
	if err != nil {
		panic(err)
	}

	service.CreateCollageMetaData(col, slog.Default())

	ciRepo, err := repository.NewCollageImgageRepository(c, ctx)
	if err != nil {
		panic(err)
	}
	defer ciRepo.Close()
	ci, err := ciRepo.Create(col.ID)
	if err != nil {
		panic(err)
	}
	service.GenerateCollage(ci, slog.Default())
}
