package utils

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"strings"
	"testing"
)

func TestCalculateAverageColor(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	testCases := []struct {
		name      string
		fillColor color.RGBA
		expected  color.RGBA
	}{
		{
			name:      "red",
			fillColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expected:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:      "green",
			fillColor: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expected:  color.RGBA{R: 0, G: 255, B: 0, A: 255},
		},
		{
			name:      "blue",
			fillColor: color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expected:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
		},
		{
			name:      "transparent",
			fillColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},
			expected:  color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			draw.Draw(
				img,
				img.Bounds(),
				&image.Uniform{test.fillColor},
				image.Point{},
				draw.Src,
			)
			result := CalculateAverageColor(img)
			if result != test.expected {
				t.Errorf(
					"Test %s failed: expected %v got %v\n",
					test.name,
					test.expected,
					result,
				)
			}
		})
	}
}

func TestImageFileToRGBA(t *testing.T) {
	testCases := []struct {
		name      string
		input     func() io.Reader
		wantErr   bool
		wantColor *color.RGBA // Use pointer to handle nil case for errors
	}{
		{
			name: "Success - Red PNG",
			input: func() io.Reader {
				f, _ := os.Open("../../testimages/RED.png")
				return f
			},
			wantErr:   false,
			wantColor: &color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name: "Invalid Image Data",
			input: func() io.Reader {
				return strings.NewReader("not an image")
			},
			wantErr: true,
		},
		{
			name: "Empty Reader",
			input: func() io.Reader {
				return strings.NewReader("")
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := tc.input()
			if f, ok := reader.(*os.File); ok {
				defer f.Close()
			}
			rgba, err := ImageFileToRGBA(reader)
			if (err != nil) != tc.wantErr {
				t.Errorf(
					"ImageFileToRGBA() error = %v, wantErr %v",
					err,
					tc.wantErr,
				)
				return
			}
			if tc.wantErr {
				return
			}
			if tc.wantColor != nil {
				result := CalculateAverageColor(rgba)
				if result != *tc.wantColor {
					t.Errorf("ImageFileToRGBA() got color = %v, want %v", result, *tc.wantColor)
				}
			}
		})
	}
}
