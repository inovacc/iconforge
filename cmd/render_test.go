package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// resetRenderFlags resets all render global flag variables to their defaults.
func resetRenderFlags(t *testing.T) {
	t.Helper()

	renderSVGPath = ""
	renderPNGPath = ""
	renderOutputDir = "build/icons/png"
	renderSizesStr = "512,256,128,64,48,32,16"
}

// executeRender sets up the renderCmd output buffers and calls its RunE directly.
func executeRender(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	renderCmd.SetOut(buf)
	renderCmd.SetErr(buf)
	err := renderCmd.RunE(renderCmd, nil)
	return buf.String(), err
}

func TestRender_FromSVG(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()
	svgDir := t.TempDir()
	svgPath := writeTestSVG(t, svgDir)

	renderSVGPath = svgPath
	renderOutputDir = tmpDir

	_, err := executeRender(t)
	if err != nil {
		t.Fatalf("render --svg failed: %v", err)
	}

	// Verify default sizes were rendered
	defaultSizes := []int{512, 256, 128, 64, 48, 32, 16}
	for _, size := range defaultSizes {
		filename := fmt.Sprintf("%dx%d.png", size, size)
		p := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", filename)
		}
	}
}

func TestRender_FromPNG(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()
	pngDir := t.TempDir()
	pngPath := writeTestPNG(t, pngDir)

	renderPNGPath = pngPath
	renderOutputDir = tmpDir

	_, err := executeRender(t)
	if err != nil {
		t.Fatalf("render --png failed: %v", err)
	}

	// Verify default sizes were rendered
	defaultSizes := []int{512, 256, 128, 64, 48, 32, 16}
	for _, size := range defaultSizes {
		filename := fmt.Sprintf("%dx%d.png", size, size)
		p := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", filename)
		}
	}
}

func TestRender_CustomSizes(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()
	svgDir := t.TempDir()
	svgPath := writeTestSVG(t, svgDir)

	renderSVGPath = svgPath
	renderOutputDir = tmpDir
	renderSizesStr = "64,32"

	_, err := executeRender(t)
	if err != nil {
		t.Fatalf("render --sizes 64,32 failed: %v", err)
	}

	// Verify only the requested sizes exist
	wantSizes := []int{64, 32}
	for _, size := range wantSizes {
		filename := fmt.Sprintf("%dx%d.png", size, size)
		p := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("expected %s to exist", filename)
		}
	}

	// Verify sizes NOT requested do NOT exist
	unwantedSizes := []int{512, 256, 128, 48, 16}
	for _, size := range unwantedSizes {
		filename := fmt.Sprintf("%dx%d.png", size, size)
		p := filepath.Join(tmpDir, filename)
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("expected %s to NOT exist when not in --sizes", filename)
		}
	}
}

func TestRender_MutuallyExclusive(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()

	renderSVGPath = "some.svg"
	renderPNGPath = "some.png"
	renderOutputDir = tmpDir

	_, err := executeRender(t)
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags, got nil")
	}
	expected := "--svg and --png are mutually exclusive"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestRender_NoSource(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()
	renderOutputDir = tmpDir

	_, err := executeRender(t)
	if err == nil {
		t.Fatal("expected error when no source specified, got nil")
	}
	expected := "either --svg or --png must be provided"
	if err.Error() != expected {
		t.Errorf("error = %q, want %q", err.Error(), expected)
	}
}

func TestParseSizes_Valid(t *testing.T) {
	tests := []struct {
		input string
		want  []int
	}{
		{"512,256,128", []int{512, 256, 128}},
		{"64", []int{64}},
		{"16, 32, 48", []int{16, 32, 48}},
		{"1024", []int{1024}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseSizes(tt.input)
			if err != nil {
				t.Fatalf("parseSizes(%q) returned error: %v", tt.input, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseSizes(%q) returned %d sizes, want %d", tt.input, len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseSizes(%q)[%d] = %d, want %d", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseSizes_Invalid(t *testing.T) {
	tests := []struct {
		input   string
		wantErr string
	}{
		{"abc", "invalid size"},
		{"0", "out of range"},
		{"-1", "out of range"},
		{"3000", "out of range"},
		{"512,abc,128", "invalid size"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := parseSizes(tt.input)
			if err == nil {
				t.Fatalf("parseSizes(%q) expected error, got nil", tt.input)
			}
			if !containsStr(err.Error(), tt.wantErr) {
				t.Errorf("parseSizes(%q) error = %q, want containing %q", tt.input, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestAppName_Explicit(t *testing.T) {
	got := AppName("myapp")
	if got != "myapp" {
		t.Errorf("AppName('myapp') = %q, want 'myapp'", got)
	}
}

func TestAppName_Empty(t *testing.T) {
	got := AppName("")
	if got == "" {
		t.Error("AppName('') should not return empty string")
	}
}

func TestRender_InvalidSizes(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()
	svgDir := t.TempDir()
	svgPath := writeTestSVG(t, svgDir)

	renderSVGPath = svgPath
	renderOutputDir = tmpDir
	renderSizesStr = "abc"

	_, err := executeRender(t)
	if err == nil {
		t.Fatal("expected error for invalid sizes, got nil")
	}
	if !containsStr(err.Error(), "invalid size") {
		t.Errorf("expected 'invalid size' error, got: %v", err)
	}
}

func TestRender_InvalidSVGPath(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()

	renderSVGPath = filepath.Join(tmpDir, "nonexistent.svg")
	renderOutputDir = tmpDir

	_, err := executeRender(t)
	if err == nil {
		t.Fatal("expected error for nonexistent SVG, got nil")
	}
}

func TestRender_InvalidPNGPath(t *testing.T) {
	resetRenderFlags(t)
	tmpDir := t.TempDir()

	renderPNGPath = filepath.Join(tmpDir, "nonexistent.png")
	renderOutputDir = tmpDir

	_, err := executeRender(t)
	if err == nil {
		t.Fatal("expected error for nonexistent PNG, got nil")
	}
}
