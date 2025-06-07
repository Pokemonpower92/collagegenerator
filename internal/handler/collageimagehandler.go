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
		l.Error("Error getting CollageImages", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info(fmt.Sprintf("Found %d CollageImages", len(collageImages)))
	response.WriteSuccessResponse(w, 200, collageImages)
}

func (cih *CollageImageHandler) GetByCollageId(
	w http.ResponseWriter,
	r *http.Request,
	l *slog.Logger,
) {
	l.Info("Getting CollageImage by Collage ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	collageImage, err := cih.repo.GetByResourceId(id)
	if err != nil {
		l.Error(
			"Error getting CollageImage by Collage id",
			"error", err,
			"collage_id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	l.Info("Found CollageImage", "collage_image", collageImage[0])
	response.WriteSuccessResponse(w, 200, collageImage[0])
}

func (cih *CollageImageHandler) CreateCollageImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating CollageImage")
	var req CreateCollageImageRequest
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
	existing, err := cih.repo.Get(req.CollageID)
	if err == nil && existing != nil {
		l.Info("CollageImage already exists", "id", req.CollageID)
		response.WriteSuccessResponse(w, 200, existing)
		return
	}
	collageImage, err := cih.repo.Create(req.CollageID)
	if err != nil {
		l.Error("Error creating CollageImage", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Created CollageImage", "id", collageImage.ID)
	go service.GenerateCollage(collageImage, l)
	response.WriteSuccessResponse(w, 201, collageImage)
}
