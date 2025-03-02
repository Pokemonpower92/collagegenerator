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
	"github.com/pokemonpower92/collagegenerator/internal/imageprocessing"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/service"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
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
		Name:        "SeedSet",
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
		// Three shades of red
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Bright red
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 192, G: 0, B: 0, A: 255}, // Medium red
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 128, G: 0, B: 0, A: 255}, // Dark red
		},
		// Three shades of green
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 255, B: 0, A: 255}, // Bright green
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 192, B: 0, A: 255}, // Medium green
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 128, B: 0, A: 255}, // Dark green
		},
		// Three shades of blue
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 0, B: 255, A: 255}, // Bright blue
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 0, B: 192, A: 255}, // Medium blue
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 0, B: 128, A: 255}, // Dark blue
		},
		// Three grayscale values
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 128, G: 128, B: 128, A: 255}, // Medium gray
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 0, G: 0, B: 0, A: 255}, // Black
		},
		// Flesh tones
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 224, B: 196, A: 255}, // Light flesh tone
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 238, G: 207, B: 180, A: 255}, // Medium light flesh tone
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 224, G: 172, B: 105, A: 255}, // Medium flesh tone
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 141, G: 85, B: 36, A: 255}, // Dark flesh tone
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 198, G: 134, B: 66, A: 255}, // Medium dark flesh tone
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 89, G: 47, B: 42, A: 255}, // Deep flesh tone
		},
		// Purples
		{
			id:    uuid.New(),
			color: color.RGBA{R: 147, G: 112, B: 219, A: 255}, // Medium purple
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 75, G: 0, B: 130, A: 255}, // Indigo
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 186, G: 85, B: 211, A: 255}, // Medium orchid
		},
		// Yellows
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 255, B: 0, A: 255}, // Bright yellow
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 215, B: 0, A: 255}, // Gold
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 218, G: 165, B: 32, A: 255}, // Goldenrod
		},
		// Pinks
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 192, B: 203, A: 255}, // Pink
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 255, G: 20, B: 147, A: 255}, // Deep pink
		},
		{
			id:    uuid.New(),
			color: color.RGBA{R: 219, G: 112, B: 147, A: 255}, // Pale violet red
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
		average := imageprocessing.CalculateAverageColor(imageSetImage)
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
	tiRepo, err := repository.NewTargetImageRepository(c, ctx)
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
	service.CreateCollageMetaData(col)
	if err != nil {
		panic(err)
	}
}
