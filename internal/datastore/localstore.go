package datastore

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
)

type LocalStore struct {
	Directory string
	logger    *log.Logger
}

func NewLocalStore(directory string) *LocalStore {
	return &LocalStore{
		Directory: directory,
		logger:    log.New(log.Writer(), "localstore ", log.LstdFlags),
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

func (s *LocalStore) GetImages(filePath string) ([]*image.RGBA, error) {
	s.logger.Printf("Reading images from directory: %s", s.Directory)
	d, err := os.ReadDir(s.Directory)
	if err != nil {
		s.logger.Printf("Failed to open directory: %s", err)
		return nil, err
	}
	imageSet := make([]*image.RGBA, 0)
	for i, item := range d {
		path := fmt.Sprintf("%s/%s", s.Directory, item.Name())
		s.logger.Printf("Decoding file %d: %s", i, path)
		f, err := os.Open(path)
		if err != nil {
			s.logger.Printf("Failed to open file: %s", item.Name())
			return nil, err
		}
		defer f.Close()
		im, _, err := image.Decode(f)
		if err != nil {
			s.logger.Printf("Failed to decode image: %s", err)
			return nil, err
		}
		b := im.Bounds()
		rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), im, b.Min, draw.Src)
		imageSet = append(imageSet, rgba)
	}
	return imageSet, nil
}
