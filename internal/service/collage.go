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

type ACRepoExtender interface {
	GetByImageSetId(id uuid.UUID) ([]*sqlc.AverageColor, error)
}

func CreateCollage(collage *sqlc.Collage) {
	serviceContext, cancel := context.WithTimeout(
		context.Background(),
		time.Second*30,
	)
	defer cancel()
	postgresConfig := config.NewPostgresConfig()
	acRepo, err := repository.NewAverageColorRepository(
		postgresConfig,
		serviceContext,
	)
	if err != nil {
		panic("Couldn't create repo")
	}
	store := datastore.NewStore()
	service := newCollageService(collage, acRepo, store)
	service.determineImagePlacements()
}

type collageService struct {
	logger     *ServiceLogger
	numThreads int
	collage    *sqlc.Collage
	acRepo     ACRepoExtender
	resolution *config.ResolutionConfig
	sectionMap []uuid.UUID
	store      datastore.Store
}

func newCollageService(
	collage *sqlc.Collage,
	acRepo ACRepoExtender,
	store datastore.Store,
) *collageService {
	logger := NewServiceLogger("collage")
	resolution := config.NewResolutionConfig()
	sectionMap := make(
		[]uuid.UUID,
		resolution.XSections*resolution.YSections,
	)
	return &collageService{
		logger:     logger,
		numThreads: 10,
		collage:    collage,
		acRepo:     acRepo,
		resolution: resolution,
		sectionMap: sectionMap,
		store:      store,
	}
}

func (cs *collageService) getAverageColors() ([]*sqlc.AverageColor, error) {
	averageColors, err := cs.acRepo.GetByImageSetId(cs.collage.ImageSetID)
	if err != nil {
		return nil, errors.New("Failed to get average colors")
	}
	return averageColors, nil
}

// Get the local average color value of the collage's
// target image by scaling it down to X_SECTIONSxY_SECTIONS
func (cs *collageService) getSectionAverageColors() ([]*color.RGBA, error) {
	targetImageReader, err := cs.store.GetFile(cs.collage.TargetImageID)
	if err != nil {
		cs.logger.Printf(
			"Failed to load target image: %s\n",
			cs.collage.TargetImageID,
		)
		return nil, err
	}
	targetImage, _, err := image.Decode(targetImageReader)
	if err != nil {
		cs.logger.Printf(
			"Failed to decode target image: %s\n",
			cs.collage.TargetImageID,
		)
		return nil, err
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
	return averageColors, nil
}

// findImagesForSections finds the image set image that best fits the given
// chunk of sections of the target image by comparing the local color of the
// section to the average color of images in the image set.
// It processes a chunk of sections in parallel.
func (cs *collageService) findImagesForSections(
	startSection int,
	numSections int,
	sectionAverages *[]*color.RGBA,
	imageSetAverageColors *[]*sqlc.AverageColor,
) {
	cs.logger.Printf(
		"Finding image for sections: %d-%d\n",
		startSection,
		startSection+numSections-1,
	)
	for section := startSection; section < startSection+numSections; section++ {
		var bestFit uuid.UUID
		bestDistance := math.MaxFloat64
		sectionAverage := (*sectionAverages)[section]
		for _, averageColor := range *imageSetAverageColors {
			imageSetAverage := &color.RGBA{
				R: uint8(averageColor.R),
				G: uint8(averageColor.G),
				B: uint8(averageColor.B),
				A: uint8(averageColor.A),
			}
			distance := utils.ColorDistance(*imageSetAverage, *sectionAverage)
			if distance < bestDistance {
				bestFit = averageColor.ID
				bestDistance = distance
			}
		}
		cs.sectionMap[section] = bestFit
	}
}

// determineImagePlacements processes the target image in batches by:
//  1. Scaling the target image to a configured resolution where the number of pixels corresponds
//     to the final collage resolution
//  2. Retrieves the average colors of each image in the image set.
//  3. Concurrently finds the best fit image by average color for each section in batches.
//  4. Encodes the placements in a metadata file stored by collage id for deferred creation.
func (cs *collageService) determineImagePlacements() {
	cs.logger.Printf("Finding image placements\n")
	totalSections := cs.resolution.XSections * cs.resolution.YSections
	chunkSize := totalSections / cs.numThreads
	sectionAverages, err := cs.getSectionAverageColors()
	if err != nil {
		cs.logger.Printf(
			"Error getting section averages\n: %s",
			err,
		)
		return
	}
	imageSetAverages, err := cs.getAverageColors()
	if err != nil {
		cs.logger.Printf(
			"Error getting image set average colors\n: %s",
			err,
		)
		return
	}
	var wg sync.WaitGroup
	for thread := 0; thread < cs.numThreads; thread++ {
		wg.Add(1)
		threadNum := thread
		go func() {
			defer wg.Done()
			cs.findImagesForSections(
				threadNum*chunkSize,
				chunkSize,
				&sectionAverages,
				&imageSetAverages,
			)
		}()
	}
	wg.Wait()
	cs.createMetaDataFile()
}

// createMetaDataFile generates a file containing the
// image placements for the collage.
// The meta data file will be stored in the configured
// store location under the name of the collage.
func (cs *collageService) createMetaDataFile() {
	var buf bytes.Buffer
	metaData := CollageMetaData{
		Resolution: *cs.resolution,
		SectionMap: cs.sectionMap,
	}
	err := json.NewEncoder(&buf).Encode(metaData)
	if err != nil {
		cs.logger.Printf("Error encoding metadata file: %s\n", err)
		return
	}
	err = cs.store.PutFile(cs.collage.ID, &buf)
	if err != nil {
		cs.logger.Printf("Error storing metadata file: %s\n", err)
	}
}
