package datastore

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStore struct {
	Directory string
	logger    *log.Logger
}

func NewLocalStore(directory string) *LocalStore {
	return &LocalStore{
		Directory: directory,
		logger:    log.New(log.Writer(), "", log.LstdFlags),
	}
}

func imageToRGBA(src image.Image) *image.RGBA {
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

func (s *LocalStore) GetRGBA(id uuid.UUID) (*image.RGBA, error) {
	f, err := s.GetFile(id)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(f)
	if err != nil {
		s.logger.Printf("Failed to decode image: %s", err)
		return nil, err
	}
	b := im.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgba, rgba.Bounds(), im, b.Min, draw.Src)
	return rgba, nil
}

func (s *LocalStore) GetFile(id uuid.UUID) (io.Reader, error) {
	path := fmt.Sprintf("%s/%s", s.Directory, id.String())
	s.logger.Printf("Getting File: %s", path)
	f, err := os.Open(path)
	if err != nil {
		s.logger.Printf("Failed to open File: %s", id.String())
		return nil, err
	}
	return f, nil
}

func (s *LocalStore) PutFile(id uuid.UUID, reader io.Reader) error {
	dst, err := os.Create(filepath.Join(s.Directory, id.String()))
	if err != nil {
		return err
	}
	defer dst.Close()
	s.logger.Printf("Created File destination: %s", dst.Name())
	if _, err := io.Copy(dst, reader); err != nil {
		return err
	}
	s.logger.Printf("Successfully stored File")
	return nil
}
