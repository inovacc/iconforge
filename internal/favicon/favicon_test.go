package favicon

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// newTestImage creates a solid-color RGBA image of the given size.
func newTestImage(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			img.SetRGBA(x, y, color.RGBA{R: 0x42, G: 0x84, B: 0xC6, A: 0xFF})
		}
	}
	return img
}

func TestGenerateFavicons_AllFilesCreated(t *testing.T) {
	dir := t.TempDir()

	images := map[int]*image.RGBA{
		16:  newTestImage(16),
		32:  newTestImage(32),
		180: newTestImage(180),
		192: newTestImage(192),
		512: newTestImage(512),
	}

	if err := GenerateFavicons(images, dir); err != nil {
		t.Fatalf("GenerateFavicons: %v", err)
	}

	expected := []string{
		"favicon.ico",
		"favicon-16x16.png",
		"favicon-32x32.png",
		"apple-touch-icon.png",
		"android-chrome-192x192.png",
		"android-chrome-512x512.png",
		"site.webmanifest",
	}

	for _, name := range expected {
		path := filepath.Join(dir, name)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected file %s not found: %v", name, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("file %s is empty", name)
		}
	}
}

func TestGenerateFavicons_ICOMagicBytes(t *testing.T) {
	dir := t.TempDir()

	images := map[int]*image.RGBA{
		512: newTestImage(512),
	}

	if err := GenerateFavicons(images, dir); err != nil {
		t.Fatalf("GenerateFavicons: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "favicon.ico"))
	if err != nil {
		t.Fatalf("read favicon.ico: %v", err)
	}

	// ICO magic: 00 00 01 00
	if len(data) < 4 {
		t.Fatal("favicon.ico too small")
	}
	if data[0] != 0x00 || data[1] != 0x00 || data[2] != 0x01 || data[3] != 0x00 {
		t.Errorf("favicon.ico magic bytes = %x %x %x %x, want 00 00 01 00",
			data[0], data[1], data[2], data[3])
	}
}

func TestGenerateFavicons_PNGSizes(t *testing.T) {
	dir := t.TempDir()

	images := map[int]*image.RGBA{
		512: newTestImage(512),
	}

	if err := GenerateFavicons(images, dir); err != nil {
		t.Fatalf("GenerateFavicons: %v", err)
	}

	checks := map[string]int{
		"favicon-16x16.png":          16,
		"favicon-32x32.png":          32,
		"apple-touch-icon.png":       180,
		"android-chrome-192x192.png": 192,
		"android-chrome-512x512.png": 512,
	}

	for filename, wantSize := range checks {
		f, err := os.Open(filepath.Join(dir, filename))
		if err != nil {
			t.Errorf("open %s: %v", filename, err)
			continue
		}

		img, err := png.Decode(f)
		_ = f.Close()
		if err != nil {
			t.Errorf("decode %s: %v", filename, err)
			continue
		}

		bounds := img.Bounds()
		if bounds.Dx() != wantSize || bounds.Dy() != wantSize {
			t.Errorf("%s size = %dx%d, want %dx%d",
				filename, bounds.Dx(), bounds.Dy(), wantSize, wantSize)
		}
	}
}

func TestGenerateFavicons_WebmanifestValid(t *testing.T) {
	dir := t.TempDir()

	images := map[int]*image.RGBA{
		512: newTestImage(512),
	}

	if err := GenerateFavicons(images, dir); err != nil {
		t.Fatalf("GenerateFavicons: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "site.webmanifest"))
	if err != nil {
		t.Fatalf("read site.webmanifest: %v", err)
	}

	var manifest webManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("unmarshal site.webmanifest: %v", err)
	}

	if len(manifest.Icons) != 2 {
		t.Fatalf("manifest icons count = %d, want 2", len(manifest.Icons))
	}

	if manifest.Icons[0].Sizes != "192x192" {
		t.Errorf("first icon size = %s, want 192x192", manifest.Icons[0].Sizes)
	}
	if manifest.Icons[1].Sizes != "512x512" {
		t.Errorf("second icon size = %s, want 512x512", manifest.Icons[1].Sizes)
	}
}

func TestGenerateFavicons_ResizesFromLargest(t *testing.T) {
	dir := t.TempDir()

	// Only provide 512, all others should be resized
	images := map[int]*image.RGBA{
		512: newTestImage(512),
	}

	if err := GenerateFavicons(images, dir); err != nil {
		t.Fatalf("GenerateFavicons: %v", err)
	}

	// Verify all PNGs exist and have correct sizes
	checks := map[string]int{
		"favicon-16x16.png":          16,
		"favicon-32x32.png":          32,
		"apple-touch-icon.png":       180,
		"android-chrome-192x192.png": 192,
		"android-chrome-512x512.png": 512,
	}

	for filename, wantSize := range checks {
		f, err := os.Open(filepath.Join(dir, filename))
		if err != nil {
			t.Errorf("open %s: %v", filename, err)
			continue
		}

		img, err := png.Decode(f)
		_ = f.Close()
		if err != nil {
			t.Errorf("decode %s: %v", filename, err)
			continue
		}

		bounds := img.Bounds()
		if bounds.Dx() != wantSize || bounds.Dy() != wantSize {
			t.Errorf("%s size = %dx%d, want %dx%d",
				filename, bounds.Dx(), bounds.Dy(), wantSize, wantSize)
		}
	}
}

func TestSizes(t *testing.T) {
	sizes := Sizes()
	if len(sizes) != len(faviconSizes) {
		t.Fatalf("Sizes() len = %d, want %d", len(sizes), len(faviconSizes))
	}

	// Verify it returns a copy
	sizes[0] = 9999
	if faviconSizes[0] == 9999 {
		t.Error("Sizes() did not return a copy")
	}
}
