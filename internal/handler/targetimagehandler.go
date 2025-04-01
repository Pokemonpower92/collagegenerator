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

type TargetImageReader = repository.TargetImageReader
type TargetImageWriter = repository.TargetImageWriter

type TargetImageHandler struct {
	reader TargetImageReader
	writer TargetImageWriter
}

func NewTargetImageHandler(reader TargetImageReader, writer TargetImageWriter) *TargetImageHandler {
	return &TargetImageHandler{reader, writer}
}

func (tih *TargetImageHandler) GetTargetImages(w http.ResponseWriter, _ *http.Request, l *slog.Logger) error {
	l.Info("Getting TargetImages")
	targetImages, err := tih.reader.GetAll()
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d target images.", len(targetImages)))
	response.WriteResponse(w, http.StatusOK, targetImages)
	return nil
}

func (tih *TargetImageHandler) GetTargetImageById(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting TargetImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	targetImage, err := tih.reader.Get(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found TargetImage: %v", targetImage))
	response.WriteResponse(w, http.StatusOK, targetImage)
	return nil
}

func (tih *TargetImageHandler) CreateTargetImage(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Creating targetimage")
	var req CreateTargetImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	targetImage, err := tih.writer.Create(sqlc.CreateTargetImageParams{
		ID:          req.TargetImageId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Created target image with id: %s", targetImage.ID))
	response.WriteResponse(w, http.StatusCreated, targetImage)
	return nil
}
