# Go Identicon Generator 🎨

A simple and flexible library for generating GitHub-style identicon images in Go. Create unique, deterministic avatars for your users or projects.

![Example Identicons](https://github.blog/wp-content/uploads/2013/08/a3c4e2a0-04df-11e3-824c-7378e6550707.png?fit=2384%2C784)

## Features

- ✨ Generate unique identicons from any integer seed
- 🎨 Support for both color and grayscale output
- 🔧 Customizable image size and grid complexity
- 💪 Zero external dependencies
- 🚀 Simple, clean API

## Installation

```bash
go get github.com/M1chlCZ/identicon-generator-go
```

## Quick Start

```go
package main

import (
    "github.com/M1chlCZ/identicon-generator-go"
)

func main() {
    // Generate with default settings
    img := identicon.GenerateIdenticon(123)
    identicon.SaveImage(img, "simple.png")

    // Custom configuration
    config := identicon.DefaultConfig()
    config.Width = 800
    config.Height = 800
    config.GridSize = 7
    config.Grayscale = true

    img = identicon.GenerateIdenticonWithConfig(123, config)
    identicon.SaveImage(img, "custom.png")
}
```

## CLI Tool

The repository includes a command-line tool for batch generating identicons:

```bash
go run cmd/main.go -w 500 -h 500 -grid 6 -count 10 -gray
```

### CLI Options

- `-w`, `-width`: Image width (default: 420)
- `-h`, `-height`: Image height (default: 420)
- `-g`, `-grid`: Grid size (default: 5)
- `-c`, `-count`: Number of images to generate (default: 150)
- `-gray`: Generate grayscale images (default: false)

## Configuration

```go
type Config struct {
    Width     int  // Width of the generated image in pixels
    Height    int  // Height of the generated image in pixels
    GridSize  int  // Complexity of the pattern (e.g., 5 for 5x5 grid)
    Grayscale bool // Whether to generate grayscale images
}
```

## Support

If this project helps you, you can support its development.

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/m1chl)
<a href="https://buymeacoffee.com/xdujifyzvt"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy me a coffee" height="30"></a>

## License

MIT License - feel free to use this in your own projects!

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

---
Made with ❤️ in Go