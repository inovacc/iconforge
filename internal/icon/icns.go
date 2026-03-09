package icon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"
)

// ICNS icon types mapped by pixel size
var icnsTypes = map[int]string{
	16:   "icp4", // 16x16
	32:   "icp5", // 32x32
	64:   "icp6", // 64x64
	128:  "ic07", // 128x128
	256:  "ic08", // 256x256
	512:  "ic09", // 512x512
	1024: "ic10", // 1024x1024
}

// WriteICNS creates a macOS ICNS file from multiple images.
func WriteICNS(path string, images map[int]*image.RGBA) error {
	// Encode each image as PNG and build ICNS entries
	type icnsEntry struct {
		osType string
		data   []byte
	}

	var entries []icnsEntry

	for size, img := range images {
		osType, ok := icnsTypes[size]
		if !ok {
			continue // Skip unsupported sizes
		}

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return fmt.Errorf("encode png for size %d: %w", size, err)
		}

		entries = append(entries, icnsEntry{
			osType: osType,
			data:   buf.Bytes(),
		})
	}

	if len(entries) == 0 {
		return fmt.Errorf("no valid ICNS sizes provided")
	}

	// Calculate total file size
	totalSize := 8 // ICNS header
	for _, e := range entries {
		totalSize += 8 + len(e.data) // entry header + data
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create icns: %w", err)
	}
	defer func() { _ = f.Close() }()

	// Write ICNS magic header: 'icns' + total length
	if _, err := f.Write([]byte("icns")); err != nil {
		return fmt.Errorf("write magic: %w", err)
	}
	if err := binary.Write(f, binary.BigEndian, uint32(totalSize)); err != nil {
		return fmt.Errorf("write total size: %w", err)
	}

	// Write each entry
	for _, e := range entries {
		// Write OSType (4 bytes)
		if _, err := f.Write([]byte(e.osType)); err != nil {
			return fmt.Errorf("write ostype %s: %w", e.osType, err)
		}
		// Write entry length (header + data)
		if err := binary.Write(f, binary.BigEndian, uint32(8+len(e.data))); err != nil {
			return fmt.Errorf("write entry size %s: %w", e.osType, err)
		}
		// Write PNG data
		if _, err := f.Write(e.data); err != nil {
			return fmt.Errorf("write data %s: %w", e.osType, err)
		}
	}

	return nil
}
