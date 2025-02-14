package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/response"
)

type FileHandler struct {
	l     *log.Logger
	store datastore.Store
}

func NewFileHandler() *FileHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
	store := datastore.NewStore()
	return &FileHandler{l: l, store: store}
}

func (ish *FileHandler) GetFiles(w http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (ish *FileHandler) GetFileById(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Getting File by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	image, err := ish.store.GetFile(id)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, image)
	if err != nil {
		return err
	}
	ish.l.Printf("Got File: %s", id)
	return nil
}

func (ish *FileHandler) StoreFile(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Storing File")
	id := uuid.New()
	if err := ish.store.PutFile(id, r.Body); err != nil {
		return err
	}
	ish.l.Printf("Stored File")
	response.WriteResponse(w, http.StatusCreated, id)
	return nil
}
