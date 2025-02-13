package utils

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
)

func ImageFileToRGBA(reader io.Reader) (*image.RGBA, error) {
	im, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	b := im.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), im, b.Min, draw.Src)
	return dst, nil
}

func CalculateAverageColor(image *image.RGBA) color.RGBA {
	r, g, b, a, totalPixels := 0, 0, 0, 0, 0
	bounds := image.Bounds()
	for y := bounds.Min.Y; y != bounds.Max.Y; y++ {
		for x := bounds.Min.X; x != bounds.Max.X; x++ {
			pixel := image.RGBAAt(x, y)
			r += int(pixel.R)
			g += int(pixel.G)
			b += int(pixel.B)
			a += int(pixel.A)
			totalPixels++
		}
	}
	return color.RGBA{
		R: uint8(r / totalPixels),
		G: uint8(g / totalPixels),
		B: uint8(b / totalPixels),
		A: uint8(a / totalPixels),
	}
}

func ColorDistance(c1, c2 color.RGBA) float64 {
	rf1, gf1, bf1 := float64(c1.R), float64(c1.G), float64(c1.B)
	rf2, gf2, bf2 := float64(c2.R), float64(c2.G), float64(c2.B)

	rMean := (rf1 + rf2) / 2.0

	r := rf1 - rf2
	g := gf1 - gf2
	b := bf1 - bf2

	// Weights based on human perception
	rWeight := 2.0 + rMean/256.0
	gWeight := 4.0
	bWeight := 2.0 + (255.0-rMean)/256.0

	return math.Sqrt(
		rWeight*r*r +
			gWeight*g*g +
			bWeight*b*b,
	)
}
