package favicon

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/inovacc/iconforge/internal/icon"
)

// faviconSizes lists all sizes needed for favicon generation.
var faviconSizes = []int{16, 32, 180, 192, 512}

// webManifest represents the structure of a site.webmanifest file.
type webManifest struct {
	Name  string            `json:"name"`
	Icons []webManifestIcon `json:"icons"`
}

type webManifestIcon struct {
	Src   string `json:"src"`
	Sizes string `json:"sizes"`
	Type  string `json:"type"`
}

// Sizes returns the list of image sizes needed for favicon generation.
func Sizes() []int {
	return append([]int(nil), faviconSizes...)
}

// GenerateFavicons creates all web-standard favicon files in outputDir.
// The images map should contain pre-rendered images keyed by pixel size.
// If a required size is missing, the largest available image is resized.
func GenerateFavicons(images map[int]*image.RGBA, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create favicon dir: %w", err)
	}

	// Ensure all required sizes are available, resizing from the largest if needed.
	resolved := resolveImages(images)

	// favicon.ico — multi-size ICO with 16x16 and 32x32
	icoImages := map[int]*image.RGBA{
		16: resolved[16],
		32: resolved[32],
	}
	icoPath := filepath.Join(outputDir, "favicon.ico")
	if err := icon.WriteICO(icoPath, icoImages); err != nil {
		return fmt.Errorf("write favicon.ico: %w", err)
	}

	// Individual PNG files
	pngFiles := map[string]int{
		"favicon-16x16.png":          16,
		"favicon-32x32.png":          32,
		"apple-touch-icon.png":       180,
		"android-chrome-192x192.png": 192,
		"android-chrome-512x512.png": 512,
	}

	for filename, size := range pngFiles {
		path := filepath.Join(outputDir, filename)
		if err := icon.WritePNG(path, resolved[size]); err != nil {
			return fmt.Errorf("write %s: %w", filename, err)
		}
	}

	// site.webmanifest
	manifest := webManifest{
		Name: "App",
		Icons: []webManifestIcon{
			{
				Src:   "android-chrome-192x192.png",
				Sizes: "192x192",
				Type:  "image/png",
			},
			{
				Src:   "android-chrome-512x512.png",
				Sizes: "512x512",
				Type:  "image/png",
			},
		},
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal webmanifest: %w", err)
	}

	manifestPath := filepath.Join(outputDir, "site.webmanifest")
	if err := os.WriteFile(manifestPath, manifestData, 0o644); err != nil {
		return fmt.Errorf("write site.webmanifest: %w", err)
	}

	return nil
}

// resolveImages ensures all required favicon sizes are present.
// Missing sizes are created by resizing the largest available image.
func resolveImages(images map[int]*image.RGBA) map[int]*image.RGBA {
	// Find the largest available image
	var largest *image.RGBA
	var largestSize int
	for size, img := range images {
		if size > largestSize {
			largestSize = size
			largest = img
		}
	}

	resolved := make(map[int]*image.RGBA, len(faviconSizes))
	for _, size := range faviconSizes {
		if img, ok := images[size]; ok {
			resolved[size] = img
		} else if largest != nil {
			resolved[size] = icon.ResizeImage(largest, size)
		}
	}

	return resolved
}
