package cmd

import (
	"fmt"
	"image"
	"io"
)

// previewIcon renders a small ANSI-color preview of the icon to the terminal.
// It uses Unicode half-block characters (▀) to pack two pixel rows per line,
// with 24-bit ANSI color codes. The image is downscaled to fit previewSize columns.
func previewIcon(w io.Writer, img *image.RGBA, previewSize int) {
	if previewSize <= 0 {
		previewSize = 32
	}

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Scale to fit previewSize columns
	scale := float64(srcW) / float64(previewSize)
	cols := previewSize
	rows := int(float64(srcH) / scale)

	// Ensure even number of rows for half-block rendering
	if rows%2 != 0 {
		rows++
	}

	_, _ = fmt.Fprintln(w)
	for y := 0; y < rows; y += 2 {
		for x := 0; x < cols; x++ {
			// Sample top pixel
			sx := int(float64(x) * scale)
			sy := int(float64(y) * scale)
			tr, tg, tb, ta := samplePixel(img, sx, sy, srcW, srcH)

			// Sample bottom pixel
			sy2 := int(float64(y+1) * scale)
			br, bg, bb, ba := samplePixel(img, sx, sy2, srcW, srcH)

			if ta < 128 && ba < 128 {
				_, _ = fmt.Fprint(w, " ")
			} else if ta < 128 {
				_, _ = fmt.Fprintf(w, "\033[38;2;%d;%d;%dm▄\033[0m", br, bg, bb)
			} else if ba < 128 {
				_, _ = fmt.Fprintf(w, "\033[38;2;%d;%d;%dm▀\033[0m", tr, tg, tb)
			} else {
				_, _ = fmt.Fprintf(w, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm▀\033[0m", tr, tg, tb, br, bg, bb)
			}
		}
		_, _ = fmt.Fprintln(w)
	}
	_, _ = fmt.Fprintln(w)
}

func samplePixel(img *image.RGBA, x, y, maxW, maxH int) (r, g, b, a uint8) {
	if x >= maxW {
		x = maxW - 1
	}
	if y >= maxH {
		y = maxH - 1
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B, c.A
}
