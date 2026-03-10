package platform

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeTestImages(sizes ...int) map[int]*image.RGBA {
	images := make(map[int]*image.RGBA, len(sizes))
	for _, size := range sizes {
		img := image.NewRGBA(image.Rect(0, 0, size, size))
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				idx := (y*size + x) * 4
				img.Pix[idx] = 255   // R
				img.Pix[idx+1] = 0   // G
				img.Pix[idx+2] = 0   // B
				img.Pix[idx+3] = 255 // A
			}
		}
		images[size] = img
	}
	return images
}

func TestCreateAppBundle(t *testing.T) {
	tests := []struct {
		name           string
		cfg            DarwinConfig
		imageSizes     []int
		wantErr        bool
		wantExecutable string
		wantBundleID   string
		wantVersion    string
	}{
		{
			name: "full config",
			cfg: DarwinConfig{
				AppName:    "TestApp",
				BundleID:   "com.test.testapp",
				Version:    "2.1.0",
				Copyright:  "Copyright 2026 Test Corp",
				Executable: "testapp-bin",
			},
			imageSizes:     []int{16, 32, 128, 256, 512},
			wantExecutable: "testapp-bin",
			wantBundleID:   "com.test.testapp",
			wantVersion:    "2.1.0",
		},
		{
			name: "defaults for empty fields",
			cfg: DarwinConfig{
				AppName: "MinimalApp",
			},
			imageSizes:     []int{128},
			wantExecutable: "MinimalApp",
			wantBundleID:   "com.example.MinimalApp",
			wantVersion:    "1.0.0",
		},
		{
			name: "custom executable only",
			cfg: DarwinConfig{
				AppName:    "MyApp",
				Executable: "myapp-runner",
			},
			imageSizes:     []int{32, 64},
			wantExecutable: "myapp-runner",
			wantBundleID:   "com.example.MyApp",
			wantVersion:    "1.0.0",
		},
		{
			name: "custom bundle ID only",
			cfg: DarwinConfig{
				AppName:  "AnotherApp",
				BundleID: "org.custom.another",
			},
			imageSizes:     []int{256},
			wantExecutable: "AnotherApp",
			wantBundleID:   "org.custom.another",
			wantVersion:    "1.0.0",
		},
		{
			name: "all fields populated",
			cfg: DarwinConfig{
				AppName:    "FullApp",
				BundleID:   "io.full.app",
				Version:    "3.5.7",
				Copyright:  "All rights reserved",
				Executable: "full-runner",
			},
			imageSizes:     []int{16, 32, 64, 128, 256, 512, 1024},
			wantExecutable: "full-runner",
			wantBundleID:   "io.full.app",
			wantVersion:    "3.5.7",
		},
		{
			name: "single small icon",
			cfg: DarwinConfig{
				AppName: "TinyApp",
				Version: "0.1.0",
			},
			imageSizes:     []int{16},
			wantExecutable: "TinyApp",
			wantBundleID:   "com.example.TinyApp",
			wantVersion:    "0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.OutputDir = t.TempDir()
			images := makeTestImages(tt.imageSizes...)

			err := CreateAppBundle(tt.cfg, images)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CreateAppBundle() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			bundlePath := filepath.Join(tt.cfg.OutputDir, tt.cfg.AppName+".app")
			contentsPath := filepath.Join(bundlePath, "Contents")
			macosPath := filepath.Join(contentsPath, "MacOS")
			resourcesPath := filepath.Join(contentsPath, "Resources")

			for _, dir := range []string{bundlePath, contentsPath, macosPath, resourcesPath} {
				info, err := os.Stat(dir)
				if err != nil {
					t.Fatalf("directory %s does not exist: %v", dir, err)
				}
				if !info.IsDir() {
					t.Errorf("%s is not a directory", dir)
				}
			}

			icnsPath := filepath.Join(resourcesPath, "icon.icns")
			info, err := os.Stat(icnsPath)
			if err != nil {
				t.Fatalf("icon.icns does not exist: %v", err)
			}
			if info.Size() == 0 {
				t.Error("icon.icns is empty")
			}

			plistPath := filepath.Join(contentsPath, "Info.plist")
			data, err := os.ReadFile(plistPath)
			if err != nil {
				t.Fatalf("failed to read Info.plist: %v", err)
			}

			content := string(data)

			if !strings.Contains(content, `<?xml version="1.0"`) {
				t.Error("Info.plist missing XML declaration")
			}
			if !strings.Contains(content, `<plist version="1.0">`) {
				t.Error("Info.plist missing plist version declaration")
			}

			expectedStrings := map[string]string{
				"CFBundleExecutable":       tt.wantExecutable,
				"CFBundleIdentifier":       tt.wantBundleID,
				"CFBundleName":             tt.cfg.AppName,
				"CFBundleShortVersionString": tt.wantVersion,
				"CFBundleVersion":          tt.wantVersion,
			}

			for key, value := range expectedStrings {
				if !strings.Contains(content, "<key>"+key+"</key>") {
					t.Errorf("Info.plist missing key %q", key)
				}
				if !strings.Contains(content, "<string>"+value+"</string>") {
					t.Errorf("Info.plist missing value %q for key %q", value, key)
				}
			}

			if !strings.Contains(content, "<key>CFBundleIconFile</key>") {
				t.Error("Info.plist missing CFBundleIconFile key")
			}
			if !strings.Contains(content, "<string>icon</string>") {
				t.Error("Info.plist missing icon value for CFBundleIconFile")
			}

			if !strings.Contains(content, "<key>CFBundlePackageType</key>") {
				t.Error("Info.plist missing CFBundlePackageType")
			}
			if !strings.Contains(content, "<string>APPL</string>") {
				t.Error("Info.plist missing APPL package type")
			}

			if !strings.Contains(content, "<key>NSHighResolutionCapable</key>") {
				t.Error("Info.plist missing NSHighResolutionCapable")
			}

			if tt.cfg.Copyright != "" {
				if !strings.Contains(content, "<string>"+tt.cfg.Copyright+"</string>") {
					t.Errorf("Info.plist missing copyright %q", tt.cfg.Copyright)
				}
			}
		})
	}
}

func TestCreateIconset(t *testing.T) {
	images := makeTestImages(16, 32, 64, 128, 256, 512, 1024)
	outputDir := t.TempDir()

	iconsetPath, err := CreateIconset(outputDir, "TestApp", images)
	if err != nil {
		t.Fatalf("CreateIconset() error = %v", err)
	}

	wantDir := filepath.Join(outputDir, "TestApp.iconset")
	if iconsetPath != wantDir {
		t.Errorf("iconsetPath = %q, want %q", iconsetPath, wantDir)
	}

	info, err := os.Stat(iconsetPath)
	if err != nil {
		t.Fatalf("iconset dir does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("iconset path is not a directory")
	}

	expectedFiles := map[string]int{
		"icon_16x16.png":      16,
		"icon_16x16@2x.png":   32,
		"icon_32x32.png":      32,
		"icon_32x32@2x.png":   64,
		"icon_128x128.png":    128,
		"icon_128x128@2x.png": 256,
		"icon_256x256.png":    256,
		"icon_256x256@2x.png": 512,
		"icon_512x512.png":    512,
		"icon_512x512@2x.png": 1024,
	}

	for filename, wantSize := range expectedFiles {
		path := filepath.Join(iconsetPath, filename)
		f, err := os.Open(path)
		if err != nil {
			t.Errorf("missing file %s: %v", filename, err)
			continue
		}

		img, err := png.Decode(f)
		_ = f.Close()
		if err != nil {
			t.Errorf("failed to decode %s: %v", filename, err)
			continue
		}

		bounds := img.Bounds()
		if bounds.Dx() != wantSize || bounds.Dy() != wantSize {
			t.Errorf("%s: got %dx%d, want %dx%d", filename, bounds.Dx(), bounds.Dy(), wantSize, wantSize)
		}
	}
}

func TestCreateAppBundle_InvalidOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a file to block directory creation
	blockingFile := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to write blocking file: %v", err)
	}

	cfg := DarwinConfig{
		AppName:   "FailApp",
		OutputDir: filepath.Join(blockingFile, "subdir"),
	}
	images := makeTestImages(128)

	err := CreateAppBundle(cfg, images)
	if err == nil {
		t.Error("expected error for invalid output dir, got nil")
	}
}

func TestCreateAppBundle_UnsupportedICNSSizes(t *testing.T) {
	cfg := DarwinConfig{
		AppName:   "NoICNSApp",
		OutputDir: t.TempDir(),
	}
	// Only unsupported ICNS sizes (48 is not in icnsTypes map)
	images := makeTestImages(48)

	err := CreateAppBundle(cfg, images)
	if err == nil {
		t.Error("expected error when all image sizes are unsupported for ICNS, got nil")
	}
}

func TestCreateIconset_InvalidOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	blockingFile := filepath.Join(tmpDir, "blocker")
	if err := os.WriteFile(blockingFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to write blocking file: %v", err)
	}

	images := makeTestImages(512)
	_, err := CreateIconset(filepath.Join(blockingFile, "sub"), "App", images)
	if err == nil {
		t.Error("expected error for invalid output dir, got nil")
	}
}

func TestWriteInfoPlist_AllDefaults(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := DarwinConfig{
		AppName: "DefaultApp",
	}

	err := writeInfoPlist(cfg, tmpDir)
	if err != nil {
		t.Fatalf("writeInfoPlist() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "Info.plist"))
	if err != nil {
		t.Fatalf("failed to read plist: %v", err)
	}

	content := string(data)

	// Check defaults
	if !strings.Contains(content, "<string>DefaultApp</string>") {
		t.Error("expected default executable to be AppName")
	}
	if !strings.Contains(content, "<string>com.example.DefaultApp</string>") {
		t.Error("expected default bundleID")
	}
	if !strings.Contains(content, "<string>1.0.0</string>") {
		t.Error("expected default version 1.0.0")
	}
}

func TestWriteInfoPlist_InvalidOutputDir(t *testing.T) {
	cfg := DarwinConfig{
		AppName: "FailApp",
	}

	err := writeInfoPlist(cfg, filepath.Join(t.TempDir(), "nonexistent", "deep"))
	if err == nil {
		t.Error("expected error for invalid output dir, got nil")
	}
}

func TestCreateIconset_MissingSizes(t *testing.T) {
	// Only provide 512 — everything else must be resized.
	images := makeTestImages(512)
	outputDir := t.TempDir()

	iconsetPath, err := CreateIconset(outputDir, "ResizeApp", images)
	if err != nil {
		t.Fatalf("CreateIconset() error = %v", err)
	}

	expectedFiles := map[string]int{
		"icon_16x16.png":      16,
		"icon_16x16@2x.png":   32,
		"icon_32x32.png":      32,
		"icon_32x32@2x.png":   64,
		"icon_128x128.png":    128,
		"icon_128x128@2x.png": 256,
		"icon_256x256.png":    256,
		"icon_256x256@2x.png": 512,
		"icon_512x512.png":    512,
		"icon_512x512@2x.png": 1024,
	}

	for filename, wantSize := range expectedFiles {
		path := filepath.Join(iconsetPath, filename)
		f, err := os.Open(path)
		if err != nil {
			t.Errorf("missing file %s: %v", filename, err)
			continue
		}

		img, err := png.Decode(f)
		_ = f.Close()
		if err != nil {
			t.Errorf("failed to decode %s: %v", filename, err)
			continue
		}

		bounds := img.Bounds()
		if bounds.Dx() != wantSize || bounds.Dy() != wantSize {
			t.Errorf("%s: got %dx%d, want %dx%d", filename, bounds.Dx(), bounds.Dy(), wantSize, wantSize)
		}
	}
}
