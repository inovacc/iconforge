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
- **Description:** Add unit tests for internal/icon (ICO, ICNS, PNG encoders)
- **Effort:** Medium
- **Category:** Tech Debt

- **Priority:** P1
- **Description:** Add unit tests for internal/svg (RenderToImage, RenderToImages)
- **Effort:** Medium
- **Category:** Tech Debt

- **Priority:** P1
- **Description:** Add unit tests for internal/detect (DetectFramework, GenerateFrameworkIcons)
- **Effort:** Small
- **Category:** Tech Debt

- **Priority:** P1
- **Description:** Add unit tests for internal/platform (WriteVersionInfo, CreateAppBundle, CreateDesktopEntry)
- **Effort:** Medium
- **Category:** Tech Debt

### P2 - This Quarter

- **Priority:** P2
- **Description:** Customize aicontext.go TODO sections (categories, project structure) for IconForge
- **Effort:** Small
- **Category:** Tech Debt
- **Source:** `cmd/aicontext.go:130`, `cmd/aicontext.go:146`

- **Priority:** P2
- **Description:** Add goversioninfo support as alternative to rsrc for .syso generation
- **Effort:** Medium
- **Category:** Feature

- **Priority:** P2
- **Description:** Support `.iconset` directory output for macOS (for use with `iconutil`)
- **Effort:** Small
- **Category:** Feature

- **Priority:** P2
- **Description:** Add `--watch` flag to forge command for auto-regeneration on SVG changes
- **Effort:** Medium
- **Category:** Feature

### P3 - Future

- **Priority:** P3
- **Description:** Support favicon generation (16x16, 32x32 ICO + 180x180 Apple touch + 192/512 PWA manifest)
- **Effort:** Medium
- **Category:** Feature

- **Priority:** P3
- **Description:** Add `--from-png` flag to accept PNG input instead of SVG
- **Effort:** Small
- **Category:** Feature

- **Priority:** P3
- **Description:** Interactive icon preview in terminal (sixel/iTerm2 inline images)
- **Effort:** Large
- **Category:** Feature
