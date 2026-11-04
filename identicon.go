package identicon

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
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
	// Style selects the drawing algorithm. The zero value is StyleSymmetric.
	Style Style
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

// Validate reports whether the configuration can produce a visible image.
func (c Config) Validate() error {
	switch {
	case c.Width <= 0:
		return fmt.Errorf("width must be positive, got %d", c.Width)
	case c.Height <= 0:
		return fmt.Errorf("height must be positive, got %d", c.Height)
	case c.GridSize < 1:
		return fmt.Errorf("grid size must be positive, got %d", c.GridSize)
	case c.Width < c.GridSize:
		return fmt.Errorf("width %d is smaller than grid size %d", c.Width, c.GridSize)
	case c.Height < c.GridSize:
		return fmt.Errorf("height %d is smaller than grid size %d", c.Height, c.GridSize)
	case c.Style > StyleRings:
		return fmt.Errorf("unknown style %d", c.Style)
	}
	return nil
}

// GenerateIdenticon generates an identicon using the provided seed
// It uses the default configuration
func GenerateIdenticon(seed int) image.Image {
	return GenerateIdenticonWithConfig(seed, DefaultConfig())
}

// GenerateIdenticonWithConfig generates an identicon using the provided seed and configuration
//
// The function panics when the configuration is invalid. Use [Render] to get
// an error instead.
func GenerateIdenticonWithConfig(seed int, config Config) image.Image {
	if err := config.Validate(); err != nil {
		panic(err)
	}
	hash := md5.Sum([]byte(fmt.Sprintf("%d", seed)))
	return renderSymmetric(hash[:], config, getColor)
}

// Render generates an identicon from a string seed.
//
// The seed is hashed with SHA-256, so the same seed always produces the same
// image. The configuration is validated before rendering. The output differs
// from [GenerateIdenticon], which keeps the legacy MD5 and integer-seed
// behaviour.
func Render(seed string, config Config) (image.Image, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	hash := sha256.Sum256([]byte(seed))
	switch config.Style {
	case StyleMosaic:
		return renderMosaic(hash[:], config), nil
	case StyleRings:
		return renderRings(hash[:], config), nil
	default:
		return renderSymmetric(hash[:], config, getReadableColor), nil
	}
}

// WritePNG encodes img as a PNG to w.
func WritePNG(w io.Writer, img image.Image) error {
	return png.Encode(w, img)
}

func renderSymmetric(hash []byte, config Config, colorFor func([]byte, bool) color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	cellWidth := config.Width / config.GridSize
	cellHeight := config.Height / config.GridSize

	for i := 0; i < config.GridSize; i++ {
		for j := 0; j < config.GridSize/2+1; j++ {
			if hash[(i*config.GridSize+j)%16]%2 == 0 {
				col := colorFor(hash[i%5*3:i%5*3+3], config.Grayscale)
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

func getReadableColor(bytes []byte, grayscale bool) color.Color {
	if grayscale {
		v := uint8((int(bytes[0]) + int(bytes[1]) + int(bytes[2])) / 3)
		if v > 160 {
			v = 160
		}
		return color.Gray{Y: v}
	}
	r, g, b := bytes[0], bytes[1], bytes[2]
	highest := max(r, g, b)
	if highest > 160 {
		r = uint8(int(r) * 160 / int(highest))
		g = uint8(int(g) * 160 / int(highest))
		b = uint8(int(b) * 160 / int(highest))
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func drawRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			img.Set(x+dx, y+dy, col)
		}
	}
}

// SaveImage saves the provided image to a file
func SaveImage(img image.Image, filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()
	return WritePNG(f, img)
}
