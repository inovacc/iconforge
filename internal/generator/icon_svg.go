package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// ColorPalette defines the colors used in icon generation.
type ColorPalette struct {
	Primary   string // Main brand color (hex)
	Secondary string // Complementary color (hex)
	Accent    string // Highlight color (hex)
}

// DefaultPalette returns the default IconForge color palette.
func DefaultPalette() ColorPalette {
	return ColorPalette{
		Primary:   "#4F46E5", // Deep Indigo
		Secondary: "#7C3AED", // Violet
		Accent:    "#F59E0B", // Electric Amber
	}
}

func init() {
	RegisterTemplate(TemplateInfo{
		Name:        "forge",
		Description: "Diamond facets with forge flame — the classic IconForge default",
		GenerateFn:  GenerateIconSVG,
	})
}

// GenerateIconSVG creates a modern abstract gradient SVG icon.
// Uses only oksvg-compatible elements (linearGradient, basic shapes, paths).
// The design features nested diamond facets with a forge-flame center.
func GenerateIconSVG(outputPath string, appName string, palette ColorPalette) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	svg := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="512" y2="512" gradientUnits="userSpaceOnUse">
      <stop offset="0" stop-color="%s"/>
      <stop offset="1" stop-color="%s"/>
    </linearGradient>
    <linearGradient id="flame" x1="256" y1="160" x2="256" y2="320" gradientUnits="userSpaceOnUse">
      <stop offset="0" stop-color="%s"/>
      <stop offset="1" stop-color="#F97316"/>
    </linearGradient>
  </defs>

  <!-- Background rounded rect -->
  <rect x="16" y="16" width="480" height="480" rx="96" ry="96" fill="url(#bg)"/>

  <!-- Outer diamond -->
  <polygon points="256,56 456,256 256,456 56,256"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.15" stroke-width="2"/>

  <!-- Middle diamond -->
  <polygon points="256,106 406,256 256,406 106,256"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.2" stroke-width="2"/>

  <!-- Inner diamond -->
  <polygon points="256,156 356,256 256,356 156,256"
    fill="#FFFFFF" fill-opacity="0.08" stroke="#FFFFFF" stroke-opacity="0.25" stroke-width="2"/>

  <!-- Core diamond -->
  <polygon points="256,196 316,256 256,316 196,256"
    fill="#FFFFFF" fill-opacity="0.12" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>

  <!-- Anvil base -->
  <path d="M208,286 L218,276 L298,276 L308,286 L304,296 L212,296 Z"
    fill="#FFFFFF" fill-opacity="0.9"/>

  <!-- Flame -->
  <path d="M256,184 C274,208 286,226 286,246 C286,266 274,278 256,278 C238,278 226,266 226,246 C226,226 238,208 256,184 Z"
    fill="url(#flame)"/>

  <!-- Inner flame highlight -->
  <path d="M256,208 C266,222 274,234 274,244 C274,254 266,260 256,260 C246,260 238,254 238,244 C238,234 246,222 256,208 Z"
    fill="#FFFFFF" fill-opacity="0.8"/>

  <!-- Spark dots -->
  <circle cx="230" cy="210" r="3" fill="%s" fill-opacity="0.8"/>
  <circle cx="280" cy="200" r="2.5" fill="%s" fill-opacity="0.6"/>
  <circle cx="240" cy="192" r="2" fill="#FFFFFF" fill-opacity="0.7"/>
  <circle cx="285" cy="218" r="1.5" fill="#FFFFFF" fill-opacity="0.5"/>
</svg>`, palette.Primary, palette.Secondary, palette.Accent, palette.Accent, palette.Accent)

	if err := os.WriteFile(outputPath, []byte(svg), 0o644); err != nil {
		return fmt.Errorf("write svg: %w", err)
	}

	return nil
}
