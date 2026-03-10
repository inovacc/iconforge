package cmd

import (
	"fmt"
	"image"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/inovacc/iconforge/internal/detect"
	"github.com/inovacc/iconforge/internal/generator"
	"github.com/inovacc/iconforge/internal/icon"
	"github.com/inovacc/iconforge/internal/platform"
	svgrender "github.com/inovacc/iconforge/internal/svg"
	"github.com/spf13/cobra"
)

var (
	forgeSVGPath    string
	forgePNGPath    string
	forgeOutputDir  string
	forgeAppName    string
	forgeVersion    string
	forgeCompany    string
	forgeCopyright  string
	forgeBundleID   string
	forgeArch       string
	forgeGenSVG     bool
	forgePrimary    string
	forgeSecondary  string
	forgeAccent     string
	forgeSkipWin    bool
	forgeSkipMac    bool
	forgeSkipLinux  bool
	forgeAutoDetect bool
)

// Standard icon sizes for export
var standardSizes = []int{512, 256, 128, 64, 48, 32, 16}

// Extended sizes including retina/high-res
var extendedSizes = []int{1024, 512, 256, 128, 64, 48, 32, 16}

var forgeCmd = &cobra.Command{
	Use:   "forge",
	Short: "Generate all icons from SVG source",
	Long: `Generate production-ready icons for all platforms from an SVG source.

Creates ICO (Windows), ICNS (macOS), PNG (Linux), and framework-specific
icon assets. Optionally generates a modern abstract gradient SVG icon.

Examples:
  iconforge forge --svg icon.svg --name myapp
  iconforge forge --generate --name myapp
  iconforge forge --svg icon.svg --name myapp --auto-detect`,
	RunE: runForge,
}

func init() {
	rootCmd.AddCommand(forgeCmd)

	forgeCmd.Flags().StringVar(&forgeSVGPath, "svg", "", "Path to source SVG file")
	forgeCmd.Flags().StringVar(&forgePNGPath, "from-png", "", "Path to source PNG file (alternative to --svg)")
	forgeCmd.Flags().StringVarP(&forgeOutputDir, "output", "o", "build/icons", "Output directory")
	forgeCmd.Flags().StringVar(&forgeAppName, "name", "", "Application name (auto-detected from directory if empty)")
	forgeCmd.Flags().StringVar(&forgeVersion, "version", "1.0.0", "Application version")
	forgeCmd.Flags().StringVar(&forgeCompany, "company", "", "Company name for Windows metadata")
	forgeCmd.Flags().StringVar(&forgeCopyright, "copyright", "", "Copyright notice")
	forgeCmd.Flags().StringVar(&forgeBundleID, "bundle-id", "", "macOS bundle identifier (e.g., com.example.app)")
	forgeCmd.Flags().StringVar(&forgeArch, "arch", "amd64", "Target architecture for .syso (amd64, 386, arm64)")

	forgeCmd.Flags().BoolVar(&forgeGenSVG, "generate", false, "Generate a modern abstract gradient SVG icon")
	forgeCmd.Flags().StringVar(&forgePrimary, "primary", "#4F46E5", "Primary gradient color (hex)")
	forgeCmd.Flags().StringVar(&forgeSecondary, "secondary", "#7C3AED", "Secondary gradient color (hex)")
	forgeCmd.Flags().StringVar(&forgeAccent, "accent", "#F59E0B", "Accent color (hex)")

	forgeCmd.Flags().BoolVar(&forgeSkipWin, "skip-windows", false, "Skip Windows icon generation")
	forgeCmd.Flags().BoolVar(&forgeSkipMac, "skip-macos", false, "Skip macOS icon generation")
	forgeCmd.Flags().BoolVar(&forgeSkipLinux, "skip-linux", false, "Skip Linux icon generation")
	forgeCmd.Flags().BoolVar(&forgeAutoDetect, "auto-detect", false, "Auto-detect and generate framework-specific icons")
}

func runForge(cmd *cobra.Command, _ []string) error {
	logger := slog.Default()

	// Determine app name
	appName := forgeAppName
	if appName == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}
		appName = filepath.Base(dir)
	}

	// Create output directory
	if err := os.MkdirAll(forgeOutputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Validate mutually exclusive flags
	if forgePNGPath != "" && forgeSVGPath != "" {
		return fmt.Errorf("--from-png and --svg are mutually exclusive")
	}
	if forgePNGPath != "" && forgeGenSVG {
		return fmt.Errorf("--from-png and --generate are mutually exclusive")
	}

	svgPath := forgeSVGPath

	// Generate SVG if requested
	if forgeGenSVG {
		palette := generator.ColorPalette{
			Primary:   forgePrimary,
			Secondary: forgeSecondary,
			Accent:    forgeAccent,
		}

		svgPath = filepath.Join(forgeOutputDir, appName+".svg")
		logger.Info("generating SVG icon", "path", svgPath)

		if err := generator.GenerateIconSVG(svgPath, appName, palette); err != nil {
			return fmt.Errorf("generate svg: %w", err)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated SVG: %s\n", svgPath)
	}

	var images map[int]*image.RGBA

	if forgePNGPath != "" {
		// Load PNG and resize to all extended sizes
		logger.Info("loading PNG source", "path", forgePNGPath)
		srcImg, err := icon.LoadPNG(forgePNGPath)
		if err != nil {
			return fmt.Errorf("load png: %w", err)
		}

		images = make(map[int]*image.RGBA, len(extendedSizes))
		for _, size := range extendedSizes {
			images[size] = icon.ResizeImage(srcImg, size)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Resized PNG to %d sizes\n", len(images))
	} else {
		if svgPath == "" {
			return fmt.Errorf("no source specified; use --svg, --from-png, or --generate")
		}

		// Rasterize SVG to multiple sizes
		logger.Info("rasterizing SVG", "sizes", extendedSizes)
		var err error
		images, err = svgrender.RenderToImages(svgPath, extendedSizes)
		if err != nil {
			return fmt.Errorf("rasterize svg: %w", err)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Rasterized to %d sizes\n", len(images))
	}

	// Export PNGs
	pngDir := filepath.Join(forgeOutputDir, "png")
	if err := icon.WritePNGs(pngDir, images); err != nil {
		return fmt.Errorf("write pngs: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "PNGs exported to: %s\n", pngDir)

	// Windows
	if !forgeSkipWin {
		if err := forgeWindows(cmd, appName, images); err != nil {
			return err
		}
	}

	// macOS
	if !forgeSkipMac {
		if err := forgeMacOS(cmd, appName, images); err != nil {
			return err
		}
	}

	// Linux
	if !forgeSkipLinux {
		if err := forgeLinux(cmd, appName, images); err != nil {
			return err
		}
	}

	// Framework auto-detection
	if forgeAutoDetect {
		projectDir, _ := os.Getwd()
		fw := detect.DetectFramework(projectDir)

		if fw != detect.FrameworkNone {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Detected framework: %s\n", fw)

			if err := detect.GenerateFrameworkIcons(fw, projectDir, images); err != nil {
				return fmt.Errorf("generate framework icons: %w", err)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Framework icons generated for %s\n", fw)
		}
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nForge complete!")
	return nil
}

func forgeWindows(cmd *cobra.Command, appName string, images map[int]*image.RGBA) error {
	winDir := filepath.Join(forgeOutputDir, "windows")
	if err := os.MkdirAll(winDir, 0o755); err != nil {
		return fmt.Errorf("create windows dir: %w", err)
	}

	// Write ICO
	icoPath := filepath.Join(winDir, "icon.ico")
	if err := icon.WriteICO(icoPath, images); err != nil {
		return fmt.Errorf("write ico: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Windows ICO: %s\n", icoPath)

	// Write versioninfo.json
	cfg := platform.WindowsConfig{
		AppName:     appName,
		Description: appName,
		Version:     forgeVersion,
		Company:     forgeCompany,
		Copyright:   forgeCopyright,
		ICOPath:     "icon.ico",
		OutputDir:   winDir,
	}

	viPath, err := platform.WriteVersionInfo(cfg)
	if err != nil {
		return fmt.Errorf("write versioninfo: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Windows manifest: %s\n", viPath)

	// Write app manifest
	manifestPath, err := platform.WriteManifest(appName, winDir)
	if err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Windows app manifest: %s\n", manifestPath)

	// Try to generate .syso using rsrc
	sysoPath := filepath.Join(winDir, "rsrc_windows_"+forgeArch+".syso")
	if err := platform.GenerateSyso(icoPath, sysoPath, forgeArch); err != nil {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Note: rsrc not available (install with: go install github.com/akavel/rsrc@latest)\n")
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "  You can generate .syso manually: rsrc -ico %s -o %s -arch %s\n", icoPath, sysoPath, forgeArch)
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Windows .syso: %s\n", sysoPath)
	}

	return nil
}

func forgeMacOS(cmd *cobra.Command, appName string, images map[int]*image.RGBA) error {
	macDir := filepath.Join(forgeOutputDir, "macos")

	// Write standalone ICNS
	icnsPath := filepath.Join(macDir, "icon.icns")
	if err := os.MkdirAll(macDir, 0o755); err != nil {
		return fmt.Errorf("create macos dir: %w", err)
	}

	if err := icon.WriteICNS(icnsPath, images); err != nil {
		return fmt.Errorf("write icns: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "macOS ICNS: %s\n", icnsPath)

	// Create .app bundle
	cfg := platform.DarwinConfig{
		AppName:    appName,
		BundleID:   forgeBundleID,
		Version:    forgeVersion,
		Copyright:  forgeCopyright,
		Executable: appName,
		OutputDir:  macDir,
	}

	if err := platform.CreateAppBundle(cfg, images); err != nil {
		return fmt.Errorf("create app bundle: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "macOS .app bundle: %s/%s.app\n", macDir, appName)

	return nil
}

func forgeLinux(cmd *cobra.Command, appName string, images map[int]*image.RGBA) error {
	linuxDir := filepath.Join(forgeOutputDir, "linux")

	cfg := platform.LinuxConfig{
		AppName:    appName,
		Comment:    appName,
		Exec:       appName,
		Categories: "Utility;",
		OutputDir:  linuxDir,
	}

	if err := platform.CreateDesktopEntry(cfg, images); err != nil {
		return fmt.Errorf("create linux desktop entry: %w", err)
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Linux .desktop: %s/%s.desktop\n", linuxDir, appName)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Linux icons: %s/icons/hicolor/\n", linuxDir)

	return nil
}
