package identicon

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func pixHash(t *testing.T, img image.Image) string {
	t.Helper()
	rgba, ok := img.(*image.RGBA)
	if !ok {
		t.Fatalf("image type = %T, want *image.RGBA", img)
	}
	sum := sha256.Sum256(rgba.Pix)
	return hex.EncodeToString(sum[:])
}

func TestLegacyOutputGolden(t *testing.T) {
	want := map[int]string{
		0:    "7c5e276478f13b1f85ed1210b5ac3d2e60acf6e5fc721ad07dab4cd9ec33a023",
		1:    "29af22ff14fa1f003c2cf9656257dc4cf9006f0a5885c166c2540b9e891eb3b0",
		123:  "f1b8aaa289e63c9e9715b341810c1ea33e7f219db2c82fe317b56e51fb4d8f88",
		2024: "b39ab11e324c32739888cdbfdf818802ad4670f14c71f03cd960388208954744",
	}
	for seed, hash := range want {
		if got := pixHash(t, GenerateIdenticon(seed)); got != hash {
			t.Errorf("seed %d hash = %s, want %s", seed, got, hash)
		}
	}
}

func TestRender(t *testing.T) {
	config := DefaultConfig()
	first, err := Render("alice@example.com", config)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Render("alice@example.com", config)
	if err != nil {
		t.Fatal(err)
	}
	other, err := Render("bob@example.com", config)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(first.(*image.RGBA).Pix, second.(*image.RGBA).Pix) {
		t.Fatal("same seed produced different pixels")
	}
	if bytes.Equal(first.(*image.RGBA).Pix, other.(*image.RGBA).Pix) {
		t.Fatal("different seeds produced the same pixels")
	}
	if img, err := Render("seed", Config{}); err == nil || img != nil {
		t.Fatalf("Render() = %v, %v; want nil image and an error", img, err)
	}
}

func TestStyles(t *testing.T) {
	seen := map[string]bool{}
	for _, name := range []string{"symmetric", "mosaic", "rings"} {
		style, err := ParseStyle(name)
		if err != nil {
			t.Fatal(err)
		}
		config := DefaultConfig()
		config.Style = style
		img, err := Render("styles", config)
		if err != nil {
			t.Fatal(err)
		}
		hash := pixHash(t, img)
		if seen[hash] {
			t.Fatalf("style %s duplicated another style", name)
		}
		seen[hash] = true
	}
	if _, err := ParseStyle("unknown"); err == nil {
		t.Fatal("ParseStyle accepted an unknown style")
	}
}

func TestPNGWriters(t *testing.T) {
	img, err := Render("writer", DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := WritePNG(&buffer, img); err != nil {
		t.Fatal(err)
	}
	decoded, err := png.Decode(&buffer)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Bounds() != img.Bounds() {
		t.Fatalf("bounds = %v, want %v", decoded.Bounds(), img.Bounds())
	}

	path := filepath.Join(t.TempDir(), "identicon.png")
	if err := SaveImage(img, path); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := png.Decode(file); err != nil {
		t.Fatal(err)
	}
}
