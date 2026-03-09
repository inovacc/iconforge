package svg

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// RenderToImage rasterizes an SVG file to an image.RGBA at the given size.
func RenderToImage(svgPath string, size int) (*image.RGBA, error) {
	f, err := os.Open(svgPath)
	if err != nil {
		return nil, fmt.Errorf("open svg: %w", err)
	}
	defer func() { _ = f.Close() }()

	icon, err := oksvg.ReadIconStream(f)
	if err != nil {
		return nil, fmt.Errorf("parse svg: %w", err)
	}

	w, h := float64(size), float64(size)
	icon.SetTarget(0, 0, w, h)

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill with transparent background
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.Transparent}, image.Point{}, draw.Src)

	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())
	raster := rasterx.NewDasher(size, size, scanner)
	icon.Draw(raster, 1.0)

	return img, nil
}

// RenderToImages rasterizes an SVG file to multiple sizes.
func RenderToImages(svgPath string, sizes []int) (map[int]*image.RGBA, error) {
	result := make(map[int]*image.RGBA, len(sizes))

	for _, size := range sizes {
		img, err := RenderToImage(svgPath, size)
		if err != nil {
			return nil, fmt.Errorf("render size %d: %w", size, err)
		}
		result[size] = img
	}

	return result, nil
}
