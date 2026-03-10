package icon

import (
	"image"

	"golang.org/x/image/draw"
)

// ResizeImage resizes an RGBA image to the given target size using bilinear interpolation.
func ResizeImage(src *image.RGBA, targetSize int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}
