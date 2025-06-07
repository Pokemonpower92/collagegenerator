package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CreateTargetImageRequest struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	TargetImageId uuid.UUID `json:"targetimage_id"`
}

type TargetImageHandler struct {
	repo repository.TIRepo
}

func NewTargetImageHandler(repo repository.TIRepo) *TargetImageHandler {
	return &TargetImageHandler{repo: repo}
}

func (tih *TargetImageHandler) GetTargetImages(w http.ResponseWriter, _ *http.Request, l *slog.Logger) {
	l.Info("Getting TargetImages")
	targetImages, err := tih.repo.GetAll()
	if err != nil {
		l.Error("Error getting TargetImages", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info(fmt.Sprintf("Found %d TargetImages", len(targetImages)))
	response.WriteSuccessResponse(w, 200, targetImages)
}

func (tih *TargetImageHandler) GetTargetImageById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting TargetImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	targetImage, err := tih.repo.Get(id)
	if err != nil {
		l.Error(
			"Error getting TargetImage",
			"error", err,
			"id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	l.Info("Found TargetImage", "target_image", targetImage)
	response.WriteSuccessResponse(w, 200, targetImage)
}

func (tih *TargetImageHandler) CreateTargetImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating TargetImage")
	var req CreateTargetImageRequest
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
	targetImage, err := tih.repo.Create(sqlc.CreateTargetImageParams{
		ID:          req.TargetImageId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		l.Error("Error creating TargetImage", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Created TargetImage", "id", targetImage.ID)
	response.WriteSuccessResponse(w, 201, targetImage)
}
