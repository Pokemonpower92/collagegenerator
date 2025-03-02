package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/pokemonpower92/collagegenerator/internal/datastore"
	"github.com/pokemonpower92/collagegenerator/internal/repository"
	"github.com/pokemonpower92/collagegenerator/internal/response"
	"github.com/pokemonpower92/collagegenerator/internal/service"
)

type CreateCollageImageRequest struct {
	CollageID uuid.UUID `json:"collage_id"`
}

type CollageImageHandler struct {
	l     *log.Logger
	repo  repository.CIRepo
	store datastore.Store
}

func NewCollageImageHandler(repo repository.CIRepo) *CollageImageHandler {
	l := log.New(log.Writer(), "", log.LstdFlags)
	store := datastore.NewStore()
	return &CollageImageHandler{l: l, repo: repo, store: store}
}

func (cih *CollageImageHandler) GetCollageImages(w http.ResponseWriter, _ *http.Request) error {
	cih.l.Printf("Getting CollageImages")
	collageImages, err := cih.repo.GetAll()
	if err != nil {
		return err
	}
	cih.l.Printf("Found %d CollageImages", len(collageImages))
	response.WriteResponse(w, http.StatusOK, collageImages)
	return nil
}

func (cih *CollageImageHandler) GetCollageImageByCollageId(
	w http.ResponseWriter,
	r *http.Request,
) error {
	cih.l.Printf("Getting CollageImage by ID")
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	collageImage, err := cih.repo.GetByResourceId(id)
	if err != nil {
		return err
	}
	cih.l.Printf("Found CollageImage: %v", collageImage[0])
	response.WriteResponse(w, http.StatusOK, collageImage[0])
	return nil
}

func (cih *CollageImageHandler) CreateCollageImage(w http.ResponseWriter, r *http.Request) error {
	cih.l.Printf("Creating CollageImage")
	var req CreateCollageImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	collageImage, err := cih.repo.Create(req.CollageID)
	if err != nil {
		return err
	}
	cih.l.Printf("Created CollageImage with id: %s", collageImage.ID)
	go service.GenerateCollage(collageImage)
	response.WriteResponse(w, http.StatusCreated, collageImage)
	return nil
}
