package platform

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/inovacc/iconforge/internal/icon"
)

// LinuxConfig holds configuration for Linux .desktop entry generation.
type LinuxConfig struct {
	AppName    string
	Comment    string
	Exec       string
	Categories string
	OutputDir  string
}

// CreateDesktopEntry creates a Linux .desktop file and installs icons
// into the freedesktop hicolor icon theme structure.
func CreateDesktopEntry(cfg LinuxConfig, images map[int]*image.RGBA) error {
	// Write .desktop file
	if err := writeDesktopFile(cfg); err != nil {
		return fmt.Errorf("write desktop file: %w", err)
	}

	// Write icons to hicolor theme structure
	if err := writeHicolorIcons(cfg, images); err != nil {
		return fmt.Errorf("write hicolor icons: %w", err)
	}

	return nil
}

func writeDesktopFile(cfg LinuxConfig) error {
	exec := cfg.Exec
	if exec == "" {
		exec = cfg.AppName
	}

	categories := cfg.Categories
	if categories == "" {
		categories = "Utility;"
	}

	comment := cfg.Comment
	if comment == "" {
		comment = cfg.AppName
	}

	desktop := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Comment=%s
Exec=%s
Icon=%s
Categories=%s
Terminal=false
`, cfg.AppName, comment, exec, cfg.AppName, categories)

	outPath := filepath.Join(cfg.OutputDir, cfg.AppName+".desktop")
	if err := os.MkdirAll(cfg.OutputDir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	if err := os.WriteFile(outPath, []byte(desktop), 0o644); err != nil {
		return fmt.Errorf("write desktop file: %w", err)
	}

	return nil
}

func writeHicolorIcons(cfg LinuxConfig, images map[int]*image.RGBA) error {
	iconsDir := filepath.Join(cfg.OutputDir, "icons", "hicolor")

	for size, img := range images {
		sizeDir := filepath.Join(iconsDir, fmt.Sprintf("%dx%d", size, size), "apps")
		path := filepath.Join(sizeDir, cfg.AppName+".png")

		if err := icon.WritePNG(path, img); err != nil {
			return fmt.Errorf("write %dx%d icon: %w", size, size, err)
		}
	}

	return nil
}
