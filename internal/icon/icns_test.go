package icon

import (
	"encoding/binary"
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteICNS(t *testing.T) {
	tests := []struct {
		name    string
		sizes   []int
		wantErr bool
	}{
		{
			name:    "single 16x16 image",
			sizes:   []int{16},
			wantErr: false,
		},
		{
			name:    "standard icon sizes",
			sizes:   []int{16, 32, 64, 128, 256, 512},
			wantErr: false,
		},
		{
			name:    "single 256x256 image",
			sizes:   []int{256},
			wantErr: false,
		},
		{
			name:    "large 1024x1024 image",
			sizes:   []int{1024},
			wantErr: false,
		},
		{
			name:    "two sizes",
			sizes:   []int{32, 128},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outPath := filepath.Join(tmpDir, "test.icns")

			images := make(map[int]*image.RGBA)
			for _, s := range tt.sizes {
				images[s] = newTestImage(s)
			}

			err := WriteICNS(outPath, images)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteICNS() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("failed to read ICNS file: %v", err)
			}

			if len(data) < 8 {
				t.Fatal("ICNS file too small for header")
			}

			magic := string(data[0:4])
			if magic != "icns" {
				t.Errorf("ICNS magic = %q, want %q", magic, "icns")
			}

			totalSize := binary.BigEndian.Uint32(data[4:8])
			if int(totalSize) != len(data) {
				t.Errorf("ICNS total size field = %d, actual file size = %d", totalSize, len(data))
			}
		})
	}
}

func TestWriteICNS_VerifyEntries(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "verify.icns")

	sizes := []int{16, 32, 128, 256}
	images := make(map[int]*image.RGBA)
	for _, s := range sizes {
		images[s] = newTestImage(s)
	}

	if err := WriteICNS(outPath, images); err != nil {
		t.Fatalf("WriteICNS() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read ICNS file: %v", err)
	}

	validOSTypes := map[string]bool{
		"icp4": true, "icp5": true, "icp6": true,
		"ic07": true, "ic08": true, "ic09": true, "ic10": true,
	}

	offset := 8
	entryCount := 0
	for offset < len(data) {
		if offset+8 > len(data) {
			t.Fatalf("truncated entry header at offset %d", offset)
		}

		osType := string(data[offset : offset+4])
		entrySize := binary.BigEndian.Uint32(data[offset+4 : offset+8])

		if !validOSTypes[osType] {
			t.Errorf("unknown OSType %q at offset %d", osType, offset)
		}

		if entrySize < 8 {
			t.Errorf("entry %q has invalid size %d (minimum is 8)", osType, entrySize)
		}

		if offset+int(entrySize) > len(data) {
			t.Errorf("entry %q extends beyond file (offset=%d, size=%d, fileLen=%d)",
				osType, offset, entrySize, len(data))
		}

		offset += int(entrySize)
		entryCount++
	}

	if entryCount != len(sizes) {
		t.Errorf("found %d entries, want %d", entryCount, len(sizes))
	}
}

func TestWriteICNS_EmptyMap(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "empty.icns")

	images := make(map[int]*image.RGBA)
	err := WriteICNS(outPath, images)
	if err == nil {
		t.Fatal("expected error for empty image map, got nil")
	}
}

func TestWriteICNS_UnsupportedSizesOnly(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "unsupported.icns")

	images := map[int]*image.RGBA{
		48:  newTestImage(48),
		96:  newTestImage(96),
		200: newTestImage(200),
	}

	err := WriteICNS(outPath, images)
	if err == nil {
		t.Fatal("expected error for unsupported sizes only, got nil")
	}
}

func TestWriteICNS_MixedSupportedAndUnsupported(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "mixed.icns")

	images := map[int]*image.RGBA{
		32: newTestImage(32),
		48: newTestImage(48),
		64: newTestImage(64),
	}

	if err := WriteICNS(outPath, images); err != nil {
		t.Fatalf("WriteICNS() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read ICNS file: %v", err)
	}

	totalSize := binary.BigEndian.Uint32(data[4:8])
	if totalSize < 8 {
		t.Error("ICNS file size too small")
	}

	offset := 8
	entryCount := 0
	for offset < len(data) {
		entrySize := binary.BigEndian.Uint32(data[offset+4 : offset+8])
		offset += int(entrySize)
		entryCount++
	}

	if entryCount != 2 {
		t.Errorf("expected 2 entries (32 and 64 are supported, 48 is not), got %d", entryCount)
	}
}

func TestWriteICNS_InvalidPath(t *testing.T) {
	images := map[int]*image.RGBA{16: newTestImage(16)}
	err := WriteICNS("/nonexistent/dir/test.icns", images)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestWriteICNS_FileSizeReasonable(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "reasonable.icns")

	images := map[int]*image.RGBA{
		16:  newTestImage(16),
		32:  newTestImage(32),
		128: newTestImage(128),
	}

	if err := WriteICNS(outPath, images); err != nil {
		t.Fatalf("WriteICNS() error = %v", err)
	}

	info, err := os.Stat(outPath)
	if err != nil {
		t.Fatalf("failed to stat ICNS file: %v", err)
	}

	if info.Size() < 100 {
		t.Errorf("ICNS file suspiciously small: %d bytes", info.Size())
	}

	maxReasonable := int64(10 * 1024 * 1024)
	if info.Size() > maxReasonable {
		t.Errorf("ICNS file suspiciously large: %d bytes", info.Size())
	}
}
