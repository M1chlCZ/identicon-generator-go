package main

import (
	"flag"
	"fmt"
	"github.com/M1chlCZ/identicon-generator-go/identicon"
	"os"
	"path/filepath"
)

const (
	outputDir   = "output"
	imagePrefix = "identicon"
)

func main() {
	// Parse command line flags
	config := identicon.DefaultConfig()
	count := 150 // Default count

	flag.IntVar(&config.Width, "w", config.Width, "Width of the image")
	flag.IntVar(&config.Width, "width", config.Width, "Width of the image")
	flag.IntVar(&config.Height, "h", config.Height, "Height of the image")
	flag.IntVar(&config.Height, "height", config.Height, "Height of the image")
	flag.IntVar(&config.GridSize, "g", config.GridSize, "Grid size (e.g., 5 for 5x5)")
	flag.IntVar(&config.GridSize, "grid", config.GridSize, "Grid size (e.g., 5 for 5x5)")
	flag.IntVar(&count, "c", count, "Number of images to generate")
	flag.IntVar(&count, "count", count, "Number of images to generate")
	flag.BoolVar(&config.Grayscale, "gray", config.Grayscale, "Generate grayscale images")
	flag.Parse()

	// Create output directory
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Generate images
	for i := 0; i < count; i++ {
		img := identicon.GenerateIdenticonWithConfig(i, config)
		filename := filepath.Join(outputDir, fmt.Sprintf("%s_%03d.png", imagePrefix, i+1))
		if err := identicon.SaveImage(img, filename); err != nil {
			fmt.Printf("Error saving image %s: %v\n", filename, err)
			continue
		}
		fmt.Printf("Generated image: %s\n", filename)
	}

	fmt.Println("Image generation complete.")
}
