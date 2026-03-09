# ADR-0002: Icon Generation Architecture

## Status
Accepted

## Context
Need to generate cross-platform application icons from SVG sources and embed them into Go binaries with proper metadata.

## Decision

### SVG Rasterization
- Use `oksvg` + `rasterx` for pure Go SVG to PNG rasterization (no CGo, no external tools)
- Generated SVGs must use only oksvg-compatible elements (linearGradient, basic shapes, paths — no filters, radialGradient, or text)

### Icon Formats
- **ICO** (Windows): Custom pure Go encoder storing PNG-compressed images per ICO spec
- **ICNS** (macOS): Custom pure Go encoder using standard Apple OSType codes
- **PNG** (Linux/universal): Standard `image/png` encoder

### Windows Embedding
- Generate `versioninfo.json` compatible with goversioninfo
- Generate `.syso` using `rsrc` (auto-linked by Go toolchain)
- Generate Windows application manifest (DPI-aware, supported OS declarations)

### macOS Packaging
- Generate `.app` bundle structure (Contents/MacOS/, Contents/Resources/)
- Generate `Info.plist` with proper CFBundleIconFile reference
- Place `.icns` in Contents/Resources/

### Linux Integration
- Generate `.desktop` file following freedesktop specification
- Install PNGs into `hicolor/{size}x{size}/apps/` theme structure

### Framework Detection
- Auto-detect Tauri (tauri.conf.json, src-tauri/)
- Auto-detect Electron (package.json with electron dependency)
- Auto-detect Wails (wails.json)
- Auto-detect Fyne (go.mod with fyne.io/fyne)
- Generate framework-specific icon assets automatically

## Consequences

### Positive
- Pure Go implementation — no external dependencies for core functionality
- Single command generates all platform assets
- Compatible with goversioninfo/rsrc ecosystem
- Framework-aware for Tauri/Electron/Wails/Fyne projects

### Negative
- SVG support limited to oksvg capabilities (no filters, no text rendering)
- ICO/ICNS encoders are custom implementations (not using established libraries)
- rsrc/goversioninfo must be installed separately for .syso generation
