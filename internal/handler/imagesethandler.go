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

type ImageSetHandler struct {
	repo repository.ISRepo
}

func NewImageSetHandler(repo repository.ISRepo) *ImageSetHandler {
	return &ImageSetHandler{repo: repo}
}

func (ish *ImageSetHandler) GetImageSets(w http.ResponseWriter, _ *http.Request, l *slog.Logger) {
	l.Info("Getting ImageSets")
	imageSets, err := ish.repo.GetAll()
	if err != nil {
		l.Error("Error getting ImageSets", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info(fmt.Sprintf("Found %d ImageSets", len(imageSets)))
	response.WriteSuccessResponse(w, 200, imageSets)
}

func (ish *ImageSetHandler) GetImageSetById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting ImageSet by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	imageSet, err := ish.repo.Get(id)
	if err != nil {
		l.Error(
			"Error getting ImageSet",
			"error", err,
			"id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	l.Info("Found ImageSet", "image_set", imageSet)
	response.WriteSuccessResponse(w, 200, imageSet)
}

func (ish *ImageSetHandler) CreateImageSet(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating ImageSet")
	var req sqlc.CreateImageSetParams
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
	imageSet, err := ish.repo.Create(req)
	if err != nil {
		l.Error("Error creating ImageSet", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Created ImageSet", "id", imageSet.ID)
	response.WriteSuccessResponse(w, 201, imageSet)
}
