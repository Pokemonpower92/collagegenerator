package service

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/stubs"
)

type ACRepoExtenderStub struct {
	GetByImageSetFunc func(id uuid.UUID) ([]*sqlc.AverageColor, error)
}

func (acr *ACRepoExtenderStub) GetByImageSetId(id uuid.UUID) ([]*sqlc.AverageColor, error) {
	return acr.GetByImageSetFunc(id)
}

func successRepo() ACRepoExtenderStub {
	return ACRepoExtenderStub{
		GetByImageSetFunc: func(id uuid.UUID) ([]*sqlc.AverageColor, error) {
			return []*sqlc.AverageColor{
				{
					DbID:     1,
					ID:       uuid.New(),
					FileName: "stubFile",
					R:        1,
					G:        1,
					B:        1,
					A:        1,
				},
			}, nil
		},
	}
}

func errorRepo() ACRepoExtenderStub {
	return ACRepoExtenderStub{
		GetByImageSetFunc: func(id uuid.UUID) ([]*sqlc.AverageColor, error) {
			return nil, errors.New("Stub error.")
		},
	}
}

func imageReader() io.Reader {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	imageBytes := buf.Bytes()
	return bytes.NewReader(imageBytes)
}

func textReader() io.Reader {
	return strings.NewReader("Hello, Reader!")
}

func successStore() stubs.StoreStub {
	return stubs.StoreStub{
		GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
			return nil, nil
		},
		GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
			// Return a new reader each time to avoid EOF issues
			return imageReader(), nil
		},
		PutFileFunc: func(id uuid.UUID, reader io.Reader) error {
			return nil
		},
	}
}

func TestGetAverageColors(t *testing.T) {
	testCases := []struct {
		name        string
		repo        ACRepoExtenderStub
		store       stubs.StoreStub
		collage     *sqlc.Collage
		expected    *[]*sqlc.AverageColor
		shouldError bool
	}{
		{
			name:    "Success",
			repo:    successRepo(),
			store:   successStore(),
			collage: &sqlc.Collage{},
			expected: &[]*sqlc.AverageColor{
				{
					DbID:     1,
					ID:       uuid.New(),
					FileName: "stubFile",
					R:        1,
					G:        1,
					B:        1,
					A:        1,
				},
			},
			shouldError: false,
		},
		{
			name:        "Success",
			repo:        errorRepo(),
			store:       successStore(),
			collage:     &sqlc.Collage{},
			expected:    nil,
			shouldError: true,
		},
	}
	for _, test := range testCases {
		service := newCollageService(
			test.collage,
			&test.repo,
			&test.store,
		)
		result, err := service.getAverageColors()
		if !test.shouldError && err != nil {
			t.Errorf(
				"Test %s failed: expected %v, got %v\n",
				test.name,
				test.expected,
				err,
			)
			return
		}
		if test.shouldError && err != nil {
			return
		}
		if result == nil {
			t.Errorf(
				"Test %s failed: expected %v, got %v\n",
				test.name,
				test.expected,
				result,
			)
		}
	}
}

func TestFindImagesForSections(t *testing.T) {
	redUUID := uuid.New()
	testCases := []struct {
		name               string
		sectionAverages    []*color.RGBA
		imageSetAverages   []*sqlc.AverageColor
		expectedSectionMap []uuid.UUID
	}{
		{
			name: "Matches closest color",
			sectionAverages: []*color.RGBA{
				{R: 255, G: 0, B: 0, A: 255},
			},
			imageSetAverages: []*sqlc.AverageColor{
				{
					ID: redUUID,
					R:  255, G: 0, B: 0, A: 255,
				},
				{
					ID: uuid.New(),
					R:  0, G: 255, B: 0, A: 255,
				},
			},
			expectedSectionMap: []uuid.UUID{
				redUUID,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &collageService{
				logger:     NewServiceLogger("test"),
				sectionMap: make([]uuid.UUID, len(tc.sectionAverages)),
			}
			service.findImagesForSections(
				0,
				len(tc.sectionAverages),
				&tc.sectionAverages,
				&tc.imageSetAverages,
			)
			if !reflect.DeepEqual(service.sectionMap, tc.expectedSectionMap) {
				t.Errorf("Expected section map %v, got %v",
					tc.expectedSectionMap, service.sectionMap)
			}
		})
	}
}

func TestCollageService(t *testing.T) {
	testCases := []struct {
		name        string
		repo        ACRepoExtenderStub
		store       stubs.StoreStub
		collage     *sqlc.Collage
		shouldError bool
	}{
		{
			name:  "Success",
			repo:  successRepo(),
			store: successStore(),
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				ImageSetID:    uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: false,
		},
		{
			name: "Error - GetFile fails",
			repo: successRepo(),
			store: stubs.StoreStub{
				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
					return nil, errors.New("RGBA error")
				},
				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
					return nil, errors.New("file error")
				},
				PutFileFunc: func(id uuid.UUID, reader io.Reader) error {
					return nil
				},
			},
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				ImageSetID:    uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: true,
		},
		{
			name:  "Error - GetByImageSetId fails",
			repo:  errorRepo(),
			store: successStore(),
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				ImageSetID:    uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: true,
		},
		{
			name: "Error - PutFile fails",
			repo: successRepo(),
			store: stubs.StoreStub{
				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
					return &image.RGBA{}, nil
				},
				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
					return imageReader(), nil
				},
				PutFileFunc: func(id uuid.UUID, reader io.Reader) error {
					return errors.New("put file error")
				},
			},
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				ImageSetID:    uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := newCollageService(tc.collage, &tc.repo, &tc.store)
			service.determineImagePlacements()
			if !tc.shouldError {
				for i, section := range service.sectionMap {
					if section == uuid.Nil {
						t.Errorf("Section %d has nil UUID in success case", i)
					}
				}
			}
		})
	}
}

func TestGetSectionAverageColors(t *testing.T) {
	testCases := []struct {
		name        string
		store       stubs.StoreStub
		collage     *sqlc.Collage
		shouldError bool
	}{
		{
			name:  "Success",
			store: successStore(),
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: false,
		},
		{
			name: "Error - GetFile fails",
			store: stubs.StoreStub{
				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
					return nil, errors.New("RGBA error")
				},
				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
					return nil, errors.New("file error")
				},
			},
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: true,
		},
		{
			name: "Error - Target image decode fails",
			store: stubs.StoreStub{
				GetRGBAFunc: func(id uuid.UUID) (*image.RGBA, error) {
					return nil, errors.New("RGBA error")
				},
				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
					return textReader(), nil
				},
			},
			collage: &sqlc.Collage{
				ID:            uuid.New(),
				TargetImageID: uuid.New(),
			},
			shouldError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := newCollageService(
				tc.collage,
				&ACRepoExtenderStub{},
				&tc.store,
			)
			colors, _ := service.getSectionAverageColors()
			if !tc.shouldError && colors == nil {
				t.Error("Expected colors but got nil in success case")
			}
		})
	}
}
