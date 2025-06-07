package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	"github.com/pokemonpower92/collagegenerator/internal/store"
)

func GetFileById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting File by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		l.Error(
			"Error parsing request",
			"error", err,
		)
		response.WriteErrorResponse(w, 422, err)
		return
	}
	store := store.NewStore(l)
	image, err := store.GetFile(id)
	if err != nil {
		l.Error(
			"Error getting File",
			"error", err,
			"file_id", id,
		)
		response.WriteErrorResponse(w, 404, err)
		return
	}
	_, err = io.Copy(w, image)
	if err != nil {
		l.Error(
			"Error parsing File contents",
			"error", err,
		)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Got File", "id", id)
}

func StoreFile(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Storing File")
	id := uuid.New()
	store := store.NewStore(l)
	if err := store.PutFile(id, r.Body); err != nil {
		l.Error(
			"Error storing File",
			"error", err,
		)
		response.WriteErrorResponse(w, 500, err)
		return
	}
	l.Info("Stored File", "id", id)
	response.WriteSuccessResponse(w, http.StatusCreated, id)
}
