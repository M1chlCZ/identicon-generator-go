# identicon-generator-go

A small Go library and command-line tool that generates deterministic,
GitHub-style identicon images. The same seed always produces the same image.
The module uses only the Go standard library.

## Install

```sh
go get github.com/M1chlCZ/identicon-generator-go@v0.2.0
```

## Library

`Render` accepts a string seed. It validates the configuration and returns an
error instead of panicking:

```go
package main

import (
	"log"

	identicon "github.com/M1chlCZ/identicon-generator-go"
)

func main() {
	config := identicon.DefaultConfig()
	config.Width = 800
	config.Height = 800
	config.GridSize = 7

	img, err := identicon.Render("alice@example.com", config)
	if err != nil {
		log.Fatal(err)
	}
	if err := identicon.SaveImage(img, "alice.png"); err != nil {
		log.Fatal(err)
	}
}
```

The seed is hashed with SHA-256. The output is deterministic across
architectures. `WritePNG` writes to any `io.Writer`:

```go
var buffer bytes.Buffer
if err := identicon.WritePNG(&buffer, img); err != nil {
	log.Fatal(err)
}
```

### Legacy API

`GenerateIdenticon` and `GenerateIdenticonWithConfig` keep the original
integer-seed API, MD5 hash, and pixel output:

```go
img := identicon.GenerateIdenticon(123)
identicon.SaveImage(img, "legacy.png")
```

The legacy functions panic on an invalid configuration. Use `Render` in new
code.

### Configuration

| Field | Default | Meaning |
| --- | --- | --- |
| `Width` | `420` | Image width in pixels |
| `Height` | `420` | Image height in pixels |
| `GridSize` | `5` | Grid dimension, at least 1 |
| `Grayscale` | `false` | Render gray cells only |

`Config.Validate` requires positive dimensions, a positive grid size, and at
least one pixel per grid cell.

## CLI

```sh
go run ./cmd/identicon -seed alice@example.com -out avatars
go run ./cmd/identicon -count 20 -start 1 -out avatars -width 500 -grid 6
```

| Flag | Default | Meaning |
| --- | --- | --- |
| `-seed` | empty | Render one identicon from a string seed |
| `-count`, `-c` | `150` | Number of images in batch mode |
| `-start` | `1` | First integer seed in batch mode |
| `-out` | `.` | Output directory |
| `-prefix` | `identicon` | File name prefix |
| `-width`, `-w` | `420` | Image width in pixels |
| `-height` | `420` | Image height in pixels |
| `-grid`, `-g` | `5` | Grid dimension |
| `-gray` | `false` | Render grayscale images |
| `-version` | `false` | Print the version |

Batch output uses the names `identicon_001.png`, `identicon_002.png`, and so
on. A single seed writes `identicon.png` unless `-prefix` changes the name.

## Development

```sh
go test -race ./...
go vet ./...
go build ./...
```

The modernization roadmap is in
[`docs/plans/2026-11-03-identicon-modernization.md`](docs/plans/2026-11-03-identicon-modernization.md).
It covers additional generator styles, SVG output, an HTTP avatar service, and
release binaries.

## Support

If this project helps you, you can support its development.

[![Support me on Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/m1chl)
<a href="https://buymeacoffee.com/xdujifyzvt"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy me a coffee" height="30"></a>

## License

MIT License. See [LICENSE](LICENSE).
