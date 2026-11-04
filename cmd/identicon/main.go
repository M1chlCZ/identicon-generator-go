package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	identicon "github.com/M1chlCZ/identicon-generator-go"
)

const version = "0.3.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "identicon:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	config := identicon.DefaultConfig()
	var seed, outDir, prefix, styleName string
	var count, start int
	var showVersion bool

	flags := flag.NewFlagSet("identicon", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "usage: identicon [flags]")
		flags.PrintDefaults()
	}
	flags.StringVar(&seed, "seed", "", "render one identicon from this string seed")
	flags.IntVar(&count, "count", 150, "number of images in batch mode")
	flags.IntVar(&count, "c", 150, "shorthand for -count")
	flags.IntVar(&start, "start", 1, "first seed in batch mode")
	flags.StringVar(&outDir, "out", ".", "output directory")
	flags.StringVar(&prefix, "prefix", "identicon", "file name prefix")
	flags.StringVar(&styleName, "style", "symmetric", "symmetric, mosaic or rings")
	flags.IntVar(&config.Width, "width", config.Width, "image width in pixels")
	flags.IntVar(&config.Width, "w", config.Width, "shorthand for -width")
	flags.IntVar(&config.Height, "height", config.Height, "image height in pixels")
	flags.IntVar(&config.GridSize, "grid", config.GridSize, "grid size (positive)")
	flags.IntVar(&config.GridSize, "g", config.GridSize, "shorthand for -grid")
	flags.BoolVar(&config.Grayscale, "gray", config.Grayscale, "generate grayscale images")
	flags.BoolVar(&showVersion, "version", false, "print the version")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}
	if showVersion {
		fmt.Println("identicon", version)
		return nil
	}
	style, err := identicon.ParseStyle(styleName)
	if err != nil {
		return err
	}
	config.Style = style
	if err := config.Validate(); err != nil {
		return err
	}
	if prefix == "" || strings.ContainsAny(prefix, `/\`) {
		return errors.New("prefix must be a file name without path separators")
	}
	if start < 0 {
		return errors.New("start must not be negative")
	}

	if seed != "" {
		return writeOne(seed, filepath.Join(outDir, prefix+".png"), config)
	}
	if count < 1 {
		return errors.New("count must be at least 1")
	}
	for i := 0; i < count; i++ {
		value := start + i
		name := fmt.Sprintf("%s_%03d.png", prefix, value)
		if err := writeOne(strconv.Itoa(value), filepath.Join(outDir, name), config); err != nil {
			return err
		}
	}
	return nil
}

func writeOne(seed, path string, config identicon.Config) error {
	img, err := identicon.Render(seed, config)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := identicon.SaveImage(img, path); err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}
