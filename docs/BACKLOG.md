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

_No P2 items._

### P3 - Future

- **Priority:** P3
- **Description:** Sixel/iTerm2 inline image protocol support for higher-fidelity terminal preview
- **Effort:** Medium
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
| P2: Internalize winres | Replaced goversioninfo with internalized tc-hib/winres v0.3.1 in pkg/winres/ (ADR-0003) | 2026-03-26 |
| P2: cmd/ test coverage | Boosted from 77.9% to 80.3% with 15+ new tests (templates, preview, Linux/macOS/Windows paths, flags) | 2026-03-26 |
| P3: Terminal icon preview | Added --preview flag with ANSI half-block rendering (works in all terminals) | 2026-03-26 |
| P3: Icon template system | 10 built-in SVG templates with registry, wired to --template/--list-templates flags | 2026-03-26 |
| P2: --watch flag | fsnotify-based file watching in forge command | 2026-03-10 |
