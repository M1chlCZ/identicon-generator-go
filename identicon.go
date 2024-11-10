package identicon

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

// Config holds the configuration for generating identicons
type Config struct {
	// Width of the generated image in pixels
	Width int
	// Height of the generated image in pixels
	Height int
	// GridSize determines the complexity of the pattern (e.g., 5 for 5x5 grid)
	GridSize int
	// Grayscale determines if the image should be in grayscale
	Grayscale bool
}

// DefaultConfig returns a Config with recommended default values
func DefaultConfig() Config {
	return Config{
		Width:     420,
		Height:    420,
		GridSize:  5,
		Grayscale: false,
	}
}

// GenerateIdenticon generates an identicon using the provided seed
// It uses the default configuration
func GenerateIdenticon(seed int) image.Image {
	return GenerateIdenticonWithConfig(seed, DefaultConfig())
}

// GenerateIdenticonWithConfig generates an identicon using the provided seed and configuration
func GenerateIdenticonWithConfig(seed int, config Config) image.Image {
	hash := md5.Sum([]byte(fmt.Sprintf("%d", seed)))
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	cellWidth := config.Width / config.GridSize
	cellHeight := config.Height / config.GridSize

	for i := 0; i < config.GridSize; i++ {
		for j := 0; j < config.GridSize/2+1; j++ {
			if hash[(i*config.GridSize+j)%16]%2 == 0 {
				col := getColor(hash[i%5*3:i%5*3+3], config.Grayscale)
				drawRect(img, j*cellWidth, i*cellHeight, cellWidth, cellHeight, col)
				drawRect(img, (config.GridSize-1-j)*cellWidth, i*cellHeight, cellWidth, cellHeight, col)
			}
		}
	}

	return img
}

func getColor(bytes []byte, grayscale bool) color.Color {
	if grayscale {
		v := uint8((int(bytes[0]) + int(bytes[1]) + int(bytes[2])) / 3)
		return color.Gray{Y: v}
	}
	return color.RGBA{R: bytes[0], G: bytes[1], B: bytes[2], A: 255}
}

func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

// SaveImage saves the provided image to a file
func SaveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("error closing file:", err)
			return
		}
	}(f)
	return png.Encode(f, img)
}
