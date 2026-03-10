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

// containsStr checks if s contains substr.
func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
