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

func TestCreateDesktopEntry_ErrorOnInvalidOutputDir(t *testing.T) {
	// Use a path where directory creation would fail (file as parent)
	tmpDir := t.TempDir()
	// Create a file where we need a directory
	blockingFile := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to write blocking file: %v", err)
	}

	cfg := LinuxConfig{
		AppName:   "FailApp",
		OutputDir: filepath.Join(blockingFile, "subdir"),
	}
	images := makeTestImages(32)

	err := CreateDesktopEntry(cfg, images)
	if err == nil {
		t.Error("expected error when output dir cannot be created, got nil")
	}
}

func TestCreateDesktopEntry_HicolorIconError(t *testing.T) {
	// Create a valid output dir so the desktop file succeeds,
	// but make the icons subdirectory creation fail
	tmpDir := t.TempDir()

	// Write a file at the path where "icons" dir is needed
	blockingFile := filepath.Join(tmpDir, "icons")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to write blocking file: %v", err)
	}

	cfg := LinuxConfig{
		AppName:   "IconFailApp",
		OutputDir: tmpDir,
	}
	images := makeTestImages(32)

	err := CreateDesktopEntry(cfg, images)
	if err == nil {
		t.Error("expected error when icon directory cannot be created, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "hicolor") {
		t.Errorf("error should mention hicolor icons, got: %v", err)
	}
}

func TestWriteHicolorIcons_MultipleSizes(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := LinuxConfig{
		AppName:   "MultiApp",
		OutputDir: tmpDir,
	}
	images := makeTestImages(16, 32, 48, 64, 128, 256)

	err := writeHicolorIcons(cfg, images)
	if err != nil {
		t.Fatalf("writeHicolorIcons() error = %v", err)
	}

	for size := range images {
		iconPath := filepath.Join(tmpDir, "icons", "hicolor",
			fmt.Sprintf("%dx%d", size, size), "apps", "MultiApp.png")
		if _, err := os.Stat(iconPath); os.IsNotExist(err) {
			t.Errorf("missing icon at %s", iconPath)
		}
	}
}

func TestWriteDesktopFile_AllDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := LinuxConfig{
		AppName:   "DefaultApp",
		OutputDir: tmpDir,
	}

	err := writeDesktopFile(cfg)
	if err != nil {
		t.Fatalf("writeDesktopFile() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "DefaultApp.desktop"))
	if err != nil {
		t.Fatalf("failed to read desktop file: %v", err)
	}

	content := string(data)

	// All defaults should be applied
	if !strings.Contains(content, "Exec=DefaultApp") {
		t.Error("expected Exec default to be AppName")
	}
	if !strings.Contains(content, "Categories=Utility;") {
		t.Error("expected Categories default to be Utility;")
	}
	if !strings.Contains(content, "Comment=DefaultApp") {
		t.Error("expected Comment default to be AppName")
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
