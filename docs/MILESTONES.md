# Milestones

## v0.1.0 - Core Icon Pipeline
- **Target:** TBD
- **Status:** Complete
- **Goals:**
  - [x] SVG to multi-resolution PNG rasterization
  - [x] ICO encoder (Windows)
  - [x] ICNS encoder (macOS)
  - [x] PNG export at 8 standard sizes
  - [x] CLI commands: forge, render, embed, detect
- **Test Coverage:** 0.0% (no tests yet)

## v0.2.0 - Platform Packaging
- **Target:** TBD
- **Status:** Complete
- **Goals:**
  - [x] Windows versioninfo.json + .exe.manifest
  - [x] Windows .syso via rsrc
  - [x] macOS .app bundle with Info.plist
  - [x] Linux .desktop + hicolor icon theme
  - [x] Framework auto-detection (Tauri, Electron, Wails, Fyne)
- **Test Coverage:** 0.0% (no tests yet)

## v0.3.0 - SVG Icon Generator
- **Target:** TBD
- **Status:** Complete
- **Goals:**
  - [x] Abstract gradient SVG generation
  - [x] Configurable color palette
  - [x] oksvg-compatible output

## v0.4.0 - Testing & PNG Input
- **Target:** 2026-03-09
- **Status:** Complete
- **Goals:**
  - [x] Unit tests for internal/svg (95.7%)
  - [x] Unit tests for internal/detect (86.0%)
  - [x] Unit tests for internal/icon (69.0%)
  - [x] Unit tests for internal/platform (71.1%)
  - [x] PNG input support (--from-png flag)
- **Test Coverage:** ~65% (svg 95.7%, detect 86%, platform 71.1%, icon 69%, generator 0%, cmd 0%)

## v0.5.0 - Web & Polish
- **Target:** 2026-03-10
- **Status:** Complete
- **Goals:**
  - [x] Favicon generation (ICO, Apple touch, PWA manifest)
  - [x] macOS .iconset directory output
  - [x] Generator package tests (85.7%)
  - [x] Integration tests for forge pipeline (cmd 27.9%)
  - [x] --watch flag for auto-regeneration
  - [x] goversioninfo pure-Go .syso generation
  - [x] CI/CD pipeline fixed (Go 1.25)
  - [x] README documentation updated
- **Test Coverage:** ~75% (svg 95.7%, detect 86%, generator 85.7%, favicon 85.3%, platform 73.3%, icon 69%, cmd 27.9%)

## v1.0.0 - First Stable Release
- **Target:** TBD
- **Status:** In Progress
- **Goals:**
  - [ ] All internal packages at 80%+ coverage
  - [ ] cmd package at 50%+ coverage
  - [ ] Performance benchmarks
  - [ ] Tag and release v1.0.0
