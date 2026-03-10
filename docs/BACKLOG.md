# Backlog

## Priority Levels

| Priority | Timeline |
|----------|----------|
| P1 | This sprint |
| P2 | This quarter |
| P3 | Future |

## Items

### P1 - This Sprint

_No P1 items._

### P2 - This Quarter

- **Priority:** P2
- **Description:** Add more cmd/ integration tests (currently 26.4%)
- **Effort:** Medium
- **Category:** Tech Debt

### P3 - Future

- **Priority:** P3
- **Description:** Interactive icon preview in terminal (sixel/iTerm2 inline images)
- **Effort:** Large
- **Category:** Feature

- **Priority:** P3
- **Description:** Performance benchmarks for large SVG rasterization
- **Effort:** Medium
- **Category:** Infrastructure

## Resolved

| Item | Resolution | Date |
|------|------------|------|
| P1: Unit tests for internal/icon | 69.0% coverage (ico_test, icns_test, png_test) | 2026-03-09 |
| P1: Unit tests for internal/svg | 95.7% coverage (render_test) | 2026-03-09 |
| P1: Unit tests for internal/detect | 86.0% coverage (framework_test, generate_test) | 2026-03-09 |
| P1: Unit tests for internal/platform | 71.1% coverage (windows_test, darwin_test, linux_test) | 2026-03-09 |
| P2: Customize aicontext.go | Updated with project-specific categories and structure | 2026-03-09 |
| P2: .iconset directory output | CreateIconset in platform/darwin.go | 2026-03-09 |
| P3: --from-png flag | Added to forge (--from-png) and render (--png) commands | 2026-03-09 |
| P3: Favicon generation | favicon command + internal/favicon package | 2026-03-09 |
| P1: Unit tests for internal/generator | 85.7% coverage (icon_svg_test) | 2026-03-10 |
| P1: Integration tests for cmd/ | 27.9% coverage (forge_test, render_test) | 2026-03-10 |
| P2: goversioninfo support | Pure-Go .syso generation in platform/windows.go | 2026-03-10 |
| P2: --watch flag | fsnotify-based file watching in forge command | 2026-03-10 |
