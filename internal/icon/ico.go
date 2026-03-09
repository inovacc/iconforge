package icon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"
	"sort"
)

// ICO file format structures

type icoHeader struct {
	Reserved uint16
	Type     uint16
	Count    uint16
}

type icoDirEntry struct {
	Width       byte
	Height      byte
	ColorCount  byte
	Reserved    byte
	Planes      uint16
	BitCount    uint16
	BytesInRes  uint32
	ImageOffset uint32
}

// WriteICO creates a Windows ICO file from multiple images.
// Images should be provided at different sizes (512, 256, 128, 64, 48, 32, 16).
// The 256px image is stored as PNG (standard for modern ICO files).
func WriteICO(path string, images map[int]*image.RGBA) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create ico: %w", err)
	}
	defer func() { _ = f.Close() }()

	// Sort sizes descending
	sizes := make([]int, 0, len(images))
	for size := range images {
		sizes = append(sizes, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))

	// Encode each image as PNG
	pngData := make([][]byte, len(sizes))
	for i, size := range sizes {
		var buf bytes.Buffer
		if err := png.Encode(&buf, images[size]); err != nil {
			return fmt.Errorf("encode png for size %d: %w", size, err)
		}
		pngData[i] = buf.Bytes()
	}

	// Write ICO header
	header := icoHeader{
		Reserved: 0,
		Type:     1, // ICO type
		Count:    uint16(len(sizes)),
	}
	if err := binary.Write(f, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	// Calculate offsets: header (6 bytes) + entries (16 bytes each)
	dataOffset := uint32(6 + 16*len(sizes))

	// Write directory entries
	for i, size := range sizes {
		w := byte(size)
		h := byte(size)
		if size >= 256 {
			w = 0 // 0 means 256 in ICO format
			h = 0
		}

		entry := icoDirEntry{
			Width:       w,
			Height:      h,
			ColorCount:  0,
			Reserved:    0,
			Planes:      1,
			BitCount:    32,
			BytesInRes:  uint32(len(pngData[i])),
			ImageOffset: dataOffset,
		}

		if err := binary.Write(f, binary.LittleEndian, &entry); err != nil {
			return fmt.Errorf("write dir entry: %w", err)
		}

		dataOffset += uint32(len(pngData[i]))
	}

	// Write image data
	for i, data := range pngData {
		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("write image data %d: %w", sizes[i], err)
		}
	}

	return nil
}
