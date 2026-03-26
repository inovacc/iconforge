package cmd

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// testSVGContent is a minimal valid SVG for testing the forge pipeline.
const testSVGContent = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">
  <rect x="10" y="10" width="80" height="80" fill="#4F46E5"/>
  <circle cx="50" cy="50" r="30" fill="#F59E0B"/>
</svg>`

// resetForgeFlags resets all forge global flag variables to their defaults.
func resetForgeFlags(t *testing.T) {
	t.Helper()

	forgeSVGPath = ""
	forgePNGPath = ""
	forgeOutputDir = "build/icons"
	forgeAppName = ""
	forgeVersion = "1.0.0"
	forgeCompany = ""
	forgeCopyright = ""
	forgeBundleID = ""
	forgeArch = "amd64"
	forgeGenSVG = false
	forgePrimary = "#4F46E5"
	forgeSecondary = "#7C3AED"
	forgeAccent = "#F59E0B"
	forgeSkipWin = false
	forgeSkipMac = false
	forgeSkipLinux = false
	forgeAutoDetect = false
	forgeIconset = false
	forgeFavicon = false
	forgeWatch = false
	forgeTemplate = "forge"
	forgeListTemplates = false
	forgePrompt = false
	forgePreview = false
	forgePreviewSize = 32
}

// writeTestSVG writes a valid SVG file to disk and returns its path.
func writeTestSVG(t *testing.T, dir string) string {
	t.Helper()
	p := filepath.Join(dir, "test.svg")
	if err := os.WriteFile(p, []byte(testSVGContent), 0o644); err != nil {
		t.Fatalf("write test svg: %v", err)
	}
	return p
}

// writeTestPNG creates a valid PNG file on disk and returns its path.
func writeTestPNG(t *testing.T, dir string) string {
	t.Helper()
	p := filepath.Join(dir, "test.png")
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("create test png: %v", err)
	}
	defer func() { _ = f.Close() }()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode test png: %v", err)
	}
	return p
}

// executeForge sets up the forgeCmd output buffers and calls runForge directly.
func executeForge(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	forgeCmd.SetOut(buf)
	forgeCmd.SetErr(buf)
	err := runForge(forgeCmd, nil)
	return buf.String(), err
}

func TestForge_GenerateSVG(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "testapp"
	forgeOutputDir = tmpDir

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --generate failed: %v", err)
	}

	// Verify SVG was created
	svgPath := filepath.Join(tmpDir, "testapp.svg")
	if _, err := os.Stat(svgPath); os.IsNotExist(err) {
		t.Error("expected generated SVG file to exist")
	}

	// Verify PNGs directory
	pngDir := filepath.Join(tmpDir, "png")
	if _, err := os.Stat(pngDir); os.IsNotExist(err) {
		t.Error("expected png/ directory to exist")
	}

	// Verify Windows ICO
	icoPath := filepath.Join(tmpDir, "windows", "icon.ico")
	if _, err := os.Stat(icoPath); os.IsNotExist(err) {
		t.Error("expected windows/icon.ico to exist")
	}

	// Verify macOS ICNS
	icnsPath := filepath.Join(tmpDir, "macos", "icon.icns")
	if _, err := os.Stat(icnsPath); os.IsNotExist(err) {
		t.Error("expected macos/icon.icns to exist")
	}

	// Verify Linux .desktop
	entries, err := os.ReadDir(filepath.Join(tmpDir, "linux"))
	if err != nil {
		t.Error("expected linux/ directory to exist")
	} else {
		found := false
		for _, e := range entries {
			if filepath.Ext(e.Name()) == ".desktop" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected .desktop file in linux/ directory")
		}
	}
}

func TestForge_FromSVG(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()
	svgDir := t.TempDir()
	svgPath := writeTestSVG(t, svgDir)

	forgeSVGPath = svgPath
	forgeAppName = "svgtest"
	forgeOutputDir = tmpDir

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --svg failed: %v", err)
	}

	// Verify PNGs
	pngDir := filepath.Join(tmpDir, "png")
	if _, err := os.Stat(pngDir); os.IsNotExist(err) {
		t.Error("expected png/ directory to exist")
	}

	// Verify Windows
	if _, err := os.Stat(filepath.Join(tmpDir, "windows", "icon.ico")); os.IsNotExist(err) {
		t.Error("expected windows/icon.ico to exist")
	}

	// Verify macOS
	if _, err := os.Stat(filepath.Join(tmpDir, "macos", "icon.icns")); os.IsNotExist(err) {
		t.Error("expected macos/icon.icns to exist")
	}

	// Verify Linux
	if _, err := os.Stat(filepath.Join(tmpDir, "linux")); os.IsNotExist(err) {
		t.Error("expected linux/ directory to exist")
	}
}

func TestForge_FromPNG(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()
	pngDir := t.TempDir()
	pngPath := writeTestPNG(t, pngDir)

	forgePNGPath = pngPath
	forgeAppName = "pngtest"
	forgeOutputDir = tmpDir

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --from-png failed: %v", err)
	}

	// Verify PNGs
	if _, err := os.Stat(filepath.Join(tmpDir, "png")); os.IsNotExist(err) {
		t.Error("expected png/ directory to exist")
	}

	// Verify Windows
	if _, err := os.Stat(filepath.Join(tmpDir, "windows", "icon.ico")); os.IsNotExist(err) {
		t.Error("expected windows/icon.ico to exist")
	}

	// Verify macOS
	if _, err := os.Stat(filepath.Join(tmpDir, "macos", "icon.icns")); os.IsNotExist(err) {
		t.Error("expected macos/icon.icns to exist")
	}

	// Verify Linux
	if _, err := os.Stat(filepath.Join(tmpDir, "linux")); os.IsNotExist(err) {
		t.Error("expected linux/ directory to exist")
	}
}

func TestForge_SkipPlatforms(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "skiptest"
	forgeOutputDir = tmpDir
	forgeSkipWin = true
	forgeSkipLinux = true

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --skip-windows --skip-linux failed: %v", err)
	}

	// macOS should exist
	if _, err := os.Stat(filepath.Join(tmpDir, "macos", "icon.icns")); os.IsNotExist(err) {
		t.Error("expected macos/icon.icns to exist when not skipped")
	}

	// Windows should NOT exist
	if _, err := os.Stat(filepath.Join(tmpDir, "windows")); !os.IsNotExist(err) {
		t.Error("expected windows/ directory to not exist when --skip-windows")
	}

	// Linux should NOT exist
	if _, err := os.Stat(filepath.Join(tmpDir, "linux")); !os.IsNotExist(err) {
		t.Error("expected linux/ directory to not exist when --skip-linux")
	}
}

func TestForge_MutuallyExclusive(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr string
	}{
		{
			name: "svg and from-png",
			setup: func() {
				forgeSVGPath = "some.svg"
				forgePNGPath = "some.png"
			},
			wantErr: "--from-png and --svg are mutually exclusive",
		},
		{
			name: "from-png and generate",
			setup: func() {
				forgePNGPath = "some.png"
				forgeGenSVG = true
			},
			wantErr: "--from-png and --generate are mutually exclusive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetForgeFlags(t)
			tmpDir := t.TempDir()
			forgeOutputDir = tmpDir
			forgeAppName = "test"
			tt.setup()

			_, err := executeForge(t)
			if err == nil {
				t.Fatal("expected error for mutually exclusive flags, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestForge_NoSource(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()
	forgeOutputDir = tmpDir
	forgeAppName = "test"

	_, err := executeForge(t)
	if err == nil {
		t.Fatal("expected error when no source specified, got nil")
	}
	expected := "no source specified; use --svg, --from-png, or --generate"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestForge_CustomColors(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "colortest"
	forgeOutputDir = tmpDir
	forgePrimary = "#FF0000"
	forgeSecondary = "#00FF00"
	forgeAccent = "#0000FF"

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge with custom colors failed: %v", err)
	}

	// Verify SVG was created and contains custom colors
	svgPath := filepath.Join(tmpDir, "colortest.svg")
	data, err := os.ReadFile(svgPath)
	if err != nil {
		t.Fatalf("read generated SVG: %v", err)
	}
	svgStr := string(data)
	for _, c := range []string{"#FF0000", "#00FF00", "#0000FF"} {
		if !containsStr(svgStr, c) {
			t.Errorf("generated SVG does not contain custom color %s", c)
		}
	}
}

func TestForge_AutoDetect(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
	})

	// Create the tauri marker file in a project directory
	projectDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(projectDir, "tauri.conf.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write tauri.conf.json: %v", err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	forgeGenSVG = true
	forgeAppName = "tauriapp"
	forgeOutputDir = tmpDir
	forgeAutoDetect = true

	output, execErr := executeForge(t)
	if execErr != nil {
		t.Fatalf("forge --auto-detect failed: %v", execErr)
	}

	if !containsStr(output, "Detected framework: Tauri") {
		t.Errorf("expected output to contain framework detection message, got: %s", output)
	}
}

func TestForge_WithFavicon(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "favicontest"
	forgeOutputDir = tmpDir
	forgeFavicon = true

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --generate --favicon failed: %v", err)
	}

	// Verify favicon directory was created
	faviconDir := filepath.Join(tmpDir, "favicon")
	if _, err := os.Stat(faviconDir); os.IsNotExist(err) {
		t.Error("expected favicon/ directory to exist")
	}
}

func TestForge_WithIconset(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "iconsettest"
	forgeOutputDir = tmpDir
	forgeIconset = true

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --generate --iconset failed: %v", err)
	}

	// macOS iconset should exist
	entries, err := os.ReadDir(filepath.Join(tmpDir, "macos"))
	if err != nil {
		t.Fatal("expected macos/ directory to exist")
	}

	foundIconset := false
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".iconset" || containsStr(e.Name(), "iconset") {
			foundIconset = true
			break
		}
	}
	if !foundIconset {
		t.Error("expected .iconset directory in macos/")
	}
}

func TestForge_SkipAllPlatforms(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "skipall"
	forgeOutputDir = tmpDir
	forgeSkipWin = true
	forgeSkipMac = true
	forgeSkipLinux = true

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --skip-all failed: %v", err)
	}

	// PNGs should still exist
	if _, err := os.Stat(filepath.Join(tmpDir, "png")); os.IsNotExist(err) {
		t.Error("expected png/ directory to exist even when all platforms skipped")
	}

	// Platform dirs should NOT exist
	for _, dir := range []string{"windows", "macos", "linux"} {
		if _, err := os.Stat(filepath.Join(tmpDir, dir)); !os.IsNotExist(err) {
			t.Errorf("expected %s/ directory to not exist when skipped", dir)
		}
	}
}

func TestForge_InvalidSVGPath(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeSVGPath = filepath.Join(tmpDir, "nonexistent.svg")
	forgeAppName = "test"
	forgeOutputDir = tmpDir

	_, err := executeForge(t)
	if err == nil {
		t.Fatal("expected error for nonexistent SVG, got nil")
	}
}

func TestForge_InvalidPNGPath(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgePNGPath = filepath.Join(tmpDir, "nonexistent.png")
	forgeAppName = "test"
	forgeOutputDir = tmpDir

	_, err := executeForge(t)
	if err == nil {
		t.Fatal("expected error for nonexistent PNG, got nil")
	}
}

func TestForge_ListTemplates(t *testing.T) {
	resetForgeFlags(t)
	forgeListTemplates = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --list-templates failed: %v", err)
	}

	expectedTemplates := []string{"forge", "shield", "terminal", "gear", "cube", "bolt", "leaf", "wave", "hexagon", "stack"}
	for _, name := range expectedTemplates {
		if !containsStr(output, name) {
			t.Errorf("--list-templates output missing template %q", name)
		}
	}
	if !containsStr(output, "Available icon templates") {
		t.Error("expected header 'Available icon templates' in output")
	}
}

func TestForge_Prompt(t *testing.T) {
	resetForgeFlags(t)
	forgePrompt = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --prompt failed: %v", err)
	}

	if !containsStr(output, "iconforge") {
		t.Error("expected 'iconforge' in prompt output")
	}
	if !containsStr(output, "--template") {
		t.Error("expected '--template' in prompt output")
	}
}

func TestForge_GenerateWithTemplate(t *testing.T) {
	templates := []string{"shield", "terminal", "gear", "cube", "bolt", "leaf", "wave", "hexagon", "stack"}

	for _, tmpl := range templates {
		t.Run(tmpl, func(t *testing.T) {
			resetForgeFlags(t)
			tmpDir := t.TempDir()

			forgeGenSVG = true
			forgeAppName = "tmpltest"
			forgeOutputDir = tmpDir
			forgeTemplate = tmpl
			forgeSkipMac = true
			forgeSkipLinux = true

			_, err := executeForge(t)
			if err != nil {
				t.Fatalf("forge --generate --template %s failed: %v", tmpl, err)
			}

			svgPath := filepath.Join(tmpDir, "tmpltest.svg")
			if _, err := os.Stat(svgPath); os.IsNotExist(err) {
				t.Errorf("expected SVG file for template %s", tmpl)
			}
		})
	}
}

func TestForge_InvalidTemplate(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "test"
	forgeOutputDir = tmpDir
	forgeTemplate = "nonexistent"

	_, err := executeForge(t)
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
	if !containsStr(err.Error(), "nonexistent") {
		t.Errorf("error should mention the template name, got: %v", err)
	}
}

func TestForge_WatchWithGenerate(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "watchtest"
	forgeOutputDir = tmpDir
	forgeWatch = true

	_, err := executeForge(t)
	if err == nil {
		t.Fatal("expected error for --watch with --generate, got nil")
	}
	if !containsStr(err.Error(), "cannot use --watch with --generate") {
		t.Errorf("expected watch/generate conflict error, got: %v", err)
	}
}

func TestForge_WatchNoSource(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeAppName = "watchtest"
	forgeOutputDir = tmpDir
	forgeWatch = true

	_, err := executeForge(t)
	if err == nil {
		t.Fatal("expected error for --watch with no source, got nil")
	}
}

func TestForge_LinuxOutput(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "linuxtest"
	forgeOutputDir = tmpDir
	forgeSkipWin = true
	forgeSkipMac = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge for linux failed: %v", err)
	}

	// Verify .desktop file exists
	desktopPath := filepath.Join(tmpDir, "linux", "linuxtest.desktop")
	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		t.Error("expected .desktop file to exist")
	}

	// Verify hicolor structure (linux/icons/hicolor/)
	hicolorPath := filepath.Join(tmpDir, "linux", "icons", "hicolor")
	if _, err := os.Stat(hicolorPath); os.IsNotExist(err) {
		t.Error("expected linux/icons/hicolor/ directory to exist")
	}

	if !containsStr(output, "Linux") {
		t.Error("expected 'Linux' in output")
	}
}

func TestForge_MacOSOutput(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "mactest"
	forgeOutputDir = tmpDir
	forgeSkipWin = true
	forgeSkipLinux = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge for macos failed: %v", err)
	}

	// Verify .app bundle
	appDir := filepath.Join(tmpDir, "macos", "mactest.app")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		t.Error("expected .app bundle to exist")
	}

	// Verify Info.plist
	plistPath := filepath.Join(appDir, "Contents", "Info.plist")
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		t.Error("expected Info.plist to exist")
	}

	if !containsStr(output, "macOS") {
		t.Error("expected 'macOS' in output")
	}
}

func TestForge_WindowsOutput(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "wintest"
	forgeOutputDir = tmpDir
	forgeSkipMac = true
	forgeSkipLinux = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge for windows failed: %v", err)
	}

	// Verify ICO
	icoPath := filepath.Join(tmpDir, "windows", "icon.ico")
	if _, err := os.Stat(icoPath); os.IsNotExist(err) {
		t.Error("expected windows/icon.ico to exist")
	}

	// Verify versioninfo.json
	viPath := filepath.Join(tmpDir, "windows", "versioninfo.json")
	if _, err := os.Stat(viPath); os.IsNotExist(err) {
		t.Error("expected versioninfo.json to exist")
	}

	// Verify .syso was generated (winres pure Go — should always succeed)
	sysoPath := filepath.Join(tmpDir, "windows", "resource_windows_amd64.syso")
	if _, err := os.Stat(sysoPath); os.IsNotExist(err) {
		// .syso generation is non-fatal in forge; check output for note
		if !containsStr(output, "Note:") {
			t.Error("expected .syso file or generation note in output")
		}
	} else {
		if !containsStr(output, "winres") {
			t.Error("expected 'winres' in output for .syso generation")
		}
	}
}

func TestForge_FaviconAndIconsetCombined(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "combined"
	forgeOutputDir = tmpDir
	forgeFavicon = true
	forgeIconset = true

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --favicon --iconset failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "favicon")); os.IsNotExist(err) {
		t.Error("expected favicon/ directory")
	}

	entries, _ := os.ReadDir(filepath.Join(tmpDir, "macos"))
	foundIconset := false
	for _, e := range entries {
		if containsStr(e.Name(), "iconset") {
			foundIconset = true
			break
		}
	}
	if !foundIconset {
		t.Error("expected .iconset in macos/")
	}
}

func TestForge_Arch386(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "archtest"
	forgeOutputDir = tmpDir
	forgeArch = "386"
	forgeSkipMac = true
	forgeSkipLinux = true

	_, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --arch 386 failed: %v", err)
	}

	// .syso generation is non-fatal; just verify the pipeline completed
	sysoPath := filepath.Join(tmpDir, "windows", "resource_windows_386.syso")
	if _, err := os.Stat(sysoPath); os.IsNotExist(err) {
		t.Log("note: .syso with 386 arch not created (non-fatal in forge)")
	}
}

func TestForge_CustomMetadata(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "metaapp"
	forgeOutputDir = tmpDir
	forgeVersion = "2.5.0"
	forgeCompany = "MetaCorp"
	forgeCopyright = "Copyright 2026 MetaCorp"
	forgeBundleID = "com.metacorp.metaapp"
	forgeSkipLinux = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge with metadata failed: %v", err)
	}

	// Verify versioninfo.json has the metadata
	viPath := filepath.Join(tmpDir, "windows", "versioninfo.json")
	data, err := os.ReadFile(viPath)
	if err != nil {
		t.Fatalf("read versioninfo.json: %v", err)
	}
	viStr := string(data)
	if !containsStr(viStr, "MetaCorp") {
		t.Error("versioninfo.json missing company name")
	}
	if !containsStr(viStr, "2.5.0") {
		t.Error("versioninfo.json missing version")
	}

	if output == "" {
		t.Error("expected non-empty output")
	}
}

func TestForge_Preview(t *testing.T) {
	resetForgeFlags(t)
	tmpDir := t.TempDir()

	forgeGenSVG = true
	forgeAppName = "previewtest"
	forgeOutputDir = tmpDir
	forgePreview = true
	forgePreviewSize = 16
	forgeSkipWin = true
	forgeSkipMac = true
	forgeSkipLinux = true

	output, err := executeForge(t)
	if err != nil {
		t.Fatalf("forge --preview failed: %v", err)
	}

	// ANSI preview should contain escape codes
	if !containsStr(output, "\033[") {
		t.Error("expected ANSI escape codes in preview output")
	}
	if !containsStr(output, "Forge complete!") {
		t.Error("expected 'Forge complete!' after preview")
	}
}

func TestPreviewIcon(t *testing.T) {
	// Create a small test image
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Pix[(y*64+x)*4+0] = uint8(x * 4) // R
			img.Pix[(y*64+x)*4+1] = uint8(y * 4) // G
			img.Pix[(y*64+x)*4+2] = 128           // B
			img.Pix[(y*64+x)*4+3] = 255           // A
		}
	}

	buf := new(bytes.Buffer)
	previewIcon(buf, img, 16)

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty preview output")
	}
	if !containsStr(output, "\033[") {
		t.Error("expected ANSI escape codes")
	}
	if !containsStr(output, "▀") {
		t.Error("expected half-block characters in preview")
	}
}

func TestPreviewIcon_Transparent(t *testing.T) {
	// Image with transparent pixels
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))

	buf := new(bytes.Buffer)
	previewIcon(buf, img, 8)

	output := buf.String()
	if output == "" {
		t.Error("expected non-empty output even for transparent image")
	}
}

// containsStr checks if s contains substr.
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
