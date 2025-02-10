package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/utils"
)

func getAverageColors(
	imageSetId uuid.UUID,
	serviceContext context.Context,
) ([]*sqlc.AverageColor, error) {
	postgresConfig := config.NewPostgresConfig()
	acRepo, err := repository.NewAverageColorRepository(
		postgresConfig,
		serviceContext,
	)
	if err != nil {
		return nil, errors.New("Failed to get average colors")
	}
	averageColors, err := acRepo.GetByImageSetId(imageSetId)
	if err != nil {
		return nil, errors.New("Failed to get average colors")
	}
	return averageColors, nil
}

func CreateCollage(collage *sqlc.Collage) {
	serviceContext, cancel := context.WithTimeout(
		context.Background(),
		time.Second*30,
	)
	defer cancel()
	service := newCollageService(collage, serviceContext)
	service.determineImagePlacements()
}

type collageService struct {
	logger          *ServiceLogger
	numThreads      int
	collage         *sqlc.Collage
	averageColors   []*sqlc.AverageColor
	resolution      *config.ResolutionConfig
	sectionMap      []uuid.UUID
	sectionAverages []*color.RGBA
	store           datastore.Store
}

func newCollageService(
	collage *sqlc.Collage,
	serviceContext context.Context,
) *collageService {
	logger := NewServiceLogger("collage")
	averageColors, err := getAverageColors(
		collage.ImageSetID,
		serviceContext,
	)
	if err != nil {
		logger.Printf("Failed to get image set images")
	}
	resolution := config.NewResolutionConfig()
	sectionMap := make(
		[]uuid.UUID,
		resolution.XSections*resolution.YSections,
	)
	store := datastore.NewStore()
	return &collageService{
		logger:        logger,
		numThreads:    10,
		collage:       collage,
		averageColors: averageColors,
		resolution:    resolution,
		sectionMap:    sectionMap,
		store:         store,
	}
}

// Get the local average color value of the collage's
// target image by scaling it down to X_SECTIONSxY_SECTIONS
func (cs *collageService) getSectionAverageColors() []*color.RGBA {
	targetImageReader, err := cs.store.GetFile(cs.collage.TargetImageID)
	if err != nil {
		cs.logger.Fatalf(
			"Failed to load target image: %s\n",
			cs.collage.TargetImageID,
		)
	}
	targetImage, _, err := image.Decode(targetImageReader)
	if err != nil {
		cs.logger.Printf(
			"Failed to decode target image: %s\n",
			cs.collage.TargetImageID,
		)
	}
	scaledImage := resize.Resize(
		uint(cs.resolution.XSections),
		uint(cs.resolution.YSections),
		targetImage,
		resize.Lanczos2,
	)
	bounds := scaledImage.Bounds()
	averageColors := make([]*color.RGBA, bounds.Dx()*bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := scaledImage.At(x, y)
			r, g, b, a := c.RGBA()
			color := &color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			}
			averageColors[y*bounds.Dx()+x] = color
		}
	}
	return averageColors
}

// Find the image set image that best fits the given
// section of the target image
func (cs *collageService) findImagesForSections(startSection int, numSections int) {
	cs.logger.Printf(
		"Finding image for sections: %d-%d\n",
		startSection,
		startSection+numSections-1,
	)
	for section := startSection; section < startSection+numSections; section++ {
		var bestFit uuid.UUID
		bestDistance := math.MaxFloat64
		for _, averageColor := range cs.averageColors {
			imageSetAverage := &color.RGBA{
				R: uint8(averageColor.R),
				G: uint8(averageColor.G),
				B: uint8(averageColor.B),
				A: uint8(averageColor.A),
			}
			distance := utils.ColorDistance(*imageSetAverage, *cs.sectionAverages[section])
			if distance < bestDistance {
				bestFit = averageColor.ID
				bestDistance = distance
			}
		}
		cs.sectionMap[section] = bestFit
	}
}

// Find out what image set image goes where in the collage.
func (cs *collageService) determineImagePlacements() {
	cs.logger.Printf("Finding image placements\n")
	totalSections := cs.resolution.XSections * cs.resolution.YSections
	chunkSize := totalSections / cs.numThreads
	cs.sectionAverages = cs.getSectionAverageColors()
	var wg sync.WaitGroup
	for thread := 0; thread < cs.numThreads; thread++ {
		wg.Add(1)
		threadNum := thread
		go func() {
			defer wg.Done()
			cs.findImagesForSections(threadNum*chunkSize, chunkSize)
		}()
	}
	wg.Wait()
	cs.createMetaDataFile()
}

func (cs *collageService) createMetaDataFile() {
	var buf bytes.Buffer
	metaData := CollageMetaData{
		Resolution: *cs.resolution,
		SectionMap: cs.sectionMap,
	}
	err := json.NewEncoder(&buf).Encode(metaData)
	if err != nil {
		cs.logger.Printf("Error encoding metadata file: %s\n", err)
	}
	err = cs.store.PutFile(cs.collage.ID, &buf)
	if err != nil {
		cs.logger.Printf("Error storing metadata file: %s\n", err)
	}
}
