package svg

import (
	"image"
	"os"
	"path/filepath"
	"testing"
)

const testSVG = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">
  <rect x="10" y="10" width="80" height="80" fill="#4F46E5"/>
  <circle cx="50" cy="50" r="30" fill="#F59E0B"/>
</svg>`

func writeTempSVG(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.svg")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp svg: %v", err)
	}
	return path
}

// hasNonZeroPixel reports whether img contains at least one pixel whose
// RGBA values are not all zero.
func hasNonZeroPixel(img *image.RGBA) bool {
	for i := 0; i < len(img.Pix); i += 4 {
		if img.Pix[i] != 0 || img.Pix[i+1] != 0 || img.Pix[i+2] != 0 || img.Pix[i+3] != 0 {
			return true
		}
	}
	return false
}

func TestRenderToImage(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"16x16", 16},
		{"32x32", 32},
		{"64x64", 64},
		{"128x128", 128},
		{"256x256", 256},
	}

	svgPath := writeTempSVG(t, testSVG)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := RenderToImage(svgPath, tt.size)
			if err != nil {
				t.Fatalf("RenderToImage(%d) returned error: %v", tt.size, err)
			}

			if img == nil {
				t.Fatal("RenderToImage returned nil image")
			}

			bounds := img.Bounds()
			if bounds.Dx() != tt.size || bounds.Dy() != tt.size {
				t.Errorf("expected %dx%d, got %dx%d", tt.size, tt.size, bounds.Dx(), bounds.Dy())
			}

			if !hasNonZeroPixel(img) {
				t.Error("image has all zero pixels; expected something to be drawn")
			}
		})
	}
}

func TestRenderToImages(t *testing.T) {
	tests := []struct {
		name  string
		sizes []int
	}{
		{"single size", []int{32}},
		{"two sizes", []int{16, 64}},
		{"standard icon sizes", []int{16, 32, 64, 128, 256}},
	}

	svgPath := writeTempSVG(t, testSVG)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RenderToImages(svgPath, tt.sizes)
			if err != nil {
				t.Fatalf("RenderToImages returned error: %v", err)
			}

			if len(result) != len(tt.sizes) {
				t.Fatalf("expected %d images, got %d", len(tt.sizes), len(result))
			}

			for _, size := range tt.sizes {
				img, ok := result[size]
				if !ok {
					t.Errorf("missing image for size %d", size)
					continue
				}

				if img == nil {
					t.Errorf("nil image for size %d", size)
					continue
				}

				bounds := img.Bounds()
				if bounds.Dx() != size || bounds.Dy() != size {
					t.Errorf("size %d: expected %dx%d, got %dx%d", size, size, size, bounds.Dx(), bounds.Dy())
				}

				if !hasNonZeroPixel(img) {
					t.Errorf("size %d: image has all zero pixels", size)
				}
			}
		})
	}
}

func TestRenderToImage_NonExistentFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.svg")
	img, err := RenderToImage(path, 64)
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if img != nil {
		t.Error("expected nil image on error")
	}
	if got := err.Error(); !containsSubstring(got, "open svg") {
		t.Errorf("error %q does not contain %q", got, "open svg")
	}
}

func TestRenderToImage_InvalidSVG(t *testing.T) {
	// oksvg.ReadIconStream may not error on arbitrary content; it silently
	// produces an empty icon. Verify no panic and a valid (possibly blank)
	// image is returned.
	path := writeTempSVG(t, "this is not valid svg at all")
	img, err := RenderToImage(path, 64)
	if err != nil {
		// If the parser does reject it, that is acceptable.
		return
	}
	if img == nil {
		t.Fatal("expected non-nil image when no error is returned")
	}
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		t.Errorf("expected 64x64, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderToImage_EmptyFile(t *testing.T) {
	// Similar to invalid SVG: oksvg may not reject an empty stream.
	path := writeTempSVG(t, "")
	img, err := RenderToImage(path, 64)
	if err != nil {
		return
	}
	if img == nil {
		t.Fatal("expected non-nil image when no error is returned")
	}
	bounds := img.Bounds()
	if bounds.Dx() != 64 || bounds.Dy() != 64 {
		t.Errorf("expected 64x64, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderToImages_PropagatesError(t *testing.T) {
	nonExistent := filepath.Join(t.TempDir(), "missing.svg")

	result, err := RenderToImages(nonExistent, []int{16, 32})
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}

	if result != nil {
		t.Error("expected nil result on error")
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
