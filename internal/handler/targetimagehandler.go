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
		return
	}
	l.Info(fmt.Sprintf("Found %d target images.", len(targetImages)))
	response.WriteSuccessResponse(w, http.StatusOK, targetImages)
	return
}

func (tih *TargetImageHandler) GetTargetImageById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting TargetImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return
	}
	targetImage, err := tih.repo.Get(id)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Found TargetImage: %v", targetImage))
	response.WriteSuccessResponse(w, http.StatusOK, targetImage)
	return
}

func (tih *TargetImageHandler) CreateTargetImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating targetimage")
	var req CreateTargetImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}
	targetImage, err := tih.repo.Create(sqlc.CreateTargetImageParams{
		ID:          req.TargetImageId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Created target image with id: %s", targetImage.ID))
	response.WriteSuccessResponse(w, http.StatusCreated, targetImage)
	return
}
