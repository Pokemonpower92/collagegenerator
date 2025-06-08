package service

// import (
// 	"io"
// 	"log/slog"
// 	"reflect"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/pokemonpower92/collagegenerator/config"
// 	"github.com/pokemonpower92/collagegenerator/internal/stubs"
// )

// func TestGenerate(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		collageImageId uuid.UUID
// 		logger         *slog.Logger
// 		store          stubs.StoreStub
// 	}{
// 		{
// 			name:           "Success",
// 			collageImageId: uuid.New(),
// 			logger:         slog.Default(),
// 			store:          successStore(),
// 		},
// 	}
// 	for _, test := range testCases {
// 		cg := CollageGenerator{
// 			collageImageId: test.collageImageId,
// 			logger:         test.logger,
// 			store:          &test.store,
// 		}
// 		cg.generate()
// 	}
// }

// func TestGetMetaDataFile(t *testing.T) {
// 	testCases := []struct {
// 		name           string
// 		collageImageId uuid.UUID
// 		logger         *slog.Logger
// 		store          stubs.StoreStub
// 		expected       *CollageMetaData
// 		shouldErr      bool
// 	}{
// 		{
// 			name:           "Success",
// 			collageImageId: uuid.New(),
// 			logger:         slog.Default(),
// 			store: stubs.StoreStub{
// 				GetFileFunc: func(id uuid.UUID) (io.Reader, error) {
// 					return metaDataReader(), nil
// 				},
// 			},
// 			expected: &CollageMetaData{
// 				Resolution: config.ResolutionConfig{},
// 				SectionMap: make([]uuid.UUID, 0),
// 			},
// 			shouldErr: false,
// 		},
// 	}
// 	for _, test := range testCases {
// 		cg := CollageGenerator{
// 			collageImageId: test.collageImageId,
// 			logger:         test.logger,
// 			store:          &test.store,
// 		}
// 		result, _ := cg.getMetaData()
// 		if !test.shouldErr && !reflect.DeepEqual(result, test.expected) {
// 			t.Errorf(
// 				"Test %s failed: expected %v, got %v\n",
// 				test.name,
// 				test.expected,
// 				result,
// 			)
// 		}
// 	}
// }
