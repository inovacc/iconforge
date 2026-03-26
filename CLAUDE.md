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
- `pkg/winres` — Internalized Windows resource library (pure Go .syso generation, from tc-hib/winres v0.3.1)

## Commands

| Command | Description |
|---------|-------------|
| `forge` | Generate all icons from SVG (main workflow) |
| `forge --list-templates` | List available icon templates |
| `forge --template <name>` | Use a specific template with --generate |
| `forge --preview` | Show ANSI terminal preview of the generated icon |
| `forge --prompt` | Output AI prompt for interactive icon creation with Claude Code |
| `render` | Rasterize SVG to PNG at specified sizes |
| `embed` | Generate .syso resource file for Go builds |
| `favicon` | Generate web favicons (ICO, Apple touch, PWA) |
| `detect` | Detect framework and show required assets |
| `version` | Print version information |
| `cmdtree` | Display command tree visualization |
| `aicontext` | Generate AI context documentation |

## Icon Templates

10 built-in SVG templates (use `--template <name>` with `--generate`):

| Template | Description |
|----------|-------------|
| `forge` | Diamond facets with forge flame (default) |
| `shield` | Shield with keyhole — security, VPN, auth |
| `terminal` | Terminal window — CLI tools, dev utilities |
| `gear` | Mechanical gear — system utilities |
| `cube` | Isometric cube — data, 3D, containers |
| `bolt` | Lightning bolt — speed, performance |
| `leaf` | Leaf — eco, nature, growth |
| `wave` | Wave pattern — streaming, audio |
| `hexagon` | Hexagonal grid — science, tech |
| `stack` | Stacked layers — infrastructure, DevOps |

## Claude Code Integration

Run `iconforge forge --prompt` to get a structured AI prompt that guides Claude Code
through interactive icon creation (template selection, color picking, platform options).
Claude Code can then build and execute the full iconforge command, preview the result,
and iterate on the SVG if needed.

## Test Coverage

**Current:** ~87% | **Target:** 80% (met)

| Package | Coverage |
|---------|----------|
| internal/svg | 95.7% |
| internal/generator | 94.6% |
| internal/detect | 86.0% |
| internal/favicon | 85.3% |
| internal/platform | 83.4% |
| internal/icon | 82.0% |
| cmd | 80.3% |

## Icon Sizes

Standard: 512, 256, 128, 64, 48, 32, 16
Extended (with retina): 1024, 512, 256, 128, 64, 48, 32, 16

## Conventions

- SVG icons must use oksvg-compatible elements (no filters, no radialGradient, no text)
- ICO files store all sizes as PNG-compressed entries
- ICNS uses standard Apple OSType codes (icp4-ic10)
- Windows .syso files are auto-linked by the Go toolchain
