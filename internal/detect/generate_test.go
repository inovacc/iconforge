package detect

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func makeTestImages(sizes ...int) map[int]*image.RGBA {
	images := make(map[int]*image.RGBA, len(sizes))
	for _, size := range sizes {
		img := image.NewRGBA(image.Rect(0, 0, size, size))
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}
		}
		images[size] = img
	}
	return images
}

func TestGenerateFrameworkIcons(t *testing.T) {
	tests := []struct {
		name        string
		framework   Framework
		imageSizes  []int
		wantErr     bool
		checkFiles  []string
		checkDirs   []string
	}{
		{
			name:       "Tauri generates all expected files",
			framework:  FrameworkTauri,
			imageSizes: []int{32, 128, 256, 512, 1024},
			wantErr:    false,
			checkFiles: []string{
				filepath.Join("src-tauri", "icons", "icon.ico"),
				filepath.Join("src-tauri", "icons", "icon.icns"),
				filepath.Join("src-tauri", "icons", "32x32.png"),
				filepath.Join("src-tauri", "icons", "128x128.png"),
				filepath.Join("src-tauri", "icons", "128x128@2x.png"),
				filepath.Join("src-tauri", "icons", "icon.png"),
			},
			checkDirs: []string{
				filepath.Join("src-tauri", "icons"),
			},
		},
		{
			name:       "Electron generates all expected files",
			framework:  FrameworkElectron,
			imageSizes: []int{16, 32, 48, 64, 128, 256, 512, 1024},
			wantErr:    false,
			checkFiles: []string{
				filepath.Join("build", "icon.ico"),
				filepath.Join("build", "icon.icns"),
				filepath.Join("build", "icon.png"),
			},
			checkDirs: []string{
				"build",
			},
		},
		{
			name:       "Electron uses 256 fallback when no 512",
			framework:  FrameworkElectron,
			imageSizes: []int{16, 32, 256},
			wantErr:    false,
			checkFiles: []string{
				filepath.Join("build", "icon.ico"),
				filepath.Join("build", "icon.icns"),
				filepath.Join("build", "icon.png"),
			},
		},
		{
			name:       "Wails generates expected files",
			framework:  FrameworkWails,
			imageSizes: []int{32, 128, 256, 512, 1024},
			wantErr:    false,
			checkFiles: []string{
				filepath.Join("build", "appicon.png"),
				filepath.Join("build", "windows", "icon.ico"),
			},
			checkDirs: []string{
				"build",
				filepath.Join("build", "windows"),
			},
		},
		{
			name:       "Wails uses 512 when 1024 not available",
			framework:  FrameworkWails,
			imageSizes: []int{32, 128, 256, 512},
			wantErr:    false,
			checkFiles: []string{
				filepath.Join("build", "appicon.png"),
			},
		},
		{
			name:       "Fyne generates Icon.png",
			framework:  FrameworkFyne,
			imageSizes: []int{256, 512, 1024},
			wantErr:    false,
			checkFiles: []string{
				"Icon.png",
			},
		},
		{
			name:       "Fyne uses 512 when 1024 not available",
			framework:  FrameworkFyne,
			imageSizes: []int{256, 512},
			wantErr:    false,
			checkFiles: []string{
				"Icon.png",
			},
		},
		{
			name:       "Fyne uses 256 when larger sizes not available",
			framework:  FrameworkFyne,
			imageSizes: []int{256},
			wantErr:    false,
			checkFiles: []string{
				"Icon.png",
			},
		},
		{
			name:       "Fyne errors when no suitable size",
			framework:  FrameworkFyne,
			imageSizes: []int{32, 64},
			wantErr:    true,
		},
		{
			name:       "FrameworkNone does nothing",
			framework:  FrameworkNone,
			imageSizes: []int{32, 128, 256, 512},
			wantErr:    false,
		},
		{
			name:       "Unknown framework returns error",
			framework:  Framework(99),
			imageSizes: []int{32, 128, 256, 512},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			images := makeTestImages(tt.imageSizes...)

			err := GenerateFrameworkIcons(tt.framework, dir, images)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateFrameworkIcons() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			for _, rel := range tt.checkFiles {
				path := filepath.Join(dir, rel)
				info, err := os.Stat(path)
				if err != nil {
					t.Errorf("expected file %s to exist, got error: %v", rel, err)
					continue
				}
				if info.IsDir() {
					t.Errorf("expected %s to be a file, got directory", rel)
					continue
				}
				if info.Size() == 0 {
					t.Errorf("expected %s to be non-empty", rel)
				}
			}

			for _, rel := range tt.checkDirs {
				path := filepath.Join(dir, rel)
				info, err := os.Stat(path)
				if err != nil {
					t.Errorf("expected directory %s to exist, got error: %v", rel, err)
					continue
				}
				if !info.IsDir() {
					t.Errorf("expected %s to be a directory, got file", rel)
				}
			}
		})
	}
}

func TestGenerateTauriIconsPartialSizes(t *testing.T) {
	dir := t.TempDir()
	images := makeTestImages(32, 128)

	err := GenerateFrameworkIcons(FrameworkTauri, dir, images)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	iconsDir := filepath.Join(dir, "src-tauri", "icons")
	mustExist := []string{"icon.ico", "icon.icns", "32x32.png", "128x128.png"}
	for _, name := range mustExist {
		if !fileExists(filepath.Join(iconsDir, name)) {
			t.Errorf("expected %s to exist in icons dir", name)
		}
	}

	shouldNotExist := []string{"128x128@2x.png", "icon.png"}
	for _, name := range shouldNotExist {
		if fileExists(filepath.Join(iconsDir, name)) {
			t.Errorf("expected %s to NOT exist (no 256/512 image provided)", name)
		}
	}
}

func TestGenerateElectronIconsNoLargePNG(t *testing.T) {
	dir := t.TempDir()
	images := makeTestImages(16, 32)

	err := GenerateFrameworkIcons(FrameworkElectron, dir, images)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fileExists(filepath.Join(dir, "build", "icon.png")) {
		t.Error("expected icon.png to NOT exist when no 256 or 512 image provided")
	}

	if !fileExists(filepath.Join(dir, "build", "icon.ico")) {
		t.Error("expected icon.ico to exist")
	}
}
