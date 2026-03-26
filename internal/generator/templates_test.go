package generator

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/srwiley/oksvg"
)

func TestListTemplates_AllRegistered(t *testing.T) {
	templates := ListTemplates()

	expected := []string{
		"bolt", "cube", "forge", "gear", "hexagon",
		"leaf", "shield", "stack", "terminal", "wave",
	}

	if len(templates) != len(expected) {
		t.Fatalf("ListTemplates() returned %d templates, want %d", len(templates), len(expected))
	}

	for i, tmpl := range templates {
		if tmpl.Name != expected[i] {
			t.Errorf("template[%d].Name = %q, want %q", i, tmpl.Name, expected[i])
		}
		if tmpl.Description == "" {
			t.Errorf("template %q has empty description", tmpl.Name)
		}
		if tmpl.GenerateFn == nil {
			t.Errorf("template %q has nil GenerateFn", tmpl.Name)
		}
	}
}

func TestGetTemplate_Found(t *testing.T) {
	tmpl, err := GetTemplate("forge")
	if err != nil {
		t.Fatalf("GetTemplate(forge) error = %v", err)
	}
	if tmpl.Name != "forge" {
		t.Errorf("Name = %q, want %q", tmpl.Name, "forge")
	}
}

func TestGetTemplate_NotFound(t *testing.T) {
	_, err := GetTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown template, got nil")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention the unknown name, got: %v", err)
	}
}

func TestAllTemplates_OksvgParseable(t *testing.T) {
	palette := ColorPalette{
		Primary:   "#1A2B3C",
		Secondary: "#4D5E6F",
		Accent:    "#FF9900",
	}

	for _, tmpl := range ListTemplates() {
		t.Run(tmpl.Name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := tmpl.GenerateFn(outPath, "TestApp", palette); err != nil {
				t.Fatalf("GenerateFn() error = %v", err)
			}

			icon, err := oksvg.ReadIcon(outPath, oksvg.StrictErrorMode)
			if err != nil {
				t.Fatalf("oksvg.ReadIcon() failed: %v", err)
			}

			vb := icon.ViewBox
			if vb.W != 512 || vb.H != 512 {
				t.Errorf("viewBox = %gx%g, want 512x512", vb.W, vb.H)
			}
		})
	}
}

func TestAllTemplates_ValidXML(t *testing.T) {
	palette := DefaultPalette()

	for _, tmpl := range ListTemplates() {
		t.Run(tmpl.Name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := tmpl.GenerateFn(outPath, "XMLTest", palette); err != nil {
				t.Fatalf("GenerateFn() error = %v", err)
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("reading file: %v", err)
			}

			var root struct{ XMLName xml.Name }
			if err := xml.Unmarshal(data, &root); err != nil {
				t.Fatalf("invalid XML: %v", err)
			}

			if root.XMLName.Local != "svg" {
				t.Errorf("root element = %q, want %q", root.XMLName.Local, "svg")
			}
		})
	}
}

func TestAllTemplates_ContainPaletteColors(t *testing.T) {
	palette := ColorPalette{
		Primary:   "#AABB11",
		Secondary: "#22CC33",
		Accent:    "#DD44EE",
	}

	for _, tmpl := range ListTemplates() {
		t.Run(tmpl.Name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "icon.svg")
			if err := tmpl.GenerateFn(outPath, "ColorTest", palette); err != nil {
				t.Fatalf("GenerateFn() error = %v", err)
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("reading file: %v", err)
			}
			content := string(data)

			if !strings.Contains(content, palette.Primary) {
				t.Errorf("SVG missing primary color %s", palette.Primary)
			}
			if !strings.Contains(content, palette.Secondary) {
				t.Errorf("SVG missing secondary color %s", palette.Secondary)
			}
			// Accent may not be used in all templates (e.g., wave uses only accent),
			// but all current templates do use it.
			if !strings.Contains(content, palette.Accent) {
				t.Errorf("SVG missing accent color %s", palette.Accent)
			}
		})
	}
}

func TestAllTemplates_CreateParentDirs(t *testing.T) {
	palette := DefaultPalette()

	for _, tmpl := range ListTemplates() {
		t.Run(tmpl.Name, func(t *testing.T) {
			outPath := filepath.Join(t.TempDir(), "a", "b", "icon.svg")
			if err := tmpl.GenerateFn(outPath, "DirTest", palette); err != nil {
				t.Fatalf("GenerateFn() error = %v", err)
			}
			if _, err := os.Stat(outPath); err != nil {
				t.Fatalf("output file not created: %v", err)
			}
		})
	}
}

func BenchmarkAllTemplates(b *testing.B) {
	palette := DefaultPalette()
	templates := ListTemplates()

	for _, tmpl := range templates {
		b.Run(tmpl.Name, func(b *testing.B) {
			dir := b.TempDir()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				outPath := filepath.Join(dir, "icon.svg")
				if err := tmpl.GenerateFn(outPath, "Bench", palette); err != nil {
					b.Fatalf("GenerateFn: %v", err)
				}
			}
		})
	}
}
