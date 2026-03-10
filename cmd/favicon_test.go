package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// resetFaviconFlags resets all favicon global flag variables to their defaults.
func resetFaviconFlags(t *testing.T) {
	t.Helper()
	faviconSVGPath = ""
	faviconPNGPath = ""
	faviconOutputDir = "build/favicons"
}

// executeFavicon sets up the faviconCmd output buffers and calls its RunE directly.
func executeFavicon(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	faviconCmd.SetOut(buf)
	faviconCmd.SetErr(buf)
	err := faviconCmd.RunE(faviconCmd, nil)
	return buf.String(), err
}

func TestFavicon_FromSVG(t *testing.T) {
	resetFaviconFlags(t)
	tmpDir := t.TempDir()
	svgDir := t.TempDir()
	svgPath := writeTestSVG(t, svgDir)

	faviconSVGPath = svgPath
	faviconOutputDir = tmpDir

	output, err := executeFavicon(t)
	if err != nil {
		t.Fatalf("favicon --svg failed: %v", err)
	}

	if !containsStr(output, "Favicons generated in") {
		t.Errorf("expected success message in output, got: %s", output)
	}

	// Verify favicon.ico exists
	if _, err := os.Stat(filepath.Join(tmpDir, "favicon.ico")); os.IsNotExist(err) {
		t.Error("expected favicon.ico to exist")
	}

	// Verify some standard PNG sizes
	expectedFiles := []string{
		"favicon-16x16.png",
		"favicon-32x32.png",
		"apple-touch-icon.png",
		"android-chrome-192x192.png",
		"android-chrome-512x512.png",
	}
	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(tmpDir, f)); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", f)
		}
	}

	// Verify webmanifest
	if _, err := os.Stat(filepath.Join(tmpDir, "site.webmanifest")); os.IsNotExist(err) {
		t.Error("expected site.webmanifest to exist")
	}
}

func TestFavicon_FromPNG(t *testing.T) {
	resetFaviconFlags(t)
	tmpDir := t.TempDir()
	pngDir := t.TempDir()
	pngPath := writeTestPNG(t, pngDir)

	faviconPNGPath = pngPath
	faviconOutputDir = tmpDir

	output, err := executeFavicon(t)
	if err != nil {
		t.Fatalf("favicon --png failed: %v", err)
	}

	if !containsStr(output, "Favicons generated in") {
		t.Errorf("expected success message in output, got: %s", output)
	}

	// Verify favicon.ico exists
	if _, err := os.Stat(filepath.Join(tmpDir, "favicon.ico")); os.IsNotExist(err) {
		t.Error("expected favicon.ico to exist")
	}
}

func TestFavicon_NoSource(t *testing.T) {
	resetFaviconFlags(t)
	tmpDir := t.TempDir()
	faviconOutputDir = tmpDir

	_, err := executeFavicon(t)
	if err == nil {
		t.Fatal("expected error when no source specified, got nil")
	}
	expected := "either --svg or --png must be provided"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestFavicon_MutuallyExclusive(t *testing.T) {
	resetFaviconFlags(t)
	tmpDir := t.TempDir()

	faviconSVGPath = "some.svg"
	faviconPNGPath = "some.png"
	faviconOutputDir = tmpDir

	_, err := executeFavicon(t)
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags, got nil")
	}
	expected := "--svg and --png are mutually exclusive"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestFavicon_InvalidSVG(t *testing.T) {
	resetFaviconFlags(t)
	tmpDir := t.TempDir()

	faviconSVGPath = filepath.Join(tmpDir, "nonexistent.svg")
	faviconOutputDir = tmpDir

	_, err := executeFavicon(t)
	if err == nil {
		t.Fatal("expected error for nonexistent SVG, got nil")
	}
}

func TestFavicon_InvalidPNG(t *testing.T) {
	resetFaviconFlags(t)
	tmpDir := t.TempDir()

	faviconPNGPath = filepath.Join(tmpDir, "nonexistent.png")
	faviconOutputDir = tmpDir

	_, err := executeFavicon(t)
	if err == nil {
		t.Fatal("expected error for nonexistent PNG, got nil")
	}
}
