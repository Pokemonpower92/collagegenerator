package stubs

import (
	"github.com/google/uuid"
	sqlc "github.com/pokemonpower92/collagegenerator/internal/sqlc/generated"
)

type RepositoryStub[O, R any] struct {
	GetFunc    func(id uuid.UUID) (*O, error)
	GetAllFunc func() ([]*O, error)
	CreateFunc func(req R) (*O, error)
	UpdateFunc func(id uuid.UUID, req R) (*O, error)
	DeleteFunc func(id uuid.UUID) error
}

func (r *RepositoryStub[O, R]) Get(id uuid.UUID) (*O, error) {
	return r.GetFunc(id)
}

func (r *RepositoryStub[O, R]) GetAll() ([]*O, error) {
	return r.GetAllFunc()
}

func (r *RepositoryStub[O, R]) Create(req R) (*O, error) {
	return r.CreateFunc(req)
}

func (r *RepositoryStub[O, R]) Update(id uuid.UUID, req R) (*O, error) {
	return r.UpdateFunc(id, req)
}

func (r *RepositoryStub[O, R]) Delete(id uuid.UUID) error {
	return r.DeleteFunc(id)
}

type ResourceRepositoryStub[O, R any] struct {
	RepositoryStub[O, R]
	GetByResourceIdFunc func(resourceId uuid.UUID) ([]*O, error)
}

func (r *ResourceRepositoryStub[O, R]) GetByResourceId(resourceId uuid.UUID) ([]*O, error) {
	return r.GetByResourceIdFunc(resourceId)
}

type (
	ISRepoStub RepositoryStub[sqlc.ImageSet, sqlc.CreateImageSetParams]
	TIRepoStub RepositoryStub[sqlc.TargetImage, sqlc.CreateTargetImageParams]
	ACRepoStub = ResourceRepositoryStub[sqlc.AverageColor, sqlc.CreateAverageColorParams]
	CRepoStub  RepositoryStub[sqlc.Collage, sqlc.CreateCollageParams]
	CIRepoStub ResourceRepositoryStub[sqlc.CollageImage, uuid.UUID]
)
