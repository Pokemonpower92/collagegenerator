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

type ImageSetReader = repository.ImageSetReader
type ImageSetWriter = repository.ImageSetWriter

type ImageSetHandler struct {
	reader ImageSetReader
	writer ImageSetWriter
}

func NewImageSetHandler(reader ImageSetReader, writer ImageSetWriter) *ImageSetHandler {
	return &ImageSetHandler{reader, writer}
}

func (ish *ImageSetHandler) GetImageSets(w http.ResponseWriter, _ *http.Request, l *slog.Logger) error {
	l.Info("Getting ImageSets")
	imageSets, err := ish.reader.GetAll()
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d ImageSets", len(imageSets)))
	response.WriteResponse(w, http.StatusOK, imageSets)
	return nil
}

func (ish *ImageSetHandler) GetImageSetById(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting ImageSet by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	imageSet, err := ish.reader.Get(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found ImageSet: %v", imageSet))
	response.WriteResponse(w, http.StatusOK, imageSet)
	return nil
}

func (ish *ImageSetHandler) CreateImageSet(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Creating ImageSet")
	var req sqlc.CreateImageSetParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	imageSet, err := ish.writer.Create(req)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Created ImageSet with id: %s", imageSet.ID))
	response.WriteResponse(w, http.StatusCreated, imageSet)
	return nil
}
