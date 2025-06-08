package service

// import (
// 	"errors"
// 	"image"
// 	"image/color"
// 	"io"
// 	"log/slog"
// 	"reflect"
// 	"testing"

// 	"github.com/google/uuid"
// 	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
// 	"github.com/pokemonpower92/collagegenerator/internal/stubs"
// )

// func TestGetAverageColors(t *testing.T) {
// 	testCases := []struct {
// 		name        string
// 		repo        stubs.ACRepoStub
// 		store       stubs.StoreStub
// 		collage     *sqlc.Collage
// 		expected    *[]*sqlc.AverageColor
// 		shouldError bool
// 	}{
// 		{
// 			name:    "Success",
// 			repo:    successACRepo(),
// 			store:   successStore(),
// 			collage: &sqlc.Collage{},
// 			expected: &[]*sqlc.AverageColor{
// 				{
// 					DbID:     1,
// 					ID:       uuid.New(),
// 					FileName: "stubFile",
// 					R:        1,
// 					G:        1,
// 					B:        1,
// 					A:        1,
// 				},
// 			},
// 			shouldError: false,
// 		},
// 		{
// 			name:        "Success",
// 			repo:        errorACRepo(),
// 			store:       successStore(),
// 			collage:     &sqlc.Collage{},
// 			expected:    nil,
// 			shouldError: true,
// 		},
// 	}
// 	for _, test := range testCases {
// 		service := newCollageMetaDataService(
// 			test.collage,
// 			&test.repo,
// 			&test.store,
// 			slog.Default(),
// 		)
// 		result, err := service.getAverageColors()
// 		if !test.shouldError && err != nil {
// 			t.Errorf(
// 				"Test %s failed: expected %v, got %v\n",
// 				test.name,
// 				test.expected,
// 				err,
// 			)
// 			return
// 		}
// 		if test.shouldError && err != nil {
// 			return
// 		}
// 		if result == nil {
// 			t.Errorf(
// 				"Test %s failed: expected %v, got %v\n",
// 				test.name,
// 				test.expected,
// 				result,
// 			)
// 		}
// 	}
// }

// func TestFindImagesForSections(t *testing.T) {
// 	redUUID := uuid.New()
// 	testCases := []struct {
// 		name               string
// 		sectionAverages    []color.Color
// 		imageSetAverages   []*sqlc.AverageColor
// 		expectedSectionMap []uuid.UUID
// 	}{
// 		{
// 			name: "Matches closest color",
// 			sectionAverages: []color.Color{
// 				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
// 			},
// 			imageSetAverages: []*sqlc.AverageColor{
// 				{
// 					ID: redUUID,
// 					R:  255, G: 0, B: 0, A: 255,
// 				},
// 				{
// 					ID: uuid.New(),
// 					R:  0, G: 255, B: 0, A: 255,
// 				},
// 			},
// 			expectedSectionMap: []uuid.UUID{
// 				redUUID,
// 			},
// 		},
// 	}
// 	for _, test := range testCases {
// 		t.Run(test.name, func(t *testing.T) {
// 			service := &collageMetaDataService{
// 				logger:     slog.Default(),
// 				sectionMap: make([]uuid.UUID, len(test.sectionAverages)),
// 			}
// 			service.findImagesForSections(
// 				0,
// 				len(test.sectionAverages),
// 				&test.sectionAverages,
// 				&test.imageSetAverages,
// 			)
// 			if !reflect.DeepEqual(service.sectionMap, test.expectedSectionMap) {
// 				t.Errorf("Expected section map %v, got %v",
// 					test.expectedSectionMap, service.sectionMap)
// 			}
// 		})
// 	}
// }

// func TestCollageService(t *testing.T) {
// 	testCases := []struct {
// 		name        string
// 		repo        stubs.ACRepoStub
// 		store       stubs.StoreStub
// 		collage     *sqlc.Collage
// 		shouldError bool
// 	}{
// 		{
// 			name:  "Success",
// 			repo:  successACRepo(),
// 			store: successStore(),
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				ImageSetID:    uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: false,
// 		},
// 		{
// 			name: "Error - GetFile fails",
// 			repo: successACRepo(),
// 			store: stubs.StoreStub{
// 				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
// 					return nil, errors.New("RGBA error")
// 				},
// 				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
// 					return nil, errors.New("file error")
// 				},
// 				PutFileFunc: func(id uuid.UUID, reader io.Reader) error {
// 					return nil
// 				},
// 			},
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				ImageSetID:    uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name:  "Error - GetByImageSetId fails",
// 			repo:  errorACRepo(),
// 			store: successStore(),
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				ImageSetID:    uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name: "Error - PutFile fails",
// 			repo: successACRepo(),
// 			store: stubs.StoreStub{
// 				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
// 					return &image.RGBA{}, nil
// 				},
// 				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
// 					return imageReader(), nil
// 				},
// 				PutFileFunc: func(id uuid.UUID, reader io.Reader) error {
// 					return errors.New("put file error")
// 				},
// 			},
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				ImageSetID:    uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: true,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			service := newCollageMetaDataService(
// 				tc.collage,
// 				&tc.repo,
// 				&tc.store,
// 				slog.Default(),
// 			)
// 			service.determineImagePlacements()
// 			if !tc.shouldError {
// 				for i, section := range service.sectionMap {
// 					if section == uuid.Nil {
// 						t.Errorf("Section %d has nil UUID in success case", i)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

// func TestGetSectionAverageColors(t *testing.T) {
// 	testCases := []struct {
// 		name        string
// 		store       stubs.StoreStub
// 		collage     *sqlc.Collage
// 		shouldError bool
// 	}{
// 		{
// 			name:  "Success",
// 			store: successStore(),
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: false,
// 		},
// 		{
// 			name: "Error - GetFile fails",
// 			store: stubs.StoreStub{
// 				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
// 					return nil, errors.New("RGBA error")
// 				},
// 				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
// 					return nil, errors.New("file error")
// 				},
// 			},
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: true,
// 		},
// 		{
// 			name: "Error - Target image decode fails",
// 			store: stubs.StoreStub{
// 				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
// 					return nil, errors.New("RGBA error")
// 				},
// 				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
// 					return textReader(), nil
// 				},
// 			},
// 			collage: &sqlc.Collage{
// 				ID:            uuid.New(),
// 				TargetImageID: uuid.New(),
// 			},
// 			shouldError: true,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			service := newCollageMetaDataService(
// 				tc.collage,
// 				&stubs.ACRepoStub{},
// 				&tc.store,
// 				slog.Default(),
// 			)
// 			colors, _ := service.getSectionAverageColors()
// 			if !tc.shouldError && colors == nil {
// 				t.Error("Expected colors but got nil in success case")
// 			}
// 		})
// 	}
// }
