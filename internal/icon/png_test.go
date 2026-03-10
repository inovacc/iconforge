package icon

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestWritePNG(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			name:    "16x16 image",
			size:    16,
			wantErr: false,
		},
		{
			name:    "32x32 image",
			size:    32,
			wantErr: false,
		},
		{
			name:    "48x48 image",
			size:    48,
			wantErr: false,
		},
		{
			name:    "64x64 image",
			size:    64,
			wantErr: false,
		},
		{
			name:    "128x128 image",
			size:    128,
			wantErr: false,
		},
		{
			name:    "256x256 image",
			size:    256,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outPath := filepath.Join(tmpDir, "test.png")

			img := newTestImage(tt.size)
			err := WritePNG(outPath, img)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WritePNG() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			f, err := os.Open(outPath)
			if err != nil {
				t.Fatalf("failed to open PNG file: %v", err)
			}
			defer func() { _ = f.Close() }()

			decoded, err := png.Decode(f)
			if err != nil {
				t.Fatalf("failed to decode PNG: %v", err)
			}

			bounds := decoded.Bounds()
			if bounds.Dx() != tt.size || bounds.Dy() != tt.size {
				t.Errorf("decoded size = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), tt.size, tt.size)
			}
		})
	}
}

func TestWritePNG_CreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "nested", "deep", "dir", "test.png")

	img := newTestImage(16)
	if err := WritePNG(outPath, img); err != nil {
		t.Fatalf("WritePNG() error = %v", err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("expected file to exist after WritePNG")
	}
}

func TestWritePNG_InvalidPath(t *testing.T) {
	img := newTestImage(16)
	err := WritePNG("", img)
	if err == nil {
		t.Error("expected error for empty path, got nil")
	}
}

func TestWritePNGs(t *testing.T) {
	tests := []struct {
		name    string
		sizes   []int
		wantErr bool
	}{
		{
			name:    "single image",
			sizes:   []int{16},
			wantErr: false,
		},
		{
			name:    "multiple images",
			sizes:   []int{16, 32, 48, 64, 128, 256},
			wantErr: false,
		},
		{
			name:    "two images",
			sizes:   []int{32, 128},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			images := make(map[int]*image.RGBA)
			for _, s := range tt.sizes {
				images[s] = newTestImage(s)
			}

			err := WritePNGs(tmpDir, images)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WritePNGs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			for _, s := range tt.sizes {
				filename := fmt.Sprintf("%dx%d.png", s, s)
				fullPath := filepath.Join(tmpDir, filename)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("expected file %s to exist", filename)
					continue
				}

				f, err := os.Open(fullPath)
				if err != nil {
					t.Errorf("failed to open %s: %v", filename, err)
					continue
				}

				decoded, err := png.Decode(f)
				_ = f.Close()
				if err != nil {
					t.Errorf("failed to decode %s: %v", filename, err)
					continue
				}

				bounds := decoded.Bounds()
				if bounds.Dx() != s || bounds.Dy() != s {
					t.Errorf("%s: decoded size = %dx%d, want %dx%d",
						filename, bounds.Dx(), bounds.Dy(), s, s)
				}
			}
		})
	}
}

func TestWritePNGs_EmptyMap(t *testing.T) {
	tmpDir := t.TempDir()

	images := make(map[int]*image.RGBA)
	err := WritePNGs(tmpDir, images)
	if err != nil {
		t.Fatalf("WritePNGs() with empty map error = %v", err)
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected 0 files for empty map, got %d", len(entries))
	}
}

func TestWritePNGs_CorrectFilenames(t *testing.T) {
	tmpDir := t.TempDir()

	sizes := []int{16, 32, 64}
	images := make(map[int]*image.RGBA)
	for _, s := range sizes {
		images[s] = newTestImage(s)
	}

	if err := WritePNGs(tmpDir, images); err != nil {
		t.Fatalf("WritePNGs() error = %v", err)
	}

	expectedFiles := map[string]bool{
		"16x16.png": false,
		"32x32.png": false,
		"64x64.png": false,
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}

	for _, entry := range entries {
		if _, ok := expectedFiles[entry.Name()]; ok {
			expectedFiles[entry.Name()] = true
		} else {
			t.Errorf("unexpected file: %s", entry.Name())
		}
	}

	for name, found := range expectedFiles {
		if !found {
			t.Errorf("expected file %s not found", name)
		}
	}
}

func TestLoadPNG(t *testing.T) {
	t.Run("valid PNG file", func(t *testing.T) {
		tmpDir := t.TempDir()
		pngPath := filepath.Join(tmpDir, "test.png")

		// Write a test PNG first
		orig := newTestImage(32)
		if err := WritePNG(pngPath, orig); err != nil {
			t.Fatalf("WritePNG() error = %v", err)
		}

		loaded, err := LoadPNG(pngPath)
		if err != nil {
			t.Fatalf("LoadPNG() error = %v", err)
		}

		bounds := loaded.Bounds()
		if bounds.Dx() != 32 || bounds.Dy() != 32 {
			t.Errorf("loaded image size = %dx%d, want 32x32", bounds.Dx(), bounds.Dy())
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := LoadPNG(filepath.Join(t.TempDir(), "nonexistent.png"))
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
	})

	t.Run("invalid PNG data", func(t *testing.T) {
		tmpDir := t.TempDir()
		badPath := filepath.Join(tmpDir, "bad.png")
		if err := os.WriteFile(badPath, []byte("not a png"), 0o644); err != nil {
			t.Fatalf("failed to write bad file: %v", err)
		}

		_, err := LoadPNG(badPath)
		if err == nil {
			t.Error("expected error for invalid PNG data, got nil")
		}
	})

	t.Run("NRGBA format conversion", func(t *testing.T) {
		// Create an NRGBA image (not RGBA), encode as PNG, then load
		tmpDir := t.TempDir()
		pngPath := filepath.Join(tmpDir, "nrgba.png")

		nrgba := image.NewNRGBA(image.Rect(0, 0, 16, 16))
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				nrgba.SetNRGBA(x, y, color.NRGBA{R: 100, G: 50, B: 200, A: 255})
			}
		}

		f, err := os.Create(pngPath)
		if err != nil {
			t.Fatalf("create file: %v", err)
		}
		if err := png.Encode(f, nrgba); err != nil {
			_ = f.Close()
			t.Fatalf("encode NRGBA: %v", err)
		}
		_ = f.Close()

		loaded, err := LoadPNG(pngPath)
		if err != nil {
			t.Fatalf("LoadPNG() error = %v", err)
		}

		if loaded == nil {
			t.Fatal("LoadPNG() returned nil image")
		}

		bounds := loaded.Bounds()
		if bounds.Dx() != 16 || bounds.Dy() != 16 {
			t.Errorf("loaded image size = %dx%d, want 16x16", bounds.Dx(), bounds.Dy())
		}
	})
}

func TestResizeImage(t *testing.T) {
	tests := []struct {
		name       string
		srcSize    int
		targetSize int
	}{
		{"downscale 128 to 32", 128, 32},
		{"upscale 16 to 64", 16, 64},
		{"same size 32 to 32", 32, 32},
		{"downscale 512 to 16", 512, 16},
		{"upscale 32 to 256", 32, 256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := newTestImage(tt.srcSize)
			result := ResizeImage(src, tt.targetSize)

			if result == nil {
				t.Fatal("ResizeImage() returned nil")
			}

			bounds := result.Bounds()
			if bounds.Dx() != tt.targetSize || bounds.Dy() != tt.targetSize {
				t.Errorf("resized image = %dx%d, want %dx%d",
					bounds.Dx(), bounds.Dy(), tt.targetSize, tt.targetSize)
			}
		})
	}
}

func BenchmarkResizeImage(b *testing.B) {
	src := newTestImage(512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ResizeImage(src, 64)
	}
}

func BenchmarkWritePNGs(b *testing.B) {
	sizes := []int{512, 256, 128, 64, 48, 32, 16}
	images := make(map[int]*image.RGBA, len(sizes))
	for _, s := range sizes {
		images[s] = newTestImage(s)
	}

	dir := b.TempDir()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outDir := filepath.Join(dir, fmt.Sprintf("bench_%d", i))
		if err := WritePNGs(outDir, images); err != nil {
			b.Fatalf("WritePNGs: %v", err)
		}
	}
}

func TestWritePNG_PixelDataPreserved(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "pixel.png")

	img := newTestImage(16)

	if err := WritePNG(outPath, img); err != nil {
		t.Fatalf("WritePNG() error = %v", err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer func() { _ = f.Close() }()

	decoded, err := png.Decode(f)
	if err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			or, og, ob, oa := img.At(x, y).RGBA()
			dr, dg, db, da := decoded.At(x, y).RGBA()
			if or != dr || og != dg || ob != db || oa != da {
				t.Errorf("pixel (%d,%d) mismatch: original=(%d,%d,%d,%d) decoded=(%d,%d,%d,%d)",
					x, y, or, og, ob, oa, dr, dg, db, da)
			}
		}
	}
}
