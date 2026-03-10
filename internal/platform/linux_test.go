package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateDesktopEntry(t *testing.T) {
	tests := []struct {
		name           string
		cfg            LinuxConfig
		imageSizes     []int
		wantErr        bool
		wantExec       string
		wantCategories string
		wantComment    string
	}{
		{
			name: "full config",
			cfg: LinuxConfig{
				AppName:    "TestApp",
				Comment:    "A great test application",
				Exec:       "/usr/bin/testapp",
				Categories: "Development;IDE;",
			},
			imageSizes:     []int{16, 32, 48, 128, 256},
			wantExec:       "/usr/bin/testapp",
			wantCategories: "Development;IDE;",
			wantComment:    "A great test application",
		},
		{
			name: "defaults for empty fields",
			cfg: LinuxConfig{
				AppName: "MinimalApp",
			},
			imageSizes:     []int{48},
			wantExec:       "MinimalApp",
			wantCategories: "Utility;",
			wantComment:    "MinimalApp",
		},
		{
			name: "custom exec only",
			cfg: LinuxConfig{
				AppName: "MyApp",
				Exec:    "/opt/myapp/bin/run",
			},
			imageSizes:     []int{32, 64},
			wantExec:       "/opt/myapp/bin/run",
			wantCategories: "Utility;",
			wantComment:    "MyApp",
		},
		{
			name: "custom categories only",
			cfg: LinuxConfig{
				AppName:    "GameApp",
				Categories: "Game;ActionGame;",
			},
			imageSizes:     []int{128, 256, 512},
			wantExec:       "GameApp",
			wantCategories: "Game;ActionGame;",
			wantComment:    "GameApp",
		},
		{
			name: "all fields populated",
			cfg: LinuxConfig{
				AppName:    "FullApp",
				Comment:    "Full featured application",
				Exec:       "/usr/local/bin/fullapp --start",
				Categories: "Office;Finance;",
			},
			imageSizes:     []int{16, 24, 32, 48, 64, 128, 256, 512},
			wantExec:       "/usr/local/bin/fullapp --start",
			wantCategories: "Office;Finance;",
			wantComment:    "Full featured application",
		},
		{
			name: "single icon size",
			cfg: LinuxConfig{
				AppName: "SingleIcon",
			},
			imageSizes:     []int{256},
			wantExec:       "SingleIcon",
			wantCategories: "Utility;",
			wantComment:    "SingleIcon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.OutputDir = t.TempDir()
			images := makeTestImages(tt.imageSizes...)

			err := CreateDesktopEntry(tt.cfg, images)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateDesktopEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			desktopPath := filepath.Join(tt.cfg.OutputDir, tt.cfg.AppName+".desktop")
			data, err := os.ReadFile(desktopPath)
			if err != nil {
				t.Fatalf("failed to read .desktop file: %v", err)
			}

			content := string(data)

			expectedLines := map[string]string{
				"Type":       "Application",
				"Name":       tt.cfg.AppName,
				"Comment":    tt.wantComment,
				"Exec":       tt.wantExec,
				"Icon":       tt.cfg.AppName,
				"Categories": tt.wantCategories,
				"Terminal":   "false",
			}

			for key, value := range expectedLines {
				expected := key + "=" + value
				if !strings.Contains(content, expected) {
					t.Errorf(".desktop file missing %q", expected)
				}
			}

			if !strings.HasPrefix(content, "[Desktop Entry]") {
				t.Error(".desktop file does not start with [Desktop Entry]")
			}

			iconsDir := filepath.Join(tt.cfg.OutputDir, "icons", "hicolor")
			for _, size := range tt.imageSizes {
				sizeDir := filepath.Join(iconsDir, fmt.Sprintf("%dx%d", size, size), "apps")
				iconPath := filepath.Join(sizeDir, tt.cfg.AppName+".png")

				info, err := os.Stat(iconPath)
				if err != nil {
					t.Errorf("icon at %s does not exist: %v", iconPath, err)
					continue
				}
				if info.Size() == 0 {
					t.Errorf("icon at %s is empty", iconPath)
				}
			}

			for _, size := range tt.imageSizes {
				sizeDir := filepath.Join(iconsDir, fmt.Sprintf("%dx%d", size, size), "apps")
				info, err := os.Stat(sizeDir)
				if err != nil {
					t.Errorf("hicolor directory %s does not exist: %v", sizeDir, err)
					continue
				}
				if !info.IsDir() {
					t.Errorf("%s is not a directory", sizeDir)
				}
			}
		})
	}
}

func TestCreateDesktopEntry_EmptyImages(t *testing.T) {
	cfg := LinuxConfig{
		AppName:   "EmptyApp",
		OutputDir: t.TempDir(),
	}
	images := makeTestImages()

	err := CreateDesktopEntry(cfg, images)
	if err != nil {
		t.Fatalf("CreateDesktopEntry() with empty images should succeed for desktop file: %v", err)
	}

	desktopPath := filepath.Join(cfg.OutputDir, cfg.AppName+".desktop")
	if _, err := os.Stat(desktopPath); err != nil {
		t.Errorf(".desktop file should still be created: %v", err)
	}
}
