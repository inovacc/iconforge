# Feature Requests

## Completed Features

### SVG to Multi-Resolution PNG
- **Status:** Completed
- **Description:** Rasterize SVG files to PNG at 8 sizes (1024, 512, 256, 128, 64, 48, 32, 16)
- **Implementation:** `internal/svg` using oksvg/rasterx

### Windows ICO Generation
- **Status:** Completed
- **Description:** Generate multi-size ICO files with PNG-compressed entries
- **Implementation:** `internal/icon/ico.go`

### macOS ICNS Generation
- **Status:** Completed
- **Description:** Generate ICNS files with standard Apple OSType codes
- **Implementation:** `internal/icon/icns.go`

### Windows Resource Embedding
- **Status:** Completed
- **Description:** Generate versioninfo.json, .exe.manifest, and .syso files
- **Implementation:** `internal/platform/windows.go`

### macOS .app Bundle
- **Status:** Completed
- **Description:** Create complete .app bundle with Info.plist and ICNS
- **Implementation:** `internal/platform/darwin.go`

### Linux Desktop Integration
- **Status:** Completed
- **Description:** Generate .desktop file and hicolor icon theme structure
- **Implementation:** `internal/platform/linux.go`

### Framework Auto-Detection
- **Status:** Completed
- **Description:** Detect Tauri, Electron, Wails, and Fyne frameworks and generate their required icon assets
- **Implementation:** `internal/detect/`

### SVG Icon Generator
- **Status:** Completed
- **Description:** Generate modern abstract gradient SVG icons with configurable color palette
- **Implementation:** `internal/generator/icon_svg.go`

## Proposed Features

### Favicon Generation
- **Priority:** P3
- **Status:** Proposed
- **Description:** Generate web favicon set (ICO, Apple touch, PWA manifest icons)
- **Motivation:** Many Go projects serve web UIs

### PNG Input Support
- **Priority:** P3
- **Status:** Proposed
- **Description:** Accept PNG as input (resize/convert) instead of requiring SVG
- **Motivation:** Not all projects have SVG source icons
