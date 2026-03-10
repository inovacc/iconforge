package icon

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

// LoadPNG loads a PNG file and returns it as an *image.RGBA.
func LoadPNG(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open png: %w", err)
	}
	defer func() { _ = f.Close() }()

	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode png: %w", err)
	}

	if rgba, ok := img.(*image.RGBA); ok {
		return rgba, nil
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	return rgba, nil
}
