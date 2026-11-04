# Changelog

## v0.3.0

- Add `Style` to `Config`, with `StyleSymmetric`, `StyleMosaic` and
  `StyleRings`.
- Add `ParseStyle` and `Style.String` for configuration and CLI input.
- Add `-style` to the CLI.
- Keep the legacy integer-seed output unchanged.

## v0.2.0

- Add `Config.Validate`, `Render`, `WritePNG` and a fixed `SaveImage`.
- Move the CLI to `cmd/identicon` and add seed, output, prefix and version
  flags.
- Add tests, CI, Dependabot and package documentation.
- Bump the module to Go 1.25.0.
