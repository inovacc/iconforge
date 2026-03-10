package platform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/josephspurrier/goversioninfo"
)

// VersionInfo represents the Windows VERSIONINFO resource structure
// compatible with goversioninfo's versioninfo.json format.
type VersionInfo struct {
	FixedFileInfo FixedFileInfo `json:"FixedFileInfo"`
	StringFileInfo StringFileInfo `json:"StringFileInfo"`
	VarFileInfo    VarFileInfo    `json:"VarFileInfo"`
	IconPath       string         `json:"IconPath"`
	ManifestPath   string         `json:"ManifestPath,omitempty"`
}

// FixedFileInfo contains version numbers.
type FixedFileInfo struct {
	FileVersion    VersionQuad `json:"FileVersion"`
	ProductVersion VersionQuad `json:"ProductVersion"`
	FileFlagsMask  string      `json:"FileFlagsMask"`
	FileFlags      string      `json:"FileFlags"`
	FileOS         string      `json:"FileOS"`
	FileType       string      `json:"FileType"`
	FileSubType    string      `json:"FileSubType"`
}

// VersionQuad represents a four-part version number.
type VersionQuad struct {
	Major int `json:"Major"`
	Minor int `json:"Minor"`
	Patch int `json:"Patch"`
	Build int `json:"Build"`
}

// StringFileInfo contains string metadata.
type StringFileInfo struct {
	Comments         string `json:"Comments"`
	CompanyName      string `json:"CompanyName"`
	FileDescription  string `json:"FileDescription"`
	FileVersion      string `json:"FileVersion"`
	InternalName     string `json:"InternalName"`
	LegalCopyright   string `json:"LegalCopyright"`
	LegalTrademarks  string `json:"LegalTrademarks"`
	OriginalFilename string `json:"OriginalFilename"`
	PrivateBuild     string `json:"PrivateBuild"`
	ProductName      string `json:"ProductName"`
	ProductVersion   string `json:"ProductVersion"`
	SpecialBuild     string `json:"SpecialBuild"`
}

// VarFileInfo contains translation info.
type VarFileInfo struct {
	Translation Translation `json:"Translation"`
}

// Translation contains language and charset IDs.
type Translation struct {
	LangID    string `json:"LangID"`
	CharsetID string `json:"CharsetID"`
}

// WindowsConfig holds configuration for Windows icon embedding.
type WindowsConfig struct {
	AppName     string
	Description string
	Version     string
	Company     string
	Copyright   string
	ICOPath     string
	OutputDir   string
}

// WriteVersionInfo generates a versioninfo.json file compatible with goversioninfo.
func WriteVersionInfo(cfg WindowsConfig) (string, error) {
	major, minor, patch := parseVersion(cfg.Version)

	vi := VersionInfo{
		FixedFileInfo: FixedFileInfo{
			FileVersion:    VersionQuad{Major: major, Minor: minor, Patch: patch, Build: 0},
			ProductVersion: VersionQuad{Major: major, Minor: minor, Patch: patch, Build: 0},
			FileFlagsMask:  "3f",
			FileFlags:      "00",
			FileOS:         "040004",
			FileType:       "01",
			FileSubType:    "00",
		},
		StringFileInfo: StringFileInfo{
			CompanyName:      cfg.Company,
			FileDescription:  cfg.Description,
			FileVersion:      cfg.Version,
			InternalName:     cfg.AppName,
			LegalCopyright:   cfg.Copyright,
			OriginalFilename: cfg.AppName + ".exe",
			ProductName:      cfg.AppName,
			ProductVersion:   cfg.Version,
		},
		VarFileInfo: VarFileInfo{
			Translation: Translation{
				LangID:    "0409",
				CharsetID: "04B0",
			},
		},
		IconPath: cfg.ICOPath,
	}

	outPath := filepath.Join(cfg.OutputDir, "versioninfo.json")
	data, err := json.MarshalIndent(vi, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal versioninfo: %w", err)
	}

	if err := os.WriteFile(outPath, data, 0o644); err != nil {
		return "", fmt.Errorf("write versioninfo.json: %w", err)
	}

	return outPath, nil
}

// GenerateSyso generates a .syso resource file using rsrc.
// The .syso file is automatically linked by the Go linker during build.
func GenerateSyso(icoPath, outputPath, arch string) error {
	if arch == "" {
		arch = "amd64"
	}

	args := []string{
		"-ico", icoPath,
		"-o", outputPath,
		"-arch", arch,
	}

	cmd := exec.Command("rsrc", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsrc failed: %w", err)
	}

	return nil
}

// GenerateSysoGoversioninfo generates a .syso resource file using the
// goversioninfo library (pure Go, no external tool required).
// It builds version info, icon, and manifest directly from the WindowsConfig.
func GenerateSysoGoversioninfo(cfg WindowsConfig, arch string) (string, error) {
	if arch == "" {
		arch = "amd64"
	}

	major, minor, patch := parseVersion(cfg.Version)

	vi := &goversioninfo.VersionInfo{}
	vi.FixedFileInfo.FileVersion.Major = major
	vi.FixedFileInfo.FileVersion.Minor = minor
	vi.FixedFileInfo.FileVersion.Patch = patch
	vi.FixedFileInfo.FileVersion.Build = 0
	vi.FixedFileInfo.ProductVersion.Major = major
	vi.FixedFileInfo.ProductVersion.Minor = minor
	vi.FixedFileInfo.ProductVersion.Patch = patch
	vi.FixedFileInfo.ProductVersion.Build = 0
	vi.FixedFileInfo.FileFlagsMask = "3f"
	vi.FixedFileInfo.FileFlags = "00"
	vi.FixedFileInfo.FileOS = "040004"
	vi.FixedFileInfo.FileType = "01"
	vi.FixedFileInfo.FileSubType = "00"

	vi.StringFileInfo.CompanyName = cfg.Company
	vi.StringFileInfo.FileDescription = cfg.Description
	vi.StringFileInfo.FileVersion = cfg.Version
	vi.StringFileInfo.InternalName = cfg.AppName
	vi.StringFileInfo.LegalCopyright = cfg.Copyright
	vi.StringFileInfo.OriginalFilename = cfg.AppName + ".exe"
	vi.StringFileInfo.ProductName = cfg.AppName
	vi.StringFileInfo.ProductVersion = cfg.Version

	// Set the icon path (relative to where the .syso will be written)
	vi.IconPath = cfg.ICOPath

	// Generate a manifest with DPI awareness and asInvoker
	manifestPath := filepath.Join(cfg.OutputDir, cfg.AppName+".exe.manifest")
	if _, err := os.Stat(manifestPath); err == nil {
		vi.ManifestPath = manifestPath
	}

	vi.Build()
	vi.Walk()

	outputPath := filepath.Join(cfg.OutputDir, "resource_windows_"+arch+".syso")
	if err := vi.WriteSyso(outputPath, arch); err != nil {
		return "", fmt.Errorf("goversioninfo WriteSyso: %w", err)
	}

	return outputPath, nil
}

// GenerateSysoFromJSON generates a .syso resource file from a versioninfo.json
// file using the goversioninfo library (pure Go, no external tool required).
func GenerateSysoFromJSON(versionInfoPath, outputPath, arch string) error {
	if arch == "" {
		arch = "amd64"
	}

	data, err := os.ReadFile(versionInfoPath)
	if err != nil {
		return fmt.Errorf("read versioninfo.json: %w", err)
	}

	vi := &goversioninfo.VersionInfo{}
	if err := vi.ParseJSON(data); err != nil {
		return fmt.Errorf("parse versioninfo.json: %w", err)
	}

	vi.Build()
	vi.Walk()

	if err := vi.WriteSyso(outputPath, arch); err != nil {
		return fmt.Errorf("goversioninfo WriteSyso: %w", err)
	}

	return nil
}

// GenerateSysoWithGovernVersionInfo generates a .syso using the external
// goversioninfo CLI tool. Deprecated: use GenerateSysoGoversioninfo or
// GenerateSysoFromJSON instead, which use the library directly.
func GenerateSysoWithGovernVersionInfo(versionInfoPath, outputPath string) error {
	args := []string{
		"-o", outputPath,
		versionInfoPath,
	}

	cmd := exec.Command("goversioninfo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goversioninfo failed: %w", err)
	}

	return nil
}

// WriteManifest generates a Windows application manifest XML file.
func WriteManifest(appName, outputDir string) (string, error) {
	manifest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="*"
    name="%s"
    type="win32"
  />
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
      <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
      <supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
    </application>
  </compatibility>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true/pm</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">PerMonitorV2</dpiAwareness>
    </windowsSettings>
  </application>
</assembly>`, appName)

	outPath := filepath.Join(outputDir, appName+".exe.manifest")
	if err := os.WriteFile(outPath, []byte(manifest), 0o644); err != nil {
		return "", fmt.Errorf("write manifest: %w", err)
	}

	return outPath, nil
}

func parseVersion(version string) (int, int, int) {
	var major, minor, patch int
	_, _ = fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	// Also try with 'v' prefix
	if major == 0 && minor == 0 && patch == 0 {
		_, _ = fmt.Sscanf(version, "v%d.%d.%d", &major, &minor, &patch)
	}
	return major, minor, patch
}
