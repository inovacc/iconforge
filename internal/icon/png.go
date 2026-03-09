package icon

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// WritePNG writes an image to a PNG file.
func WritePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create png: %w", err)
	}
	defer func() { _ = f.Close() }()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	return nil
}

// WritePNGs writes multiple images to PNG files in a directory.
// Files are named as "{size}x{size}.png".
func WritePNGs(dir string, images map[int]*image.RGBA) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	for size, img := range images {
		filename := fmt.Sprintf("%dx%d.png", size, size)
		path := filepath.Join(dir, filename)

		if err := WritePNG(path, img); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
	}

	return nil
}
