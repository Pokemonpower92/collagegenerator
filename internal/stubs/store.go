package stubs

import (
	"image"
	"io"

	"github.com/google/uuid"
)

type StoreStub struct {
	GetRGBAFunc func(id uuid.UUID) (*image.RGBA, error)
	GetFileFunc func(id uuid.UUID) (io.Reader, error)
	PutFileFunc func(id uuid.UUID, reader io.Reader) error
}

func (s *StoreStub) GetRGBA(id uuid.UUID) (*image.RGBA, error) {
	return s.GetRGBAFunc(id)
}

func (s *StoreStub) GetFile(id uuid.UUID) (io.Reader, error) {
	return s.GetFileFunc(id)
}

func (s *StoreStub) PutFile(id uuid.UUID, reader io.Reader) error {
	return s.PutFileFunc(id, reader)
}
