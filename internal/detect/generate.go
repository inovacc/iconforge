package detect

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/inovacc/iconforge/internal/icon"
)

// GenerateFrameworkIcons generates the required icon assets for a detected framework.
func GenerateFrameworkIcons(framework Framework, projectDir string, images map[int]*image.RGBA) error {
	switch framework {
	case FrameworkTauri:
		return generateTauriIcons(projectDir, images)
	case FrameworkElectron:
		return generateElectronIcons(projectDir, images)
	case FrameworkWails:
		return generateWailsIcons(projectDir, images)
	case FrameworkFyne:
		return generateFyneIcons(projectDir, images)
	case FrameworkNone:
		return nil
	default:
		return fmt.Errorf("unsupported framework: %s", framework)
	}
}

func generateTauriIcons(projectDir string, images map[int]*image.RGBA) error {
	iconsDir := filepath.Join(projectDir, "src-tauri", "icons")
	if err := os.MkdirAll(iconsDir, 0o755); err != nil {
		return fmt.Errorf("create tauri icons dir: %w", err)
	}

	// icon.ico
	if err := icon.WriteICO(filepath.Join(iconsDir, "icon.ico"), images); err != nil {
		return fmt.Errorf("write tauri ico: %w", err)
	}

	// icon.icns
	if err := icon.WriteICNS(filepath.Join(iconsDir, "icon.icns"), images); err != nil {
		return fmt.Errorf("write tauri icns: %w", err)
	}

	// PNG files at required sizes
	tauriPNGs := map[string]int{
		"32x32.png":       32,
		"128x128.png":     128,
		"128x128@2x.png":  256,
		"icon.png":        512,
	}

	for name, size := range tauriPNGs {
		if img, ok := images[size]; ok {
			if err := icon.WritePNG(filepath.Join(iconsDir, name), img); err != nil {
				return fmt.Errorf("write tauri %s: %w", name, err)
			}
		}
	}

	return nil
}

func generateElectronIcons(projectDir string, images map[int]*image.RGBA) error {
	buildDir := filepath.Join(projectDir, "build")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return fmt.Errorf("create electron build dir: %w", err)
	}

	// icon.ico (Windows)
	if err := icon.WriteICO(filepath.Join(buildDir, "icon.ico"), images); err != nil {
		return fmt.Errorf("write electron ico: %w", err)
	}

	// icon.icns (macOS)
	if err := icon.WriteICNS(filepath.Join(buildDir, "icon.icns"), images); err != nil {
		return fmt.Errorf("write electron icns: %w", err)
	}

	// icon.png (Linux, 512x512)
	if img, ok := images[512]; ok {
		if err := icon.WritePNG(filepath.Join(buildDir, "icon.png"), img); err != nil {
			return fmt.Errorf("write electron png: %w", err)
		}
	} else if img, ok := images[256]; ok {
		if err := icon.WritePNG(filepath.Join(buildDir, "icon.png"), img); err != nil {
			return fmt.Errorf("write electron png: %w", err)
		}
	}

	return nil
}

func generateWailsIcons(projectDir string, images map[int]*image.RGBA) error {
	buildDir := filepath.Join(projectDir, "build")
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return fmt.Errorf("create wails build dir: %w", err)
	}

	// appicon.png (source icon, use largest available)
	for _, size := range []int{1024, 512, 256} {
		if img, ok := images[size]; ok {
			if err := icon.WritePNG(filepath.Join(buildDir, "appicon.png"), img); err != nil {
				return fmt.Errorf("write wails appicon: %w", err)
			}
			break
		}
	}

	// Windows ICO
	windowsDir := filepath.Join(buildDir, "windows")
	if err := os.MkdirAll(windowsDir, 0o755); err != nil {
		return fmt.Errorf("create wails windows dir: %w", err)
	}
	if err := icon.WriteICO(filepath.Join(windowsDir, "icon.ico"), images); err != nil {
		return fmt.Errorf("write wails ico: %w", err)
	}

	return nil
}

func generateFyneIcons(projectDir string, images map[int]*image.RGBA) error {
	// Fyne uses a single PNG icon, typically 512x512 or 1024x1024
	for _, size := range []int{1024, 512, 256} {
		if img, ok := images[size]; ok {
			return icon.WritePNG(filepath.Join(projectDir, "Icon.png"), img)
		}
	}
	return fmt.Errorf("no suitable icon size for Fyne")
}
