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
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageReader = repository.CollageReader
type CollageWriter = repository.CollageWriter

type CollageHandler struct {
	reader CollageReader
	writer CollageWriter
}

func NewCollageHandler(reader CollageReader, writer CollageWriter) *CollageHandler {
	return &CollageHandler{reader, writer}
}

func (ch *CollageHandler) GetCollages(w http.ResponseWriter, _ *http.Request, l *slog.Logger) error {
	l.Info("Getting Collages")
	collages, err := ch.reader.GetAll()
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found %d Collages", len(collages)))
	response.WriteResponse(w, http.StatusOK, collages)
	return nil
}

func (ch *CollageHandler) GetCollageById(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Getting Collage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing id from path")
		return err
	}
	collage, err := ch.reader.Get(id)
	if err != nil {
		return err
	}
	l.Info(fmt.Sprintf("Found Collage: %v", collage))
	response.WriteResponse(w, http.StatusOK, collage)
	return nil
}

func (ch *CollageHandler) CreateCollage(w http.ResponseWriter, r *http.Request, l *slog.Logger) error {
	l.Info("Creating Collage")
	var req sqlc.CreateCollageParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error(fmt.Sprintf("Error parsing request: %s", err))
		return err
	}
	collage, err := ch.writer.Create(req)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating collage: %s", err))
		return err
	}
	l.Info(fmt.Sprintf("Created Collage with id: %s", collage.ID))
	go service.CreateCollageMetaData(collage, l)
	response.WriteResponse(w, http.StatusCreated, collage)
	return nil
}
