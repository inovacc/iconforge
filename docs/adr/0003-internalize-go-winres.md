# ADR-0003: Internalize go-winres

**Status:** Accepted (implemented Option B — internalized `winres` library)
**Date:** 2026-03-26

## Context

Evaluated `github.com/tc-hib/go-winres` for internalization into iconforge, which currently uses `github.com/josephspurrier/goversioninfo` for Windows .syso resource generation.

## Summary

| Field | Value |
|-------|-------|
| Module | `github.com/tc-hib/go-winres` |
| Version | v0.3.3 |
| Upstream Commit | `c4d55a3dfc2e22a4a9fb95ae0280ce106daadce9` |
| Total Go files | 6 (4 source + 2 test) |
| Total LOC | ~2,804 (includes ~2,000 lines of embedded icon data) |
| Packages | 1 (`main`) |
| Direct deps | 2 (`github.com/tc-hib/winres`, `github.com/urfave/cli/v2`) |
| License | ISC (0-clause BSD equivalent — very permissive) |

## Structure Analysis

```
.tmp/go-winres/
├── main.go                    # CLI entry point, commands (init, make, simply, extract, patch)
├── resdir.go                  # JSON resource import/export, icon/cursor/bitmap helpers
├── resdir_test.go             # Tests for resource directory
├── main_test.go               # Tests for CLI commands
├── defaulticon_go1.16.go      # Embedded default icon (go:embed, ~60K tokens)
└── defaulticon_other.go       # Fallback for pre-1.16 Go
```

- **Single package:** `package main` — **not importable as a library**
- No `internal/`, `pkg/`, or sub-packages
- No CGO, no assembly, no code generation

## Dependency Graph

```
go-winres (CLI)
├── github.com/tc-hib/winres v0.3.1        ← THE ACTUAL LIBRARY
│   ├── github.com/nfnt/resize              (image resizing)
│   └── golang.org/x/image                  (image codecs)
└── github.com/urfave/cli/v2 v2.27.1        ← CLI framework (not needed)
    ├── github.com/cpuguy83/go-md2man/v2
    ├── github.com/russross/blackfriday/v2
    └── github.com/xrash/smetrics
```

## Exported API

**None.** `go-winres` is `package main` — it exports nothing. All useful functionality comes from `github.com/tc-hib/winres`, which provides:

- `winres.ResourceSet` — core resource container
- `winres.ResourceSet.WriteObject()` — generates .syso/.obj files
- `winres.ResourceSet.SetIcon()` / `SetManifest()` / `SetVersionInfo()` — resource setters
- `winres.LoadFromEXE()` / `WriteToEXE()` — extract/patch EXE resources
- `winres.NewIconFromResizedImage()` — icon from PNG with auto-resize
- `winres.LoadICO()` — load .ico files
- `winres.Arch` — target architecture constants (ArchAMD64, ArchI386, ArchARM, ArchARM64)

## Current Project Usage

iconforge does **not** use `go-winres` or `winres`. It currently uses:

| Current Dependency | Usage |
|--------------------|-------|
| `github.com/josephspurrier/goversioninfo` | `internal/platform/windows.go` — .syso generation via `VersionInfo.WriteSyso()` |
| `rsrc` (external CLI) | Fallback for .syso generation via `exec.Command("rsrc", ...)` |

### Call sites in `internal/platform/windows.go`:
- `GenerateSysoGoversioninfo()` — primary: uses `goversioninfo.VersionInfo` struct directly
- `GenerateSysoFromJSON()` — secondary: parses versioninfo.json → .syso
- `GenerateSyso()` — fallback: shells out to `rsrc` CLI
- `GenerateSysoWithGovernVersionInfo()` — deprecated: shells out to `goversioninfo` CLI

## Packages to Internalize

### From `go-winres` itself: **None recommended**

The `go-winres` CLI is not importable. The useful patterns (JSON resource definitions, icon loading, version info handling) are ~300 LOC of glue code that could be adapted, but they depend entirely on the `winres` library.

### What should actually be internalized: `github.com/tc-hib/winres`

If the goal is to replace `goversioninfo` with a better .syso generator, the target should be `winres` — the library that `go-winres` wraps. It provides:

1. **Pure Go .syso generation** (no external tools needed)
2. **Multi-arch support** (amd64, 386, arm, arm64) — goversioninfo only does amd64
3. **Full resource types** (icons, cursors, bitmaps, manifests, version info, arbitrary data)
4. **EXE patching** (extract/replace resources in existing .exe files)
5. **Icon auto-resizing** from a single PNG source

## Import Rewrite Map

Not applicable — `go-winres` is not imported by iconforge.

If `winres` were internalized instead:
```
github.com/tc-hib/winres         → github.com/inovacc/iconforge/pkg/winres
github.com/tc-hib/winres/version → github.com/inovacc/iconforge/pkg/winres/version
```

## Minimal Subset Analysis

From `go-winres` CLI, the patterns worth adapting (~300 LOC):
- **Used:** `writeObjectFile()`, `simplySetIcon()`, `simplySetVersionInfo()`, `simplySetManifest()` patterns
- **Not needed:** `cmdInit`, `cmdExtract`, `cmdPatch`, all `urfave/cli` integration, embedded default icons, JSON resource import/export (unless iconforge needs custom resource definitions)

## Risk Assessment

### License Compatibility
- **go-winres**: ISC license — compatible with BSD-3 (iconforge's license) ✅
- **winres**: ISC license — compatible ✅

### Maintenance Burden
- `go-winres` last commit: active, v0.3.3
- `winres` is actively maintained by the same author
- Internalizing freezes the version — acceptable for a stable API

### Complexity
- No CGO, no assembly, no build tags (except `go1.16` for embed)
- No code generation
- Low complexity ✅

### Breaking Changes Risk
- Freezing at v0.3.x is acceptable — the API is stable
- `winres` has had no breaking changes since v0.2.0

## Recommended Strategy

### **Keep external** (for `go-winres`)

**Justification:** `go-winres` is a CLI tool (`package main`), not a library. It cannot be imported. Copying it into `pkg/go-winres/` would require renaming the package and stripping the CLI framework — essentially rewriting it.

### **Alternative recommendation: Add `winres` as a dependency (or internalize it)**

Instead of internalizing `go-winres`, consider one of these approaches:

#### Option A: Replace `goversioninfo` with `winres` as an external dependency
```bash
go get github.com/tc-hib/winres@v0.3.1
```
Then refactor `internal/platform/windows.go` to use `winres.ResourceSet` instead of `goversioninfo.VersionInfo`. Benefits:
- Multi-arch .syso generation (arm, arm64 in addition to amd64, 386)
- Richer resource support (cursors, bitmaps, manifest objects)
- No external tool fallback needed (`rsrc` becomes unnecessary)
- Removes `goversioninfo` + `rsrc` dependencies

#### Option B: Internalize `winres` library into `pkg/winres/`
If you want zero external deps for the core functionality, internalize `winres` instead:
- ~15 Go files, ~3,000 LOC (excluding tests)
- 2 sub-packages: `winres` + `winres/version`
- 2 transitive deps: `nfnt/resize`, `golang.org/x/image`

## Step-by-Step Execution Plan (for Option A — recommended)

1. `go get github.com/tc-hib/winres@v0.3.1`
2. Refactor `internal/platform/windows.go`:
   - Replace `goversioninfo.VersionInfo` with `winres.ResourceSet`
   - Use `rs.SetIcon()` + `rs.SetManifest()` + `rs.SetVersionInfo()` + `rs.WriteObject()`
   - Support multiple architectures via `winres.Arch` constants
3. Remove `GenerateSyso()` (rsrc fallback) and `GenerateSysoWithGovernVersionInfo()` (deprecated)
4. Update `cmd/embed.go` to use the new `winres`-based functions
5. Remove old dependencies: `go get -u && go mod tidy`
   - Removes: `goversioninfo`, `akavel/rsrc`, `cli/safeexec`
6. Run tests: `go test ./...`
7. Update CLAUDE.md to reflect the dependency change

## Decision

**Did not internalize `go-winres`** (CLI tool, `package main`). Instead, **internalized `github.com/tc-hib/winres` v0.3.1** (the library) into `pkg/winres/`.

### What was done:
1. Copied 17 source files (12 root + 5 version/) into `pkg/winres/`
2. Rewrote all imports from `github.com/tc-hib/winres` → `github.com/inovacc/iconforge/pkg/winres`
3. Added `github.com/nfnt/resize` as transitive dependency
4. Refactored `internal/platform/windows.go`:
   - Replaced `goversioninfo.VersionInfo` with `winres.ResourceSet` + `version.Info`
   - New `GenerateSysoWinres()` replaces `GenerateSysoGoversioninfo()`
   - New `GenerateSysoFromICO()` replaces `GenerateSyso()` (rsrc fallback)
   - Removed `GenerateSysoFromJSON()` and `GenerateSysoWithGovernVersionInfo()` (deprecated)
5. Updated `cmd/embed.go` and `cmd/forge.go` to use new functions
6. Updated all tests — all passing
7. Removed dependencies: `josephspurrier/goversioninfo`, `akavel/rsrc`
8. Generated tracking file at `pkg/winres/.dep-track.json`
