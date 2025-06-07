package handler

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	"github.com/pokemonpower92/collagegenerator/internal/store"
)

func GetFiles(w http.ResponseWriter, _ *http.Request) {
	return
}

func GetFileById(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Getting File by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return
	}
	store := store.NewStore(l)
	image, err := store.GetFile(id)
	if err != nil {
		return
	}
	_, err = io.Copy(w, image)
	if err != nil {
		return
	}
	l.Info(fmt.Sprintf("Got File: %s", id))
	return
}

func StoreFile(w http.ResponseWriter, r *http.Request, l *slog.Logger) {
	l.Info("Storing File")
	id := uuid.New()
	store := store.NewStore(l)
	if err := store.PutFile(id, r.Body); err != nil {
		return
	}
	l.Info("Stored File")
	response.WriteSuccessResponse(w, http.StatusCreated, id)
	return
}
