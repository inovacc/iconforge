# Roadmap

## Current Status
**Overall Progress:** 60% - Core functionality implemented

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

### Phase 5: Polish & Release [NOT STARTED]
- [ ] Unit tests (target: 80% coverage)
- [ ] Integration tests for forge pipeline
- [ ] CI/CD pipeline validation
- [ ] Documentation completion
- [ ] Performance optimization for large SVGs
- [ ] v1.0.0 release

## Test Coverage

**Current:** 0.0% | **Target:** 80%

| Package | Coverage | Status |
|---------|----------|--------|
| cmd | 0.0% | No tests |
| internal/svg | 0.0% | No tests |
| internal/icon | 0.0% | No tests |
| internal/platform | 0.0% | No tests |
| internal/detect | 0.0% | No tests |
| internal/generator | 0.0% | No tests |
