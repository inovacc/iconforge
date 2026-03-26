package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// TemplateInfo describes a registered icon template.
type TemplateInfo struct {
	Name        string
	Description string
	GenerateFn  func(outputPath, appName string, palette ColorPalette) error
}

var registry = make(map[string]TemplateInfo)

// RegisterTemplate adds a template to the global registry.
func RegisterTemplate(info TemplateInfo) {
	registry[info.Name] = info
}

// GetTemplate retrieves a template by name.
func GetTemplate(name string) (TemplateInfo, error) {
	t, ok := registry[name]
	if !ok {
		names := make([]string, 0, len(registry))
		for n := range registry {
			names = append(names, n)
		}
		sort.Strings(names)
		return TemplateInfo{}, fmt.Errorf("unknown template %q; available: %v", name, names)
	}
	return t, nil
}

// ListTemplates returns all registered templates sorted by name.
func ListTemplates() []TemplateInfo {
	list := make([]TemplateInfo, 0, len(registry))
	for _, t := range registry {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name < list[j].Name })
	return list
}

// writeSVG is a helper that writes an SVG string to outputPath, creating directories.
func writeSVG(outputPath, svg string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}
	if err := os.WriteFile(outputPath, []byte(svg), 0o644); err != nil {
		return fmt.Errorf("write svg: %w", err)
	}
	return nil
}

// svgHeader returns the common rounded-rect background header for all templates.
func svgHeader(p ColorPalette) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
  <defs>
    <linearGradient id="bg" x1="0" y1="0" x2="512" y2="512" gradientUnits="userSpaceOnUse">
      <stop offset="0" stop-color="%s"/>
      <stop offset="1" stop-color="%s"/>
    </linearGradient>
  </defs>
  <rect x="16" y="16" width="480" height="480" rx="96" ry="96" fill="url(#bg)"/>`, p.Primary, p.Secondary)
}

const svgFooter = "\n</svg>"

func init() {
	RegisterTemplate(TemplateInfo{
		Name:        "shield",
		Description: "Shield with keyhole — security, VPN, auth apps",
		GenerateFn:  generateShield,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "terminal",
		Description: "Terminal prompt — CLI tools, dev utilities",
		GenerateFn:  generateTerminal,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "gear",
		Description: "Mechanical gear — system utilities, settings",
		GenerateFn:  generateGear,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "cube",
		Description: "Isometric cube — data, 3D, containers",
		GenerateFn:  generateCube,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "bolt",
		Description: "Lightning bolt — speed, performance, power",
		GenerateFn:  generateBolt,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "leaf",
		Description: "Leaf — eco, nature, organic, growth",
		GenerateFn:  generateLeaf,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "wave",
		Description: "Wave pattern — streaming, audio, data flow",
		GenerateFn:  generateWave,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "hexagon",
		Description: "Hexagonal grid — science, tech, molecular",
		GenerateFn:  generateHexagon,
	})
	RegisterTemplate(TemplateInfo{
		Name:        "stack",
		Description: "Stacked layers — infrastructure, DevOps, platforms",
		GenerateFn:  generateStack,
	})
}

func generateShield(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Shield body -->
  <path d="M256,72 C256,72 370,92 420,108 C424,112 424,118 424,118 L424,248 C424,320 380,380 256,440 C132,380 88,320 88,248 L88,118 C88,112 92,108 92,108 C142,92 256,72 256,72 Z"
    fill="#FFFFFF" fill-opacity="0.12" stroke="#FFFFFF" stroke-opacity="0.5" stroke-width="3"/>

  <!-- Shield inner -->
  <path d="M256,100 C256,100 352,116 396,130 L396,244 C396,306 358,358 256,410 C154,358 116,306 116,244 L116,130 C160,116 256,100 256,100 Z"
    fill="#FFFFFF" fill-opacity="0.06"/>

  <!-- Keyhole circle -->
  <circle cx="256" cy="210" r="40" fill="%s" fill-opacity="0.9"/>

  <!-- Keyhole inner -->
  <circle cx="256" cy="210" r="24" fill="#FFFFFF" fill-opacity="0.85"/>

  <!-- Keyhole slot -->
  <path d="M244,224 L256,210 L268,224 L264,300 L248,300 Z" fill="%s" fill-opacity="0.9"/>

  <!-- Shield highlight -->
  <path d="M256,100 C256,100 352,116 396,130 L396,160 C340,148 256,134 256,134 C256,134 172,148 116,160 L116,130 C160,116 256,100 256,100 Z"
    fill="#FFFFFF" fill-opacity="0.1"/>`, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateTerminal(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Terminal window -->
  <rect x="80" y="100" width="352" height="312" rx="20" ry="20"
    fill="#FFFFFF" fill-opacity="0.1" stroke="#FFFFFF" stroke-opacity="0.4" stroke-width="2"/>

  <!-- Title bar -->
  <rect x="80" y="100" width="352" height="40" rx="20" ry="0"
    fill="#FFFFFF" fill-opacity="0.08"/>
  <line x1="80" y1="140" x2="432" y2="140" stroke="#FFFFFF" stroke-opacity="0.2" stroke-width="1"/>

  <!-- Window dots -->
  <circle cx="110" cy="120" r="8" fill="#E53E3E" fill-opacity="0.8"/>
  <circle cx="136" cy="120" r="8" fill="#F59E0B" fill-opacity="0.8"/>
  <circle cx="162" cy="120" r="8" fill="#38B249" fill-opacity="0.8"/>

  <!-- Prompt chevron -->
  <polyline points="120,200 160,240 120,280"
    fill="none" stroke="%s" stroke-width="12" stroke-linecap="round" stroke-linejoin="round"/>

  <!-- Cursor line -->
  <line x1="190" y1="228" x2="370" y2="228"
    stroke="%s" stroke-opacity="0.8" stroke-width="8" stroke-linecap="round"/>

  <!-- Dimmed text lines -->
  <line x1="190" y1="280" x2="320" y2="280"
    stroke="#FFFFFF" stroke-opacity="0.2" stroke-width="6" stroke-linecap="round"/>
  <line x1="190" y1="320" x2="280" y2="320"
    stroke="#FFFFFF" stroke-opacity="0.15" stroke-width="6" stroke-linecap="round"/>
  <line x1="190" y1="360" x2="350" y2="360"
    stroke="#FFFFFF" stroke-opacity="0.1" stroke-width="6" stroke-linecap="round"/>`, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateGear(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Gear teeth (12 teeth as rectangles rotated via transform) -->
  <polygon points="244,76 268,76 272,116 240,116" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="244,396 268,396 272,436 240,436" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="76,244 76,268 116,272 116,240" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="396,244 396,268 436,272 436,240" fill="#FFFFFF" fill-opacity="0.7"/>

  <polygon points="327,99 345,115 321,149 299,127" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="167,397 185,413 161,379 183,363" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="397,327 413,345 379,321 363,299" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="99,167 115,185 149,161 127,183" fill="#FFFFFF" fill-opacity="0.7"/>

  <polygon points="99,345 115,327 149,351 127,329" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="397,185 413,167 379,191 363,213" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="345,397 327,413 351,379 329,363" fill="#FFFFFF" fill-opacity="0.7"/>
  <polygon points="185,99 167,115 191,149 213,127" fill="#FFFFFF" fill-opacity="0.7"/>

  <!-- Gear body -->
  <circle cx="256" cy="256" r="140" fill="#FFFFFF" fill-opacity="0.12" stroke="#FFFFFF" stroke-opacity="0.5" stroke-width="3"/>

  <!-- Inner ring -->
  <circle cx="256" cy="256" r="100" fill="#FFFFFF" fill-opacity="0.06" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>

  <!-- Center hub -->
  <circle cx="256" cy="256" r="50" fill="%s" fill-opacity="0.9"/>

  <!-- Center dot -->
  <circle cx="256" cy="256" r="20" fill="#FFFFFF" fill-opacity="0.85"/>`, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateCube(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Cube - isometric projection -->
  <!-- Top face -->
  <polygon points="256,96 416,192 256,288 96,192"
    fill="%s" fill-opacity="0.8"/>

  <!-- Top face highlight -->
  <polygon points="256,96 416,192 256,288 96,192"
    fill="#FFFFFF" fill-opacity="0.2"/>

  <!-- Left face -->
  <polygon points="96,192 256,288 256,416 96,320"
    fill="%s" fill-opacity="0.6"/>

  <!-- Right face -->
  <polygon points="416,192 256,288 256,416 416,320"
    fill="%s" fill-opacity="0.4"/>

  <!-- Edge lines -->
  <line x1="256" y1="96" x2="256" y2="288" stroke="#FFFFFF" stroke-opacity="0.1" stroke-width="1"/>
  <line x1="96" y1="192" x2="256" y2="288" stroke="#FFFFFF" stroke-opacity="0.15" stroke-width="1"/>
  <line x1="416" y1="192" x2="256" y2="288" stroke="#FFFFFF" stroke-opacity="0.15" stroke-width="1"/>

  <!-- Outline -->
  <polygon points="256,96 416,192 416,320 256,416 96,320 96,192"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.5" stroke-width="2"/>

  <!-- Center vertical edge -->
  <line x1="256" y1="288" x2="256" y2="416" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>`, p.Accent, p.Primary, p.Secondary) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateBolt(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Outer glow circle -->
  <circle cx="256" cy="256" r="180" fill="#FFFFFF" fill-opacity="0.04"/>
  <circle cx="256" cy="256" r="140" fill="#FFFFFF" fill-opacity="0.04"/>

  <!-- Lightning bolt -->
  <polygon points="280,68 160,268 236,268 200,444 360,220 276,220 320,68"
    fill="%s" fill-opacity="0.95"/>

  <!-- Bolt inner highlight -->
  <polygon points="278,108 188,258 244,258 216,404 340,234 280,234 312,108"
    fill="#FFFFFF" fill-opacity="0.25"/>

  <!-- Spark lines -->
  <line x1="140" y1="180" x2="108" y2="156" stroke="%s" stroke-opacity="0.5" stroke-width="3" stroke-linecap="round"/>
  <line x1="380" y1="320" x2="412" y2="344" stroke="%s" stroke-opacity="0.5" stroke-width="3" stroke-linecap="round"/>
  <line x1="128" y1="320" x2="100" y2="340" stroke="%s" stroke-opacity="0.3" stroke-width="2" stroke-linecap="round"/>
  <line x1="400" y1="180" x2="424" y2="160" stroke="%s" stroke-opacity="0.3" stroke-width="2" stroke-linecap="round"/>`, p.Accent, p.Accent, p.Accent, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateLeaf(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Leaf body -->
  <path d="M256,80 C380,120 440,220 420,340 C400,420 340,460 256,440 C172,460 112,420 92,340 C72,220 132,120 256,80 Z"
    fill="%s" fill-opacity="0.8" stroke="%s" stroke-width="2"/>

  <!-- Leaf vein - center -->
  <path d="M256,110 C256,110 256,200 256,420"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.5" stroke-width="3" stroke-linecap="round"/>

  <!-- Leaf veins - left -->
  <path d="M256,200 C230,180 180,180 140,200"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2" stroke-linecap="round"/>
  <path d="M256,260 C230,240 170,240 120,270"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.25" stroke-width="2" stroke-linecap="round"/>
  <path d="M256,320 C230,305 180,310 140,340"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.2" stroke-width="2" stroke-linecap="round"/>

  <!-- Leaf veins - right -->
  <path d="M256,200 C282,180 332,180 372,200"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2" stroke-linecap="round"/>
  <path d="M256,260 C282,240 342,240 392,270"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.25" stroke-width="2" stroke-linecap="round"/>
  <path d="M256,320 C282,305 332,310 372,340"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.2" stroke-width="2" stroke-linecap="round"/>

  <!-- Stem -->
  <path d="M256,420 C256,420 270,450 280,460"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.5" stroke-width="3" stroke-linecap="round"/>`, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateWave(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Wave 1 (back) -->
  <path d="M60,320 C120,260 180,260 240,320 C300,380 360,380 452,320"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.15" stroke-width="8" stroke-linecap="round"/>

  <!-- Wave 2 -->
  <path d="M60,280 C120,220 180,220 240,280 C300,340 360,340 452,280"
    fill="none" stroke="#FFFFFF" stroke-opacity="0.25" stroke-width="8" stroke-linecap="round"/>

  <!-- Wave 3 (main) -->
  <path d="M60,240 C120,180 180,180 240,240 C300,300 360,300 452,240"
    fill="none" stroke="%s" stroke-opacity="0.9" stroke-width="10" stroke-linecap="round"/>

  <!-- Wave 4 -->
  <path d="M60,200 C120,140 180,140 240,200 C300,260 360,260 452,200"
    fill="none" stroke="%s" stroke-opacity="0.6" stroke-width="8" stroke-linecap="round"/>

  <!-- Wave 5 (front) -->
  <path d="M60,160 C120,100 180,100 240,160 C300,220 360,220 452,160"
    fill="none" stroke="%s" stroke-opacity="0.3" stroke-width="6" stroke-linecap="round"/>

  <!-- Dot accents -->
  <circle cx="240" cy="240" r="6" fill="%s"/>
  <circle cx="140" cy="200" r="4" fill="#FFFFFF" fill-opacity="0.5"/>
  <circle cx="360" cy="260" r="4" fill="#FFFFFF" fill-opacity="0.5"/>`, p.Accent, p.Accent, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateHexagon(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Center hexagon -->
  <polygon points="256,96 396,176 396,336 256,416 116,336 116,176"
    fill="%s" fill-opacity="0.3" stroke="%s" stroke-width="3"/>

  <!-- Inner hexagon -->
  <polygon points="256,156 336,204 336,300 256,348 176,300 176,204"
    fill="%s" fill-opacity="0.2" stroke="%s" stroke-opacity="0.6" stroke-width="2"/>

  <!-- Core hexagon -->
  <polygon points="256,208 296,232 296,280 256,304 216,280 216,232"
    fill="%s" fill-opacity="0.8"/>

  <!-- Connecting lines to outer nodes -->
  <line x1="256" y1="96" x2="256" y2="156" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>
  <line x1="396" y1="176" x2="336" y2="204" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>
  <line x1="396" y1="336" x2="336" y2="300" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>
  <line x1="256" y1="416" x2="256" y2="348" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>
  <line x1="116" y1="336" x2="176" y2="300" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>
  <line x1="116" y1="176" x2="176" y2="204" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>

  <!-- Vertex dots -->
  <circle cx="256" cy="96" r="8" fill="#FFFFFF" fill-opacity="0.7"/>
  <circle cx="396" cy="176" r="8" fill="#FFFFFF" fill-opacity="0.7"/>
  <circle cx="396" cy="336" r="8" fill="#FFFFFF" fill-opacity="0.7"/>
  <circle cx="256" cy="416" r="8" fill="#FFFFFF" fill-opacity="0.7"/>
  <circle cx="116" cy="336" r="8" fill="#FFFFFF" fill-opacity="0.7"/>
  <circle cx="116" cy="176" r="8" fill="#FFFFFF" fill-opacity="0.7"/>`, p.Accent, p.Accent, p.Accent, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}

func generateStack(outputPath, _ string, p ColorPalette) error {
	svg := svgHeader(p) + fmt.Sprintf(`

  <!-- Layer 3 (bottom) -->
  <polygon points="256,380 420,300 420,340 256,420 92,340 92,300"
    fill="#FFFFFF" fill-opacity="0.15" stroke="#FFFFFF" stroke-opacity="0.3" stroke-width="2"/>

  <!-- Layer 2 (middle) -->
  <polygon points="256,300 420,220 420,260 256,340 92,260 92,220"
    fill="#FFFFFF" fill-opacity="0.2" stroke="#FFFFFF" stroke-opacity="0.4" stroke-width="2"/>

  <!-- Layer 1 (top) -->
  <polygon points="256,220 420,140 420,180 256,260 92,180 92,140"
    fill="%s" fill-opacity="0.7" stroke="%s" stroke-width="2"/>

  <!-- Top face of layer 1 -->
  <polygon points="256,140 420,60 420,140 256,220 92,140 92,60"
    fill="%s" fill-opacity="0.5" stroke="%s" stroke-opacity="0.6" stroke-width="2"/>

  <!-- Top face highlight -->
  <polygon points="256,140 420,60 256,220 92,60"
    fill="#FFFFFF" fill-opacity="0.08"/>

  <!-- Center line -->
  <line x1="256" y1="140" x2="256" y2="420" stroke="#FFFFFF" stroke-opacity="0.15" stroke-width="1"/>

  <!-- Arrow up on top layer -->
  <polygon points="256,90 280,130 268,130 268,170 244,170 244,130 232,130"
    fill="#FFFFFF" fill-opacity="0.7"/>`, p.Accent, p.Accent, p.Accent, p.Accent) + svgFooter

	return writeSVG(outputPath, svg)
}
