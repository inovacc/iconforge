package cmd

import (
	"fmt"
	"image"
	"os"

	"github.com/inovacc/iconforge/internal/favicon"
	"github.com/inovacc/iconforge/internal/icon"
	svgrender "github.com/inovacc/iconforge/internal/svg"
	"github.com/spf13/cobra"
)

var (
	faviconSVGPath   string
	faviconPNGPath   string
	faviconOutputDir string
)

var faviconCmd = &cobra.Command{
	Use:   "favicon",
	Short: "Generate web-standard favicons from SVG or PNG",
	Long: `Generate all favicon assets needed for modern web applications.

Creates favicon.ico, sized PNGs (16x16, 32x32, 180x180), Android Chrome
icons (192x192, 512x512), and a site.webmanifest file.

Examples:
  iconforge favicon --svg icon.svg -o build/favicons
  iconforge favicon --png logo.png -o build/favicons`,
	RunE: runFavicon,
}

func init() {
	rootCmd.AddCommand(faviconCmd)

	faviconCmd.Flags().StringVar(&faviconSVGPath, "svg", "", "Path to source SVG file")
	faviconCmd.Flags().StringVar(&faviconPNGPath, "png", "", "Path to source PNG file (alternative to --svg)")
	faviconCmd.Flags().StringVarP(&faviconOutputDir, "output", "o", "build/favicons", "Output directory")
}

func runFavicon(cmd *cobra.Command, _ []string) error {
	if faviconSVGPath == "" && faviconPNGPath == "" {
		return fmt.Errorf("either --svg or --png must be provided")
	}
	if faviconSVGPath != "" && faviconPNGPath != "" {
		return fmt.Errorf("--svg and --png are mutually exclusive")
	}

	if err := os.MkdirAll(faviconOutputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	sizes := favicon.Sizes()
	var images map[int]*image.RGBA

	if faviconPNGPath != "" {
		srcImg, err := icon.LoadPNG(faviconPNGPath)
		if err != nil {
			return fmt.Errorf("load png: %w", err)
		}

		images = make(map[int]*image.RGBA, len(sizes))
		for _, size := range sizes {
			images[size] = icon.ResizeImage(srcImg, size)
		}
	} else {
		var err error
		images, err = svgrender.RenderToImages(faviconSVGPath, sizes)
		if err != nil {
			return fmt.Errorf("rasterize svg: %w", err)
		}
	}

	if err := favicon.GenerateFavicons(images, faviconOutputDir); err != nil {
		return fmt.Errorf("generate favicons: %w", err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Favicons generated in: %s\n", faviconOutputDir)
	return nil
}
