# IconForge

**Icons forged for every platform.**

Cross-platform application icon generator — converts SVG or PNG to production-ready icons (ICO, ICNS, PNG) and embeds them into Go binaries, macOS .app bundles, Linux .desktop entries, and Tauri/Electron projects.

## Features

- SVG/PNG to multi-resolution PNG rasterization (1024, 512, 256, 128, 64, 48, 32, 16)
- Windows ICO generation with `.syso` resource embedding (goversioninfo/rsrc compatible)
- macOS ICNS generation with `.app` bundle and `.iconset` directory support
- Linux `.desktop` entry generation with freedesktop hicolor icon theme
- Web favicon generation (ICO, Apple touch icon, PWA manifest)
- Automatic framework detection (Tauri, Electron, Wails, Fyne)
- Modern abstract gradient icon generation from color palette
- PNG input support (`--from-png` flag) as alternative to SVG
- File watching with auto-regeneration (`--watch` flag)

## Installation

```bash
go install github.com/inovacc/iconforge@latest
```

## Quick Start

```bash
# Generate icons from SVG for all platforms
iconforge forge --svg icon.svg --name myapp

# Generate a branded SVG icon + all platform icons
iconforge forge --generate --name myapp --primary "#4F46E5" --secondary "#7C3AED"

# Generate from existing PNG
iconforge forge --from-png logo.png --name myapp

# Rasterize SVG to PNGs only
iconforge render --svg icon.svg --sizes 512,256,128,64

# Generate web favicons
iconforge favicon --svg icon.svg -o build/favicons

# Auto-detect framework and generate appropriate icons
iconforge forge --svg icon.svg --name myapp --auto-detect

# Watch for changes and auto-rebuild
iconforge forge --svg icon.svg --name myapp --watch
```

## Commands

| Command | Description |
|---------|-------------|
| `forge` | Generate all icons from SVG/PNG source (main workflow) |
| `render` | Rasterize SVG/PNG to PNG at multiple sizes |
| `favicon` | Generate web favicons (ICO, Apple touch, PWA manifest) |
| `embed` | Generate .syso for Go build (Windows) |
| `detect` | Detect framework and show required assets |
| `version` | Print version information |
| `cmdtree` | Display command tree visualization |
| `aicontext` | Generate AI context documentation |

### Forge Flags

| Flag | Description |
|------|-------------|
| `--svg` | Path to source SVG file |
| `--from-png` | Path to source PNG file (alternative to `--svg`) |
| `--generate` | Generate a modern abstract gradient SVG icon |
| `--name` | Application name (auto-detected from directory if empty) |
| `--output, -o` | Output directory (default: `build/icons`) |
| `--primary` | Primary gradient color hex (default: `#4F46E5`) |
| `--secondary` | Secondary gradient color hex (default: `#7C3AED`) |
| `--accent` | Accent color hex (default: `#F59E0B`) |
| `--version` | Application version (default: `1.0.0`) |
| `--company` | Company name for Windows metadata |
| `--copyright` | Copyright notice |
| `--bundle-id` | macOS bundle identifier (e.g., `com.example.app`) |
| `--arch` | Target architecture for .syso (default: `amd64`) |
| `--skip-windows` | Skip Windows icon generation |
| `--skip-macos` | Skip macOS icon generation |
| `--skip-linux` | Skip Linux icon generation |
| `--auto-detect` | Auto-detect and generate framework-specific icons |
| `--favicon` | Also generate web favicons |
| `--iconset` | Also generate macOS .iconset directory |
| `--watch` | Watch source file and auto-regenerate on changes |

## Output Structure

```
build/icons/
  png/                    # Multi-resolution PNGs
    512x512.png
    256x256.png
    ...
  windows/
    icon.ico              # Multi-size ICO
    versioninfo.json      # goversioninfo metadata
    myapp.exe.manifest    # DPI-aware manifest
    rsrc_windows_amd64.syso  # Embedded resource (if rsrc available)
  macos/
    icon.icns             # Apple ICNS
    myapp.app/            # .app bundle with Info.plist
    myapp.iconset/        # .iconset directory (with --iconset)
  linux/
    myapp.desktop         # freedesktop .desktop file
    icons/hicolor/        # Icon theme structure
      48x48/apps/myapp.png
      256x256/apps/myapp.png
      ...
  favicon/                # Web favicons (with --favicon)
    favicon.ico
    favicon-16x16.png
    favicon-32x32.png
    apple-touch-icon.png
    android-chrome-192x192.png
    android-chrome-512x512.png
    site.webmanifest
```

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

# All checks
task check
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

BSD 3-Clause
