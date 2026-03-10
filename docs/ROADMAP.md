# Roadmap

## Current Status
**Overall Progress:** 75% - Core complete, tests in progress

## Phases

### Phase 1: Foundation [COMPLETE]
- [x] Project scaffold (structure, tooling, CI config)
- [x] CLI framework (Cobra commands: forge, render, embed, detect)
- [x] SVG rasterization engine (oksvg/rasterx)
- [x] Multi-resolution PNG export (1024, 512, 256, 128, 64, 48, 32, 16)

### Phase 2: Platform Icon Generators [COMPLETE]
- [x] ICO encoder (Windows, PNG-compressed multi-size)
- [x] ICNS encoder (macOS, Apple OSType codes)
- [x] Windows versioninfo.json (goversioninfo compatible)
- [x] Windows .exe.manifest generation (DPI-aware)
- [x] Windows .syso generation via rsrc
- [x] macOS .app bundle with Info.plist
- [x] Linux .desktop file (freedesktop spec)
- [x] Linux hicolor icon theme structure

### Phase 3: Framework Integration [COMPLETE]
- [x] Framework auto-detection (Tauri, Electron, Wails, Fyne)
- [x] Tauri icon asset generation (src-tauri/icons/)
- [x] Electron icon asset generation (build/)
- [x] Wails icon asset generation (build/)
- [x] Fyne icon asset generation (Icon.png)

### Phase 4: Icon Generation [COMPLETE]
- [x] Abstract gradient SVG icon generator
- [x] Configurable color palette (primary, secondary, accent)
- [x] oksvg-compatible SVG output

### Phase 5: Polish & Release [IN PROGRESS]
- [x] Unit tests for internal/svg (95.7% coverage)
- [x] Unit tests for internal/detect (86.0% coverage)
- [x] Unit tests for internal/icon (69.0% coverage)
- [x] Unit tests for internal/platform (73.3% coverage)
- [x] Unit tests for internal/generator (85.7% coverage)
- [x] Unit tests for internal/favicon (85.3% coverage)
- [x] PNG input support (--from-png flag)
- [x] Favicon generation (web-standard outputs)
- [x] macOS .iconset directory output
- [x] Integration tests for forge pipeline (cmd 27.9%)
- [ ] CI/CD pipeline validation
- [ ] Documentation completion
- [ ] Performance optimization for large SVGs
- [ ] v1.0.0 release

## Test Coverage

**Current:** ~80% | **Target:** 80%

| Package | Coverage | Status |
|---------|----------|--------|
| internal/svg | 95.7% | Excellent |
| internal/platform | 94.8% | Excellent |
| internal/detect | 86.0% | Good |
| internal/generator | 85.7% | Good |
| internal/favicon | 85.3% | Good |
| internal/icon | 82.0% | Good |
| cmd | 26.4% | Needs improvement |
