# IconForge - CLAUDE.md

## Overview

Cross-platform application icon generator. Converts SVG to production-ready ICO, ICNS, and PNG icons, and embeds them into Go binaries, macOS .app bundles, Linux .desktop entries, and Tauri/Electron projects.

## Build & Test

```bash
task build          # Build the binary
task test           # Run tests with coverage
task lint           # Run golangci-lint
task check          # Run all quality checks
go run . --help     # Run directly
```

## Project Structure

```
cmd/                    # CLI commands (Cobra)
  forge.go              # Main forge command (generates all icons)
  render.go             # SVG to PNG rasterization
  embed.go              # Windows .syso generation
  detect.go             # Framework detection
internal/
  svg/                  # SVG parsing and rasterization (oksvg/rasterx)
  icon/                 # ICO, ICNS, PNG encoders
  platform/             # Platform-specific generators (Windows, macOS, Linux)
  detect/               # Framework detection (Tauri, Electron, Wails, Fyne)
  favicon/              # Web favicon generation (ICO, Apple touch, PWA manifest)
  generator/            # SVG icon generator
docs/                   # Documentation
```

## Key Dependencies

- `github.com/srwiley/oksvg` + `rasterx` — Pure Go SVG rasterizer
- `github.com/spf13/cobra` — CLI framework
- External tools (optional): `rsrc`, `goversioninfo` for .syso generation

## Commands

| Command | Description |
|---------|-------------|
| `forge` | Generate all icons from SVG (main workflow) |
| `render` | Rasterize SVG to PNG at specified sizes |
| `embed` | Generate .syso resource file for Go builds |
| `favicon` | Generate web favicons (ICO, Apple touch, PWA) |
| `detect` | Detect framework and show required assets |
| `version` | Print version information |
| `cmdtree` | Display command tree visualization |
| `aicontext` | Generate AI context documentation |

## Test Coverage

**Current:** ~75% | **Target:** 80%

| Package | Coverage |
|---------|----------|
| internal/svg | 95.7% |
| internal/detect | 86.0% |
| internal/generator | 85.7% |
| internal/favicon | 85.3% |
| internal/platform | 73.3% |
| internal/icon | 69.0% |
| cmd | 27.9% |

## Icon Sizes

Standard: 512, 256, 128, 64, 48, 32, 16
Extended (with retina): 1024, 512, 256, 128, 64, 48, 32, 16

## Conventions

- SVG icons must use oksvg-compatible elements (no filters, no radialGradient, no text)
- ICO files store all sizes as PNG-compressed entries
- ICNS uses standard Apple OSType codes (icp4-ic10)
- Windows .syso files are auto-linked by the Go toolchain
