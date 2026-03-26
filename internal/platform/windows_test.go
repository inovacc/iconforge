package platform

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantMajor     int
		wantMinor     int
		wantPatch     int
	}{
		{
			name:      "standard semver",
			input:     "1.2.3",
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
		},
		{
			name:      "v-prefixed semver",
			input:     "v1.2.3",
			wantMajor: 1,
			wantMinor: 2,
			wantPatch: 3,
		},
		{
			name:      "non-numeric version string",
			input:     "dev",
			wantMajor: 0,
			wantMinor: 0,
			wantPatch: 0,
		},
		{
			name:      "empty string",
			input:     "",
			wantMajor: 0,
			wantMinor: 0,
			wantPatch: 0,
		},
		{
			name:      "major only",
			input:     "5.0.0",
			wantMajor: 5,
			wantMinor: 0,
			wantPatch: 0,
		},
		{
			name:      "high version numbers",
			input:     "10.20.30",
			wantMajor: 10,
			wantMinor: 20,
			wantPatch: 30,
		},
		{
			name:      "v-prefixed with zeros",
			input:     "v0.0.1",
			wantMajor: 0,
			wantMinor: 0,
			wantPatch: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, minor, patch := parseVersion(tt.input)
			if major != tt.wantMajor || minor != tt.wantMinor || patch != tt.wantPatch {
				t.Errorf("parseVersion(%q) = (%d, %d, %d), want (%d, %d, %d)",
					tt.input, major, minor, patch, tt.wantMajor, tt.wantMinor, tt.wantPatch)
			}
		})
	}
}

func TestWriteVersionInfo(t *testing.T) {
	tests := []struct {
		name    string
		cfg     WindowsConfig
		wantErr bool
	}{
		{
			name: "full config",
			cfg: WindowsConfig{
				AppName:     "TestApp",
				Description: "A test application",
				Version:     "1.0.0",
				Company:     "Test Corp",
				Copyright:   "Copyright 2026 Test Corp",
				ICOPath:     "icon.ico",
			},
		},
		{
			name: "v-prefixed version",
			cfg: WindowsConfig{
				AppName:     "MyApp",
				Description: "My application",
				Version:     "v1.2.3",
				Company:     "My Company",
				Copyright:   "Copyright 2026",
				ICOPath:     "app.ico",
			},
		},
		{
			name: "empty version",
			cfg: WindowsConfig{
				AppName:     "EmptyVer",
				Description: "No version",
				Version:     "",
				Company:     "",
				Copyright:   "",
				ICOPath:     "icon.ico",
			},
		},
		{
			name: "empty fields",
			cfg: WindowsConfig{
				AppName: "Minimal",
				ICOPath: "icon.ico",
			},
		},
		{
			name: "dev version string",
			cfg: WindowsConfig{
				AppName:     "DevApp",
				Description: "Development build",
				Version:     "dev",
				Company:     "Dev Inc",
				Copyright:   "",
				ICOPath:     "dev.ico",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.OutputDir = t.TempDir()

			outPath, err := WriteVersionInfo(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteVersionInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if filepath.Base(outPath) != "versioninfo.json" {
				t.Errorf("expected output file named versioninfo.json, got %s", filepath.Base(outPath))
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}

			var vi VersionInfo
			if err := json.Unmarshal(data, &vi); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			if vi.StringFileInfo.CompanyName != tt.cfg.Company {
				t.Errorf("CompanyName = %q, want %q", vi.StringFileInfo.CompanyName, tt.cfg.Company)
			}
			if vi.StringFileInfo.FileDescription != tt.cfg.Description {
				t.Errorf("FileDescription = %q, want %q", vi.StringFileInfo.FileDescription, tt.cfg.Description)
			}
			if vi.StringFileInfo.FileVersion != tt.cfg.Version {
				t.Errorf("FileVersion = %q, want %q", vi.StringFileInfo.FileVersion, tt.cfg.Version)
			}
			if vi.StringFileInfo.InternalName != tt.cfg.AppName {
				t.Errorf("InternalName = %q, want %q", vi.StringFileInfo.InternalName, tt.cfg.AppName)
			}
			if vi.StringFileInfo.LegalCopyright != tt.cfg.Copyright {
				t.Errorf("LegalCopyright = %q, want %q", vi.StringFileInfo.LegalCopyright, tt.cfg.Copyright)
			}
			if vi.StringFileInfo.OriginalFilename != tt.cfg.AppName+".exe" {
				t.Errorf("OriginalFilename = %q, want %q", vi.StringFileInfo.OriginalFilename, tt.cfg.AppName+".exe")
			}
			if vi.StringFileInfo.ProductName != tt.cfg.AppName {
				t.Errorf("ProductName = %q, want %q", vi.StringFileInfo.ProductName, tt.cfg.AppName)
			}
			if vi.StringFileInfo.ProductVersion != tt.cfg.Version {
				t.Errorf("ProductVersion = %q, want %q", vi.StringFileInfo.ProductVersion, tt.cfg.Version)
			}
			if vi.IconPath != tt.cfg.ICOPath {
				t.Errorf("IconPath = %q, want %q", vi.IconPath, tt.cfg.ICOPath)
			}

			wantMajor, wantMinor, wantPatch := parseVersion(tt.cfg.Version)
			if vi.FixedFileInfo.FileVersion.Major != wantMajor {
				t.Errorf("FileVersion.Major = %d, want %d", vi.FixedFileInfo.FileVersion.Major, wantMajor)
			}
			if vi.FixedFileInfo.FileVersion.Minor != wantMinor {
				t.Errorf("FileVersion.Minor = %d, want %d", vi.FixedFileInfo.FileVersion.Minor, wantMinor)
			}
			if vi.FixedFileInfo.FileVersion.Patch != wantPatch {
				t.Errorf("FileVersion.Patch = %d, want %d", vi.FixedFileInfo.FileVersion.Patch, wantPatch)
			}
			if vi.FixedFileInfo.FileVersion.Build != 0 {
				t.Errorf("FileVersion.Build = %d, want 0", vi.FixedFileInfo.FileVersion.Build)
			}

			if vi.VarFileInfo.Translation.LangID != "0409" {
				t.Errorf("LangID = %q, want %q", vi.VarFileInfo.Translation.LangID, "0409")
			}
			if vi.VarFileInfo.Translation.CharsetID != "04B0" {
				t.Errorf("CharsetID = %q, want %q", vi.VarFileInfo.Translation.CharsetID, "04B0")
			}
		})
	}
}

func TestWriteVersionInfo_InvalidOutputDir(t *testing.T) {
	cfg := WindowsConfig{
		AppName:   "TestApp",
		Version:   "1.0.0",
		ICOPath:   "icon.ico",
		OutputDir: filepath.Join(t.TempDir(), "nonexistent", "deeply", "nested"),
	}

	_, err := WriteVersionInfo(cfg)
	if err == nil {
		t.Error("expected error for invalid output directory, got nil")
	}
}

func TestWriteManifest(t *testing.T) {
	tests := []struct {
		name    string
		appName string
		wantErr bool
	}{
		{
			name:    "standard app name",
			appName: "TestApp",
		},
		{
			name:    "app name with spaces",
			appName: "My Application",
		},
		{
			name:    "single character name",
			appName: "a",
		},
		{
			name:    "app name with dots",
			appName: "com.example.app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir := t.TempDir()

			outPath, err := WriteManifest(tt.appName, outputDir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			expectedFilename := tt.appName + ".exe.manifest"
			if filepath.Base(outPath) != expectedFilename {
				t.Errorf("expected filename %q, got %q", expectedFilename, filepath.Base(outPath))
			}

			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("failed to read manifest: %v", err)
			}

			content := string(data)

			expectedElements := []string{
				"assemblyIdentity",
				"requestedExecutionLevel",
				"dpiAware",
				"dpiAwareness",
				"supportedOS",
				"trustInfo",
			}

			for _, elem := range expectedElements {
				if !strings.Contains(content, elem) {
					t.Errorf("manifest missing expected element %q", elem)
				}
			}

			if !strings.Contains(content, `name="`+tt.appName+`"`) {
				t.Errorf("manifest does not contain app name %q in assemblyIdentity", tt.appName)
			}

			if !strings.Contains(content, `level="asInvoker"`) {
				t.Error("manifest missing asInvoker execution level")
			}

			if !strings.Contains(content, `<dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true/pm</dpiAware>`) {
				t.Error("manifest missing dpiAware element with correct value")
			}

			if !strings.Contains(content, `PerMonitorV2`) {
				t.Error("manifest missing PerMonitorV2 dpiAwareness")
			}

			if !strings.Contains(content, `<?xml version="1.0"`) {
				t.Error("manifest missing XML declaration")
			}
		})
	}
}

func TestWriteManifest_InvalidOutputDir(t *testing.T) {
	_, err := WriteManifest("TestApp", filepath.Join(t.TempDir(), "nonexistent", "deeply", "nested"))
	if err == nil {
		t.Error("expected error for invalid output directory, got nil")
	}
}

func TestWriteManifest_AllSupportedOSGuids(t *testing.T) {
	outputDir := t.TempDir()
	outPath, err := WriteManifest("GuidApp", outputDir)
	if err != nil {
		t.Fatalf("WriteManifest() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}

	content := string(data)

	guids := []struct {
		id      string
		version string
	}{
		{"{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}", "Windows 10"},
		{"{1f676c76-80e1-4239-95bb-83d0f6d0da78}", "Windows 7"},
		{"{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}", "Windows Vista"},
		{"{35138b9a-5d96-4fbd-8e2d-a2440225f93a}", "Windows 8"},
		{"{e2011457-1546-43c5-a5fe-008deee3d3f0}", "Windows 8.1"},
	}

	for _, g := range guids {
		if !strings.Contains(content, g.id) {
			t.Errorf("manifest missing supportedOS GUID for %s: %s", g.version, g.id)
		}
	}
}

func TestGenerateSysoWinres(t *testing.T) {
	tests := []struct {
		name    string
		cfg     WindowsConfig
		arch    string
		wantErr bool
	}{
		{
			name: "basic config amd64",
			cfg: WindowsConfig{
				AppName:     "TestApp",
				Description: "Test application",
				Version:     "1.0.0",
				Company:     "Test Corp",
				Copyright:   "Copyright 2026",
			},
			arch: "amd64",
		},
		{
			name: "default arch",
			cfg: WindowsConfig{
				AppName: "DefaultArch",
				Version: "2.0.0",
			},
			arch: "",
		},
		{
			name: "386 arch",
			cfg: WindowsConfig{
				AppName: "X86App",
				Version: "1.0.0",
			},
			arch: "386",
		},
		{
			name: "v-prefixed version",
			cfg: WindowsConfig{
				AppName: "VPrefixed",
				Version: "v3.2.1",
			},
			arch: "amd64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cfg.OutputDir = t.TempDir()
			tt.cfg.ICOPath = ""

			outPath, err := GenerateSysoWinres(tt.cfg, tt.arch)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateSysoWinres() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if _, err := os.Stat(outPath); err != nil {
				t.Fatalf(".syso file was not created: %v", err)
			}

			info, err := os.Stat(outPath)
			if err != nil {
				t.Fatalf("stat: %v", err)
			}
			if info.Size() == 0 {
				t.Error(".syso file is empty")
			}

			expectedArch := tt.arch
			if expectedArch == "" {
				expectedArch = "amd64"
			}
			expectedName := "resource_windows_" + expectedArch + ".syso"
			if filepath.Base(outPath) != expectedName {
				t.Errorf("output filename = %q, want %q", filepath.Base(outPath), expectedName)
			}
		})
	}
}

func TestGenerateSysoWinres_WithManifest(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := WindowsConfig{
		AppName:     "ManifestApp",
		Description: "App with manifest",
		Version:     "1.0.0",
		OutputDir:   tmpDir,
		ICOPath:     "",
	}

	manifestPath := filepath.Join(tmpDir, "ManifestApp.exe.manifest")
	if err := os.WriteFile(manifestPath, []byte(`<?xml version="1.0"?><assembly/>`), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	outPath, err := GenerateSysoWinres(cfg, "amd64")
	if err != nil {
		t.Fatalf("GenerateSysoWinres() error = %v", err)
	}

	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf(".syso not created: %v", err)
	}
}

func TestGenerateSysoFromICO(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "resource.syso")

	// GenerateSysoFromICO requires a valid ICO file
	// With no ICO file, it should fail
	err := GenerateSysoFromICO("/nonexistent/icon.ico", outputPath, "amd64")
	if err == nil {
		t.Error("expected error for nonexistent ICO, got nil")
	}
}

func TestGenerateSysoFromICO_DefaultArch(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "resource.syso")

	err := GenerateSysoFromICO("/nonexistent/icon.ico", outputPath, "")
	if err == nil {
		t.Error("expected error for nonexistent ICO, got nil")
	}
}

func TestWriteVersionInfo_VersionQuadFields(t *testing.T) {
	cfg := WindowsConfig{
		AppName:   "QuadApp",
		Version:   "v3.7.12",
		ICOPath:   "icon.ico",
		OutputDir: t.TempDir(),
	}

	outPath, err := WriteVersionInfo(cfg)
	if err != nil {
		t.Fatalf("WriteVersionInfo() error = %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var vi VersionInfo
	if err := json.Unmarshal(data, &vi); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if vi.FixedFileInfo.ProductVersion != vi.FixedFileInfo.FileVersion {
		t.Errorf("ProductVersion quad %+v != FileVersion quad %+v",
			vi.FixedFileInfo.ProductVersion, vi.FixedFileInfo.FileVersion)
	}

	if vi.FixedFileInfo.FileFlagsMask != "3f" {
		t.Errorf("FileFlagsMask = %q, want %q", vi.FixedFileInfo.FileFlagsMask, "3f")
	}
	if vi.FixedFileInfo.FileFlags != "00" {
		t.Errorf("FileFlags = %q, want %q", vi.FixedFileInfo.FileFlags, "00")
	}
	if vi.FixedFileInfo.FileOS != "040004" {
		t.Errorf("FileOS = %q, want %q", vi.FixedFileInfo.FileOS, "040004")
	}
	if vi.FixedFileInfo.FileType != "01" {
		t.Errorf("FileType = %q, want %q", vi.FixedFileInfo.FileType, "01")
	}
	if vi.FixedFileInfo.FileSubType != "00" {
		t.Errorf("FileSubType = %q, want %q", vi.FixedFileInfo.FileSubType, "00")
	}
}
