package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type collageService struct {
	l             *ServiceLogger
	collage       *sqlc.Collage
	averageColors []*sqlc.AverageColor
	store         datastore.Store
	imageMap      map[int]string
}

func getAverageColors(
	imageSetId uuid.UUID,
	serviceContext context.Context,
) ([]*sqlc.AverageColor, error) {
	postgresConfig := config.NewPostgresConfig()
	acRepo, err := repository.NewAverageColorRepository(postgresConfig, serviceContext)
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
	serviceContext, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	service := newCollageService(collage, serviceContext)
	service.determineImagePlacements()
	service.placeImagesInCollage()
}

func newCollageService(
	collage *sqlc.Collage,
	serviceContext context.Context,
) *collageService {
	l := NewServiceLogger("collage")
	averageColors, err := getAverageColors(collage.ImageSetID, serviceContext)
	if err != nil {
		l.Fatalf("Failed to get image set images")
	}
	store := datastore.NewStore()
	imageMap := make(map[int]string)
	return &collageService{
		l:             l,
		collage:       collage,
		averageColors: averageColors,
		store:         store,
		imageMap:      imageMap,
	}
}

// Find the image set image that best fits the given
// section of the target image
func (cs *collageService) findImageForSection(section int) {
	cs.l.Printf("Finding image for section: %d\n", section)
	// Calculate the average color of the section.
}

// Calculate the number of sections the target image can
// be split into given the configured resolution of the
// collage
func (cs *collageService) calculateSections() int {
	return 1
}

// Find out what image set image goes where in the collage.
func (cs *collageService) determineImagePlacements() {
	cs.l.Printf("Finding image placements\n")
	var wg sync.WaitGroup
	numSections := cs.calculateSections()
	for section := 0; section < numSections; section++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cs.findImageForSection(section)
		}()
	}
	wg.Wait()
}

func (cs *collageService) placeImagesInCollage() {
	cs.l.Printf("Placing images in the collage\n")
}
