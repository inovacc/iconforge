package icon

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func newTestImage(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	return img
}

func TestWriteICO(t *testing.T) {
	tests := []struct {
		name      string
		sizes     []int
		wantErr   bool
		wantCount int
	}{
		{
			name:      "single 16x16 image",
			sizes:     []int{16},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:      "standard icon sizes",
			sizes:     []int{16, 32, 48, 64, 128, 256},
			wantErr:   false,
			wantCount: 6,
		},
		{
			name:      "single 256x256 image",
			sizes:     []int{256},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:      "two sizes",
			sizes:     []int{32, 64},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:      "large 512x512 image",
			sizes:     []int{512},
			wantErr:   false,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outPath := filepath.Join(tmpDir, "test.ico")

			images := make(map[int]*image.RGBA)
			for _, s := range tt.sizes {
				images[s] = newTestImage(s)
			}

			err := WriteICO(outPath, images)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteICO() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("failed to read ICO file: %v", err)
			}

			if len(data) < 6 {
				t.Fatal("ICO file too small for header")
			}

			reserved := binary.LittleEndian.Uint16(data[0:2])
			icoType := binary.LittleEndian.Uint16(data[2:4])
			count := binary.LittleEndian.Uint16(data[4:6])

			if reserved != 0 {
				t.Errorf("ICO reserved field = %d, want 0", reserved)
			}
			if icoType != 1 {
				t.Errorf("ICO type = %d, want 1", icoType)
			}
			if int(count) != tt.wantCount {
				t.Errorf("ICO entry count = %d, want %d", count, tt.wantCount)
			}

			expectedMinSize := 6 + 16*tt.wantCount
			if len(data) < expectedMinSize {
				t.Errorf("ICO file size %d < minimum expected %d", len(data), expectedMinSize)
			}
		})
	}
}

func TestWriteICO_VerifyEntries(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "verify.ico")

	sizes := []int{16, 32, 256}
	images := make(map[int]*image.RGBA)
	for _, s := range sizes {
		images[s] = newTestImage(s)
	}

	if err := WriteICO(outPath, images); err != nil {
		t.Fatalf("WriteICO() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read ICO file: %v", err)
	}

	count := int(binary.LittleEndian.Uint16(data[4:6]))
	if count != 3 {
		t.Fatalf("expected 3 entries, got %d", count)
	}

	for i := 0; i < count; i++ {
		offset := 6 + i*16
		entry := data[offset : offset+16]

		w := entry[0]
		h := entry[1]
		planes := binary.LittleEndian.Uint16(entry[4:6])
		bitCount := binary.LittleEndian.Uint16(entry[6:8])
		bytesInRes := binary.LittleEndian.Uint32(entry[8:12])
		imageOffset := binary.LittleEndian.Uint32(entry[12:16])

		if w != h {
			t.Errorf("entry %d: width %d != height %d", i, w, h)
		}
		if planes != 1 {
			t.Errorf("entry %d: planes = %d, want 1", i, planes)
		}
		if bitCount != 32 {
			t.Errorf("entry %d: bitCount = %d, want 32", i, bitCount)
		}
		if bytesInRes == 0 {
			t.Errorf("entry %d: bytesInRes is 0", i)
		}
		if int(imageOffset+bytesInRes) > len(data) {
			t.Errorf("entry %d: data extends beyond file (offset=%d, size=%d, fileLen=%d)",
				i, imageOffset, bytesInRes, len(data))
		}
	}
}

func TestWriteICO_ReadBack(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "readback.ico")

	sizes := []int{16, 32, 48}
	images := make(map[int]*image.RGBA)
	for _, s := range sizes {
		images[s] = newTestImage(s)
	}

	if err := WriteICO(outPath, images); err != nil {
		t.Fatalf("WriteICO() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read ICO file: %v", err)
	}

	count := int(binary.LittleEndian.Uint16(data[4:6]))

	for i := 0; i < count; i++ {
		entryOffset := 6 + i*16
		entry := data[entryOffset : entryOffset+16]
		bytesInRes := binary.LittleEndian.Uint32(entry[8:12])
		imageOffset := binary.LittleEndian.Uint32(entry[12:16])

		pngData := data[imageOffset : imageOffset+bytesInRes]
		img, err := png.Decode(bytes.NewReader(pngData))
		if err != nil {
			t.Errorf("entry %d: failed to decode embedded PNG: %v", i, err)
			continue
		}
		bounds := img.Bounds()
		if bounds.Dx() != bounds.Dy() {
			t.Errorf("entry %d: image is not square (%dx%d)", i, bounds.Dx(), bounds.Dy())
		}
	}
}

func TestWriteICO_InvalidPath(t *testing.T) {
	images := map[int]*image.RGBA{16: newTestImage(16)}
	err := WriteICO("/nonexistent/dir/test.ico", images)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestWriteICO_EmptyMap(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "empty.ico")

	images := make(map[int]*image.RGBA)
	err := WriteICO(outPath, images)
	if err != nil {
		t.Fatalf("WriteICO() with empty map error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read ICO file: %v", err)
	}

	count := binary.LittleEndian.Uint16(data[4:6])
	if count != 0 {
		t.Errorf("expected 0 entries for empty map, got %d", count)
	}
}

func TestWriteICO_256SizeUsesZeroByte(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "size256.ico")

	images := map[int]*image.RGBA{256: newTestImage(256)}

	if err := WriteICO(outPath, images); err != nil {
		t.Fatalf("WriteICO() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read ICO file: %v", err)
	}

	w := data[6]
	h := data[7]
	if w != 0 || h != 0 {
		t.Errorf("256px entry should have width=0, height=0 in ICO format, got w=%d h=%d", w, h)
	}
}
