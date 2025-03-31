package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/imageprocessing"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CreateAverageColorRequest struct {
	ImagesetID     uuid.UUID `json:"imageset_id"`
	AverageColorID uuid.UUID `json:"averagecolor_id"`
}

type AverageColorReader = repository.AverageColorReader
type AverageColorWriter = repository.AverageColorWriter

type AverageColorHandler struct {
	reader AverageColorReader
	writer AverageColorWriter
}

func NewAverageColorHandler(
	reader AverageColorReader,
	writer AverageColorWriter,
) *AverageColorHandler {
	return &AverageColorHandler{reader, writer}
}

func (ach *AverageColorHandler) GetAverageColors(w http.ResponseWriter, _ *http.Request, l *slog.Logger) error {
	l.Info("Getting AverageColors")
	averageColors, err := ach.reader.GetAll()
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d AverageColors", len(averageColors)))
	response.WriteResponse(w, http.StatusOK, averageColors)
	return nil
}

func (ach *AverageColorHandler) GetAverageColorById(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting AverageColor by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColor, err := ach.reader.Get(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found AverageColor: %v", averageColor))
	response.WriteResponse(w, http.StatusOK, averageColor)
	return nil
}

func (ach *AverageColorHandler) GetByImageSetId(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting AverageColor by ImageSet ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	averageColors, err := ach.reader.GetByResourceId(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d AverageColors", len(averageColors)))
	response.WriteResponse(w, http.StatusOK, averageColors)
	return nil
}

func (ach *AverageColorHandler) CreateAverageColor(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Creating AverageColor")
	var req CreateAverageColorRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	store := datastore.NewStore(l)
	image, err := store.GetRGBA(req.AverageColorID)
	if err != nil {
		return err
	}
	average := imageprocessing.CalculateAverageColor(image)
	averageColor, err := ach.writer.Create(sqlc.CreateAverageColorParams{
		ID:         req.AverageColorID,
		ImagesetID: req.ImagesetID,
		FileName:   req.AverageColorID.String(),
		R:          int32(average.R),
		G:          int32(average.G),
		B:          int32(average.B),
		A:          int32(average.A),
	})
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Created AverageColor with id: %s", averageColor.ID))
	response.WriteResponse(w, http.StatusCreated, averageColor)
	return nil
}
