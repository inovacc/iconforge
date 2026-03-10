# Backlog

## Priority Levels

| Priority | Timeline |
|----------|----------|
| P1 | This sprint |
| P2 | This quarter |
| P3 | Future |

## Items

### P1 - This Sprint

- **Priority:** P1
- **Description:** Add unit tests for internal/generator (GenerateIconSVG)
- **Effort:** Small
- **Category:** Tech Debt

- **Priority:** P1
- **Description:** Integration tests for forge/render pipeline (cmd/ package)
- **Effort:** Medium
- **Category:** Tech Debt

### P2 - This Quarter

- **Priority:** P2
- **Description:** Add goversioninfo support as alternative to rsrc for .syso generation
- **Effort:** Medium
- **Category:** Feature

- **Priority:** P2
- **Description:** Add `--watch` flag to forge command for auto-regeneration on SVG changes
- **Effort:** Medium
- **Category:** Feature

### P3 - Future

- **Priority:** P3
- **Description:** Interactive icon preview in terminal (sixel/iTerm2 inline images)
- **Effort:** Large
- **Category:** Feature

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
