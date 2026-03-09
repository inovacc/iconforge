package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/inovacc/iconforge/internal/icon"
	svgrender "github.com/inovacc/iconforge/internal/svg"
	"github.com/spf13/cobra"
)

var (
	renderSVGPath   string
	renderOutputDir string
	renderSizesStr  string
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Rasterize SVG to PNG at multiple sizes",
	Long: `Rasterize an SVG file to PNG images at specified sizes.

Examples:
  iconforge render --svg icon.svg --sizes 256,128,64,48,32,16
  iconforge render --svg icon.svg -o build/png`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		sizes, err := parseSizes(renderSizesStr)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(renderOutputDir, 0o755); err != nil {
			return fmt.Errorf("create output dir: %w", err)
		}

		images, err := svgrender.RenderToImages(renderSVGPath, sizes)
		if err != nil {
			return fmt.Errorf("render: %w", err)
		}

		if err := icon.WritePNGs(renderOutputDir, images); err != nil {
			return fmt.Errorf("write pngs: %w", err)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Rendered %d PNGs to %s\n", len(sizes), renderOutputDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVar(&renderSVGPath, "svg", "", "Path to source SVG file")
	renderCmd.Flags().StringVarP(&renderOutputDir, "output", "o", "build/icons/png", "Output directory")
	renderCmd.Flags().StringVar(&renderSizesStr, "sizes", "256,128,64,48,32,16", "Comma-separated icon sizes")

	_ = renderCmd.MarkFlagRequired("svg")
}

func parseSizes(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	sizes := make([]int, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		size, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid size %q: %w", p, err)
		}
		if size <= 0 || size > 2048 {
			return nil, fmt.Errorf("size %d out of range (1-2048)", size)
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

// AppName returns the application name, deriving it from the working directory if empty.
func AppName(name string) string {
	if name != "" {
		return name
	}
	dir, err := os.Getwd()
	if err != nil {
		return "app"
	}
	return filepath.Base(dir)
}
