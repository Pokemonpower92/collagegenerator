package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	"github.com/nfnt/resize"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/client"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageGenerator struct {
	logger         *slog.Logger
	store          client.FileClient
	collageImageId uuid.UUID
	collageId      uuid.UUID
	imageCache     *lru.Cache
}

// GenerateCollage generates the final image for a collage.
func GenerateCollage(collageImage *sqlc.CollageImage, logger *slog.Logger) {
	store := client.NewFileClient("http://filestore:8081/files", logger.With())
	cache, err := lru.New(1000)
	if err != nil {
		panic(err)
	}
	generator := CollageGenerator{
		logger:         logger,
		store:          *store,
		collageImageId: collageImage.ID,
		collageId:      collageImage.CollageID,
		imageCache:     cache,
	}
	generator.generate()
}

func (cg *CollageGenerator) generate() {
	cg.logger.Info(fmt.Sprintf("Generating collage for collage %s\n", cg.collageId))
	// Get metadata file as collageMetaData
	metaData, err := cg.getMetaData()
	if err != nil {
		cg.logger.Error(fmt.Sprintf("Error getting collage metadata: %v\n", err))
		return
	}
	// Create blank canvas.
	canvas, err := cg.createCanvas(metaData.Resolution)
	if err != nil {
		cg.logger.Error(fmt.Sprintf("Error getting creating canvas: %v\n", err))
		return
	}
	// Build image concurrently.
	if err := cg.buildFinalImage(canvas, metaData); err != nil {
		cg.logger.Info(fmt.Sprintf("Error generating final image: %s\n", err))
	}
	cg.logger.Info(fmt.Sprintf("Final image generated\n"))
}

func (cg *CollageGenerator) getMetaData() (*CollageMetaData, error) {
	metaDataFile, err := cg.store.GetFile(cg.collageId)
	if err != nil {
		return nil, err
	}
	var metaData CollageMetaData
	decoder := json.NewDecoder(metaDataFile)
	decoder.Decode(&metaData)
	return &metaData, nil
}

func (cg *CollageGenerator) createCanvas(
	resolution config.ResolutionConfig,
) (*image.RGBA, error) {
	canvas := image.NewRGBA(image.Rect(
		0,
		0,
		resolution.CollageWidth,
		resolution.CollageHeight,
	))
	return canvas, nil
}

func (cg *CollageGenerator) getImage(id uuid.UUID) (*image.RGBA, error) {
	if cached, exists := cg.imageCache.Get(id); exists {
		return cached.(*image.RGBA), nil
	}
	imageFile, err := cg.store.GetFile(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get image file: %w", err)
	}
	im, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	bounds := im.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, im, bounds.Min, draw.Src)
	cg.imageCache.Add(id, rgba)
	return rgba, nil
}

func (cg *CollageGenerator) fillSections(
	start int,
	chunkSize int,
	canvas *image.RGBA,
	metaData *CollageMetaData,
) {
	for section := start; section < start+chunkSize; section++ {
		// Load image for section.
		im, err := cg.getImage(metaData.SectionMap[section])
		if err != nil {
			cg.logger.Error(fmt.Sprintf(
				"Error getting fill image: %v\n",
				err,
			))
		}
		// Scale it.
		sectionWidth := metaData.Resolution.SectionWidth
		sectionHeight := metaData.Resolution.SectionHeight
		scaledImage := resize.Resize(
			uint(sectionWidth),
			uint(sectionHeight),
			im,
			resize.Lanczos2,
		)
		// Draw it where it needs to go on the canvas.
		numColumns := metaData.Resolution.CollageWidth / sectionWidth
		row := section / numColumns
		col := section % numColumns
		startingPoint := image.Point{
			X: sectionWidth * col,
			Y: sectionHeight * row,
		}
		bounds := image.Rectangle{
			Min: startingPoint,
			Max: startingPoint.Add(scaledImage.Bounds().Size()),
		}
		draw.Draw(
			canvas,
			bounds,
			scaledImage,
			scaledImage.Bounds().Min,
			draw.Src,
		)
	}
}

func (cg *CollageGenerator) buildFinalImage(
	canvas *image.RGBA,
	metaData *CollageMetaData,
) error {
	numThreads := 10
	chunkSize := len(metaData.SectionMap) / numThreads
	var wg sync.WaitGroup
	for thread := 0; thread < numThreads; thread++ {
		wg.Add(1)
		startingSection := thread * chunkSize
		go func(start int) {
			defer wg.Done()
			cg.fillSections(
				start,
				chunkSize,
				canvas,
				metaData,
			)
		}(startingSection)
	}
	wg.Wait()
	var buf bytes.Buffer
	err := png.Encode(&buf, canvas)
	if err != nil {
		return err
	}
	cg.store.PutFile(cg.collageImageId, &buf)
	return nil
}
