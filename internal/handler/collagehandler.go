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

type CollageHandler struct {
	repo repository.CRepo
}

func NewCollageHandler(repo repository.CRepo) *CollageHandler {
	return &CollageHandler{repo: repo}
}

func (ch *CollageHandler) GetCollages(w http.ResponseWriter, _ *http.Request, l *slog.Logger) {
	l.Info("Getting Collages")
	collages, err := ch.repo.GetAll()
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Found %d Collages", len(collages)))
	response.WriteSuccessResponse(w, http.StatusOK, collages)
	return
}

func (ch *CollageHandler) GetCollageById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting Collage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing id from path")
		return
	}
	collage, err := ch.repo.Get(id)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Found Collage: %v", collage))
	response.WriteSuccessResponse(w, http.StatusOK, collage)
	return
}

func (ch *CollageHandler) CreateCollage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating Collage")
	var req sqlc.CreateCollageParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		l.Error(fmt.Sprintf("Error parsing request: %s", err))
		return
	}
	collage, err := ch.repo.Create(req)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating collage: %s", err))
		return
	}
	l.Info(fmt.Sprintf("Created Collage with id: %s", collage.ID))
	go service.CreateCollageMetaData(collage, l)
	response.WriteSuccessResponse(w, http.StatusCreated, collage)
	return
}
