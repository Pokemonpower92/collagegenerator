package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/config"
	"github.com/pokemonpower92/collagegenerator/internal/client"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type CollageHandler struct {
	repo   repository.CRepo
	sender client.MessageSender
}

func NewCollageHandler(
	repo repository.CRepo,
	sender client.MessageSender,
) *CollageHandler {
	return &CollageHandler{repo, sender}
}

func (ch *CollageHandler) GetCollages(w http.ResponseWriter, _ *http.Request, l *slog.Logger) {
	l.Info("Getting Collages")
	collages, err := ch.repo.GetAll()
	if err != nil {
		l.Error("Error getting Collages", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info(fmt.Sprintf("Found %d Collages", len(collages)))
	response.WriteSuccessResponse(w, 200, collages)
}

func (ch *CollageHandler) GetCollageById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting Collage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error("Error parsing request", "error", err)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	collage, err := ch.repo.Get(id)
	if err != nil {
		l.Error(
			"Error getting Collage",
			"error", err,
			"id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	l.Info("Found Collage", "collage", collage)
	response.WriteSuccessResponse(w, 200, collage)
}

func (ch *CollageHandler) CreateCollage(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Creating Collage")
	var req sqlc.CreateCollageParams
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
	collage, err := ch.repo.Create(req)
	if err != nil {
		l.Error("Error creating Collage", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Created Collage", "id", collage.ID)
	collageJSON, err := json.Marshal(collage)
	if err != nil {
		l.Error("Error marshaling Collage", "error", err)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	err = ch.sender.Send(config.METADATA_QUEUE(), string(collageJSON), r.Context())
	if err != nil {
		l.Error(
			"Error sending metadata job",
			"error", err,
		)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	response.WriteSuccessResponse(w, 201, collage)
}
