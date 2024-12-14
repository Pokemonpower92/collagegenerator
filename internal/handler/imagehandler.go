package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/response"
)

type ImageHandler struct {
	l     *log.Logger
	store datastore.Store
}

func NewImageHandler() *ImageHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
	store := datastore.NewStore()
	return &ImageHandler{l: l, store: store}
}

func (ish *ImageHandler) GetImages(w http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (ish *ImageHandler) GetImageById(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Getting Image by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	image, err := ish.store.GetImage(id)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, image)
	if err != nil {
		return err
	}
	return nil
}

func (ish *ImageHandler) StoreImage(w http.ResponseWriter, r *http.Request) error {
	ish.l.Printf("Storing image")
	id := uuid.New()
	if err := ish.store.PutImage(id, r.Body); err != nil {
		return err
	}
	ish.l.Printf("Stored image")
	response.WriteResponse(w, http.StatusCreated, id)
	return nil
}
