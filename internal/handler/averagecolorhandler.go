package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/client"
	"github.com/pokemonpower92/collagegenerator/internal/imageprocessing"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
	"github.com/pokemonpower92/collagegenerator/internal/store"
)

type CreateAverageColorRequest struct {
	ImagesetID     uuid.UUID `json:"imageset_id"`
	AverageColorID uuid.UUID `json:"averagecolor_id"`
}

type AverageColorHandler struct {
	repo repository.ACRepo
}

func NewAverageColorHandler(repo repository.ACRepo) *AverageColorHandler {
	return &AverageColorHandler{repo: repo}
}

func (ach *AverageColorHandler) GetAverageColors(w http.ResponseWriter, _ *http.Request, l *slog.Logger) {
	l.Info("Getting AverageColors")
	averageColors, err := ach.repo.GetAll()
	if err != nil {
		l.Error("Error getting AverageColors", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info(fmt.Sprintf("Found %d AverageColors", len(averageColors)))
	response.WriteSuccessResponse(w, 200, averageColors)
}

func (ach *AverageColorHandler) GetAverageColorById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting AverageColor by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	averageColor, err := ach.repo.Get(id)
	if err != nil {
		l.Error(
			"Error getting AverageColor",
			"error", err,
			"id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	l.Info("Found AverageColor", "average_color", averageColor)
	response.WriteSuccessResponse(w, 200, averageColor)
}

func (ach *AverageColorHandler) GetByImageSetId(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting AverageColor by ImageSet ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	averageColors, err := ach.repo.GetByResourceId(id)
	if err != nil {
		l.Error(
			"Error getting AverageColor by ImageSet id",
			"error", err,
			"image_set_id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	l.Info(fmt.Sprintf("Found %d AverageColors", len(averageColors)))
	response.WriteSuccessResponse(w, 200, averageColors)
}

func (ach *AverageColorHandler) CreateAverageColor(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating AverageColor")
	var req CreateAverageColorRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error(
			"Error parsing request",
			"error", err,
			"request", r.Body,
		)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	existing, err := ach.repo.Get(req.AverageColorID)
	if err == nil && existing != nil {
		l.Info("AverageColor already exists", "id", req.AverageColorID)
		response.WriteSuccessResponse(w, 200, existing)
		return
	}
	fileReader := client.NewFileReader("http://filestore:8081/files/", l)
	fileResponse, err := fileReader.GetFile(req.AverageColorID)
	if err != nil {
		l.Error("Error getting file", "error", err)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	image, err := store.GetRGBA(fileResponse)
	if err != nil {
		l.Error("Error converting file", "err", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	average := imageprocessing.CalculateAverageColor(image)
	averageColor, err := ach.repo.Create(sqlc.CreateAverageColorParams{
		ID:         req.AverageColorID,
		ImagesetID: req.ImagesetID,
		FileName:   req.AverageColorID.String(),
		R:          int32(average.R),
		G:          int32(average.G),
		B:          int32(average.B),
		A:          int32(average.A),
	})
	if err != nil {
		l.Error("Error creating AverageColor", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Created AverageColor", "id", averageColor.ID)
	response.WriteSuccessResponse(w, 201, averageColor)
}
