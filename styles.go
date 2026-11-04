package identicon

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strings"
)

// Style selects the identicon drawing algorithm.
type Style uint8

const (
	// StyleSymmetric is the original mirrored grid.
	StyleSymmetric Style = iota
	// StyleMosaic is a denser mirrored grid with three colors.
	StyleMosaic
	// StyleRings draws concentric colored rings.
	StyleRings
)

// String returns the CLI name of the style.
func (s Style) String() string {
	switch s {
	case StyleSymmetric:
		return "symmetric"
	case StyleMosaic:
		return "mosaic"
	case StyleRings:
		return "rings"
	default:
		return "unknown"
	}
}

// ParseStyle converts a CLI name into a Style.
func ParseStyle(name string) (Style, error) {
	switch strings.ToLower(name) {
	case "symmetric":
		return StyleSymmetric, nil
	case "mosaic":
		return StyleMosaic, nil
	case "rings":
		return StyleRings, nil
	default:
		return 0, fmt.Errorf("unknown style %q", name)
	}
}

func renderMosaic(hash []byte, config Config) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	grid := config.GridSize * 2
	if grid > config.Width {
		grid = config.Width
	}
	if grid > config.Height {
		grid = config.Height
	}
	cellWidth := config.Width / grid
	cellHeight := config.Height / grid
	palette := readablePalette(hash, config.Grayscale)

	for y := 0; y < grid; y++ {
		for x := 0; x < grid/2+1; x++ {
			value := hash[(y*grid+x)%len(hash)]
			if value%3 == 0 {
				continue
			}
			col := palette[int(value)%len(palette)]
			drawRect(img, x*cellWidth, y*cellHeight, cellWidth, cellHeight, col)
			drawRect(img, (grid-1-x)*cellWidth, y*cellHeight, cellWidth, cellHeight, col)
		}
	}
	return img
}

func renderRings(hash []byte, config Config) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, config.Width, config.Height))
	background := color.Color(color.White)
	if config.Grayscale {
		background = color.Gray{Y: 255}
	}
	draw.Draw(img, img.Bounds(), image.NewUniform(background), image.Point{}, draw.Src)

	centerX := float64(config.Width) / 2
	centerY := float64(config.Height) / 2
	radius := float64(min(config.Width, config.Height)) / 2
	if radius <= 0 {
		return img
	}
	palette := readablePalette(hash, config.Grayscale)
	bands := config.GridSize * 2

	for y := 0; y < config.Height; y++ {
		for x := 0; x < config.Width; x++ {
			distance := math.Hypot(float64(x)+0.5-centerX, float64(y)+0.5-centerY)
			if distance > radius {
				continue
			}
			band := int(distance / radius * float64(bands))
			if band%2 == 1 {
				continue
			}
			img.Set(x, y, palette[(band/2)%len(palette)])
		}
	}
	return img
}

func readablePalette(hash []byte, grayscale bool) []color.Color {
	return []color.Color{
		getReadableColor(hash[0:3], grayscale),
		getReadableColor(hash[3:6], grayscale),
		getReadableColor(hash[6:9], grayscale),
	}
}
