package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	"github.com/pokemonpower92/collagegenerator/internal/service"
)

type CreateCollageImageRequest struct {
	CollageID uuid.UUID `json:"collage_id"`
}

type CollageImageHandler struct {
	repo repository.CIRepo
}

func NewCollageImageHandler(repo repository.CIRepo) *CollageImageHandler {
	return &CollageImageHandler{repo: repo}
}

func (cih *CollageImageHandler) GetCollageImages(w http.ResponseWriter, _ *http.Request, l *slog.Logger) {
	l.Info("Getting CollageImages")
	collageImages, err := cih.repo.GetAll()
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Found %d CollageImages", len(collageImages)))
	response.WriteSuccessResponse(w, http.StatusOK, collageImages)
	return
}

func (cih *CollageImageHandler) GetCollageImageByCollageId(
	w http.ResponseWriter,
	r *http.Request,
	l *slog.Logger,
) {
	l.Info("Getting CollageImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return
	}
	collageImage, err := cih.repo.GetByResourceId(id)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Found CollageImage: %v", collageImage[0]))
	response.WriteSuccessResponse(w, http.StatusOK, collageImage[0])
	return
}

func (cih *CollageImageHandler) CreateCollageImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating CollageImage")
	var req CreateCollageImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return
	}
	collageImage, err := cih.repo.Create(req.CollageID)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Created CollageImage with id: %s", collageImage.ID))
	go service.GenerateCollage(collageImage, l)
	response.WriteSuccessResponse(w, http.StatusCreated, collageImage)
	return
}
