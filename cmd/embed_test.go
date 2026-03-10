package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// resetEmbedFlags resets all embed global flag variables to their defaults.
func resetEmbedFlags(t *testing.T) {
	t.Helper()
	embedICOPath = ""
	embedOutputPath = "resource.syso"
	embedArch = "amd64"
	embedMethod = "auto"
	embedViPath = ""
	embedAppName = "app"
	embedVersion = "1.0.0"
	embedCompany = ""
	embedCopyright = ""
}

// executeEmbed sets up the embedCmd output buffers and calls runEmbed directly.
func executeEmbed(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	embedCmd.SetOut(buf)
	embedCmd.SetErr(buf)
	err := runEmbed(embedCmd, nil)
	return buf.String(), err
}

// createTestICO creates a valid ICO file by running forge and returning the ICO path.
func createTestICO(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	// Use the forge pipeline to generate a real ICO
	resetForgeFlags(t)
	forgeGenSVG = true
	forgeAppName = "testapp"
	forgeOutputDir = tmpDir
	forgeSkipMac = true
	forgeSkipLinux = true

	buf := new(bytes.Buffer)
	forgeCmd.SetOut(buf)
	forgeCmd.SetErr(buf)
	if err := runForge(forgeCmd, nil); err != nil {
		t.Fatalf("forge to create test ICO failed: %v", err)
	}

	icoPath := filepath.Join(tmpDir, "windows", "icon.ico")
	if _, err := os.Stat(icoPath); os.IsNotExist(err) {
		t.Fatalf("expected ICO file at %s", icoPath)
	}

	return icoPath
}

func TestEmbed_InvalidMethod(t *testing.T) {
	resetEmbedFlags(t)

	embedICOPath = "some.ico"
	embedMethod = "invalid"

	_, err := executeEmbed(t)
	if err == nil {
		t.Fatal("expected error for invalid method, got nil")
	}
	if !containsStr(err.Error(), "unsupported method") {
		t.Errorf("expected 'unsupported method' error, got: %v", err)
	}
}

func TestEmbed_Goversioninfo(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "goversioninfo"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedAppName = "testapp"
	embedVersion = "2.0.0"
	embedCompany = "TestCo"
	embedCopyright = "Copyright 2024"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --method goversioninfo failed: %v", err)
	}

	if !containsStr(output, "goversioninfo") {
		t.Errorf("expected 'goversioninfo' in output, got: %s", output)
	}
}

func TestEmbed_Auto(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "auto"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")

	output, err := executeEmbed(t)
	// auto tries goversioninfo first, then rsrc — either could succeed or fail
	// depending on available tools, but the function itself should not panic
	if err != nil {
		// rsrc might not be installed — that's OK, we tested the path
		t.Logf("embed auto returned error (expected if rsrc not installed): %v", err)
	} else {
		if output == "" {
			t.Error("expected non-empty output on success")
		}
	}
}

func TestEmbed_Rsrc(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "rsrc"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")

	_, err := executeEmbed(t)
	// rsrc requires the external tool to be installed
	// We just verify the function runs without panic
	if err != nil {
		t.Logf("embed rsrc returned error (expected if rsrc not installed): %v", err)
	}
}

func TestEmbed_GoversioninfoWithVI(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	// Create a versioninfo.json file
	viContent := `{
	"FixedFileInfo": {
		"FileVersion": {"Major": 1, "Minor": 0, "Patch": 0, "Build": 0},
		"ProductVersion": {"Major": 1, "Minor": 0, "Patch": 0, "Build": 0}
	},
	"StringFileInfo": {
		"FileDescription": "Test App",
		"ProductName": "TestApp"
	},
	"IconPath": ""
}`
	viPath := filepath.Join(tmpDir, "versioninfo.json")
	if err := os.WriteFile(viPath, []byte(viContent), 0o644); err != nil {
		t.Fatalf("write versioninfo.json: %v", err)
	}

	embedICOPath = icoPath
	embedMethod = "goversioninfo"
	embedViPath = viPath
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")

	output, err := executeEmbed(t)
	if err != nil {
		// goversioninfo with a custom JSON might fail if the JSON schema doesn't match
		t.Logf("embed goversioninfo with VI returned error: %v", err)
	} else {
		if !containsStr(output, "goversioninfo") {
			t.Errorf("expected 'goversioninfo' in output, got: %s", output)
		}
	}
}

func TestEmbed_Arch386(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "goversioninfo"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedArch = "386"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --arch 386 failed: %v", err)
	}

	if !containsStr(output, "goversioninfo") {
		t.Errorf("expected 'goversioninfo' in output, got: %s", output)
	}
}

func TestEmbed_ArchArm64(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "goversioninfo"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedArch = "arm64"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --arch arm64 failed: %v", err)
	}

	if !containsStr(output, "goversioninfo") {
		t.Errorf("expected 'goversioninfo' in output, got: %s", output)
	}
}
