package generator

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/srwiley/oksvg"
)

func TestDefaultPalette(t *testing.T) {
	p := DefaultPalette()
	if p.Primary != "#4F46E5" {
		t.Errorf("Primary = %q, want %q", p.Primary, "#4F46E5")
	}
	if p.Secondary != "#7C3AED" {
		t.Errorf("Secondary = %q, want %q", p.Secondary, "#7C3AED")
	}
	if p.Accent != "#F59E0B" {
		t.Errorf("Accent = %q, want %q", p.Accent, "#F59E0B")
	}
}

func TestGenerateIconSVG(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		palette ColorPalette
	}{
		{
			name:    "default palette",
			appName: "TestApp",
			palette: DefaultPalette(),
		},
		{
			name:    "custom palette",
			appName: "CustomApp",
			palette: ColorPalette{
				Primary:   "#FF0000",
				Secondary: "#00FF00",
				Accent:    "#0000FF",
			},
		},
		{
			name:    "dark palette",
			appName: "DarkApp",
			palette: ColorPalette{
				Primary:   "#1A1A2E",
				Secondary: "#16213E",
				Accent:    "#E94560",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := GenerateIconSVG(outPath, tt.appName, tt.palette); err != nil {
				t.Fatalf("GenerateIconSVG() error = %v", err)
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("reading output file: %v", err)
			}

			content := string(data)

			// File should be non-empty and start with XML declaration.
			if len(data) == 0 {
				t.Fatal("generated SVG is empty")
			}
			if !strings.HasPrefix(content, "<?xml") {
				t.Error("SVG does not start with <?xml declaration")
			}
		})
	}
}

func TestGenerateIconSVG_XMLStructure(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "icon.svg")
	if err := GenerateIconSVG(outPath, "XMLTest", DefaultPalette()); err != nil {
		t.Fatalf("GenerateIconSVG() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	// Verify the file is valid XML.
	var root struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(data, &root); err != nil {
		t.Fatalf("SVG is not valid XML: %v", err)
	}

	if root.XMLName.Local != "svg" {
		t.Errorf("root element = %q, want %q", root.XMLName.Local, "svg")
	}
}

func TestGenerateIconSVG_ColorPaletteInOutput(t *testing.T) {
	tests := []struct {
		name    string
		palette ColorPalette
	}{
		{
			name:    "default colors",
			palette: DefaultPalette(),
		},
		{
			name: "custom colors",
			palette: ColorPalette{
				Primary:   "#AABBCC",
				Secondary: "#DDEEFF",
				Accent:    "#112233",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := GenerateIconSVG(outPath, "ColorTest", tt.palette); err != nil {
				t.Fatalf("GenerateIconSVG() error = %v", err)
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("reading output file: %v", err)
			}
			content := string(data)

			if !strings.Contains(content, tt.palette.Primary) {
				t.Errorf("SVG does not contain primary color %s", tt.palette.Primary)
			}
			if !strings.Contains(content, tt.palette.Secondary) {
				t.Errorf("SVG does not contain secondary color %s", tt.palette.Secondary)
			}
			if !strings.Contains(content, tt.palette.Accent) {
				t.Errorf("SVG does not contain accent color %s", tt.palette.Accent)
			}
		})
	}
}

func TestGenerateIconSVG_OksvgParseable(t *testing.T) {
	tests := []struct {
		name    string
		palette ColorPalette
	}{
		{"default palette", DefaultPalette()},
		{"custom palette", ColorPalette{Primary: "#FF5733", Secondary: "#33FF57", Accent: "#3357FF"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := GenerateIconSVG(outPath, "OksvgTest", tt.palette); err != nil {
				t.Fatalf("GenerateIconSVG() error = %v", err)
			}

			icon, err := oksvg.ReadIcon(outPath, oksvg.StrictErrorMode)
			if err != nil {
				t.Fatalf("oksvg.ReadIcon() could not parse SVG: %v", err)
			}

			vb := icon.ViewBox
			if vb.W != 512 || vb.H != 512 {
				t.Errorf("viewBox = %gx%g, want 512x512", vb.W, vb.H)
			}
		})
	}
}

func TestGenerateIconSVG_DifferentAppNames(t *testing.T) {
	// The current implementation does not embed the app name in the SVG,
	// but we verify that the function works correctly with different names
	// and produces valid output in each case.
	names := []string{"AppAlpha", "AppBeta", "My Icon App", ""}

	for _, name := range names {
		t.Run("appName="+name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := GenerateIconSVG(outPath, name, DefaultPalette()); err != nil {
				t.Fatalf("GenerateIconSVG(%q) error = %v", name, err)
			}

			info, err := os.Stat(outPath)
			if err != nil {
				t.Fatalf("stat output file: %v", err)
			}
			if info.Size() == 0 {
				t.Error("generated SVG file is empty")
			}
		})
	}
}

func BenchmarkGenerateIconSVG(b *testing.B) {
	dir := b.TempDir()
	palette := DefaultPalette()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outPath := filepath.Join(dir, fmt.Sprintf("bench_%d.svg", i))
		if err := GenerateIconSVG(outPath, "BenchApp", palette); err != nil {
			b.Fatalf("GenerateIconSVG: %v", err)
		}
	}
}

func TestGenerateIconSVG_InvalidOutputPath(t *testing.T) {
	// Use a path where the file itself cannot be created (e.g. a null byte
	// on Unix or a reserved device name trick). We use a read-only directory
	// to force a write failure on all platforms.
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0o555); err != nil {
		t.Fatalf("creating readonly dir: %v", err)
	}

	// Attempt to write into a path that treats an existing file as a directory.
	blocker := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("creating blocker file: %v", err)
	}

	// MkdirAll should fail because "blocker" is a file, not a directory.
	badPath := filepath.Join(blocker, "subdir", "icon.svg")
	err := GenerateIconSVG(badPath, "Test", DefaultPalette())
	if err == nil {
		t.Fatal("expected error for invalid output path, got nil")
	}
}

func TestGenerateIconSVG_CreatesParentDirs(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "a", "b", "c", "icon.svg")
	if err := GenerateIconSVG(outPath, "DirTest", DefaultPalette()); err != nil {
		t.Fatalf("GenerateIconSVG() error = %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("output file does not exist: %v", err)
	}
}

func TestGenerateIconSVG_FileContent(t *testing.T) {
	outPath := filepath.Join(t.TempDir(), "icon.svg")
	if err := GenerateIconSVG(outPath, "ContentTest", DefaultPalette()); err != nil {
		t.Fatalf("GenerateIconSVG() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	content := string(data)

	// Check essential SVG elements are present.
	checks := []struct {
		desc    string
		substr  string
	}{
		{"xml declaration", "<?xml"},
		{"svg open tag", "<svg"},
		{"svg namespace", `xmlns="http://www.w3.org/2000/svg"`},
		{"viewBox", `viewBox="0 0 512 512"`},
		{"defs section", "<defs>"},
		{"linearGradient", "linearGradient"},
		{"rect element", "<rect"},
		{"polygon element", "<polygon"},
		{"circle element", "<circle"},
		{"path element", "<path"},
		{"closing svg tag", "</svg>"},
	}

	for _, c := range checks {
		if !strings.Contains(content, c.substr) {
			t.Errorf("SVG missing %s (expected substring %q)", c.desc, c.substr)
		}
	}
}
