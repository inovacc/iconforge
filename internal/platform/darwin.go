package platform

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/inovacc/iconforge/internal/icon"
)

// DarwinConfig holds configuration for macOS .app bundle generation.
type DarwinConfig struct {
	AppName    string
	BundleID   string
	Version    string
	Copyright  string
	Executable string
	OutputDir  string
}

// CreateAppBundle creates a macOS .app bundle directory structure with icon.
func CreateAppBundle(cfg DarwinConfig, images map[int]*image.RGBA) error {
	bundlePath := filepath.Join(cfg.OutputDir, cfg.AppName+".app")
	contentsPath := filepath.Join(bundlePath, "Contents")
	macosPath := filepath.Join(contentsPath, "MacOS")
	resourcesPath := filepath.Join(contentsPath, "Resources")

	// Create directory structure
	for _, dir := range []string{macosPath, resourcesPath} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", dir, err)
		}
	}

	// Write ICNS
	icnsPath := filepath.Join(resourcesPath, "icon.icns")
	if err := icon.WriteICNS(icnsPath, images); err != nil {
		return fmt.Errorf("write icns: %w", err)
	}

	// Write Info.plist
	if err := writeInfoPlist(cfg, contentsPath); err != nil {
		return fmt.Errorf("write plist: %w", err)
	}

	return nil
}

// writeInfoPlist generates the Info.plist file for the .app bundle.
func writeInfoPlist(cfg DarwinConfig, contentsPath string) error {
	executable := cfg.Executable
	if executable == "" {
		executable = cfg.AppName
	}

	bundleID := cfg.BundleID
	if bundleID == "" {
		bundleID = "com.example." + cfg.AppName
	}

	version := cfg.Version
	if version == "" {
		version = "1.0.0"
	}

	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>en</string>
	<key>CFBundleExecutable</key>
	<string>%s</string>
	<key>CFBundleIconFile</key>
	<string>icon</string>
	<key>CFBundleIdentifier</key>
	<string>%s</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>%s</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>%s</string>
	<key>CFBundleVersion</key>
	<string>%s</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.13</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>NSHumanReadableCopyright</key>
	<string>%s</string>
</dict>
</plist>`, executable, bundleID, cfg.AppName, version, version, cfg.Copyright)

	outPath := filepath.Join(contentsPath, "Info.plist")
	if err := os.WriteFile(outPath, []byte(plist), 0o644); err != nil {
		return fmt.Errorf("write Info.plist: %w", err)
	}

	return nil
}
