# IconForge

**Icons forged for every platform.**

Cross-platform application icon generator — converts SVG to production-ready icons (ICO, ICNS, PNG) and embeds them into Go binaries, macOS .app bundles, Linux .desktop entries, and Tauri/Electron projects.

## Features

- SVG to multi-resolution PNG rasterization (512, 256, 128, 64, 48, 32, 16)
- Windows ICO generation with `.syso` resource embedding (goversioninfo/rsrc compatible)
- macOS ICNS generation with `.app` bundle support
- Linux `.desktop` entry generation with freedesktop icon installation
- Automatic framework detection (Tauri, Electron, Wails, Fyne)
- Modern abstract gradient icon generation from color palette
- `versioninfo.json` manifest generation for Windows metadata

## Installation

```bash
go install github.com/inovacc/iconforge@latest
```

## Usage

```bash
iconforge --help
```

## Commands

| Command | Description |
|---------|-------------|
| `forge` | Generate all icons from SVG source (main workflow) |
| `render` | Rasterize SVG to PNG at multiple sizes |
| `embed` | Generate .syso for Go build (Windows) |
| `detect` | Detect framework and show required assets |
| `version` | Print version information |
| `cmdtree` | Display command tree visualization |
| `aicontext` | Generate AI context documentation |

## Development

```bash
# Build
task build

# Run
task run

# Test
task test

# Lint
task lint
```

## Release

```bash
# Create a snapshot release
task release:snapshot

# Create a production release (requires git tag)
git tag v1.0.0
task release
```

## License

MIT
