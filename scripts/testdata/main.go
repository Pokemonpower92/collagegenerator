package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func createRGBImages(logger *log.Logger) error {
	logger.Println("Creating RGB test images")
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	testFiles := []struct {
		fileName  string
		fillColor color.RGBA
	}{
		{
			fileName:  "RED.png",
			fillColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			fileName:  "GREEN.png",
			fillColor: color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			fileName:  "BLUE.png",
			fillColor: color.RGBA{R: 0, G: 0, B: 255, A: 255},
		},
		{
			fileName:  "BLACK.png",
			fillColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			fileName:  "WHITE.png",
			fillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			fileName:  "TRANSPARENT.png",
			fillColor: color.RGBA{R: 255, G: 255, B: 255, A: 0},
		},
	}
	for _, file := range testFiles {
		draw.Draw(
			img,
			img.Bounds(),
			&image.Uniform{file.fillColor},
			image.Point{},
			draw.Src,
		)
		fullPath := filepath.Join("./resources/testimages", file.fileName)
		f, err := os.Create(fullPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		if err := png.Encode(f, img); err != nil {
			f.Close()
			logger.Println("Failed to encode image.")
		}
		f.Close()
	}
	return nil
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Generating test images...")
	if err := createRGBImages(logger); err != nil {
		logger.Printf("Error creating RGB test images: %s\n", err)
	}
}
