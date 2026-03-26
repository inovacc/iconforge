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

func TestEmbed_Winres(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "auto"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedAppName = "testapp"
	embedVersion = "2.0.0"
	embedCompany = "TestCo"
	embedCopyright = "Copyright 2024"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --method auto failed: %v", err)
	}

	if !containsStr(output, "winres") {
		t.Errorf("expected 'winres' in output, got: %s", output)
	}
}

func TestEmbed_GoversioninfoBackCompat(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	// "goversioninfo" method should map to winres
	embedICOPath = icoPath
	embedMethod = "goversioninfo"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedAppName = "testapp"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --method goversioninfo (back compat) failed: %v", err)
	}

	if !containsStr(output, "winres") {
		t.Errorf("expected 'winres' in output, got: %s", output)
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
	if err != nil {
		t.Fatalf("embed auto failed: %v", err)
	}
	if output == "" {
		t.Error("expected non-empty output on success")
	}
}

func TestEmbed_Simple(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "simple"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed simple failed: %v", err)
	}
	if !containsStr(output, "winres") {
		t.Errorf("expected 'winres' in output, got: %s", output)
	}
}

func TestEmbed_RsrcBackCompat(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	// "rsrc" method should map to simple
	embedICOPath = icoPath
	embedMethod = "rsrc"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --method rsrc (back compat) failed: %v", err)
	}
	if !containsStr(output, "winres") {
		t.Errorf("expected 'winres' in output, got: %s", output)
	}
}

func TestEmbed_Arch386(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "auto"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedArch = "386"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --arch 386 failed: %v", err)
	}

	if !containsStr(output, "winres") {
		t.Errorf("expected 'winres' in output, got: %s", output)
	}
}

func TestEmbed_ArchArm64(t *testing.T) {
	resetEmbedFlags(t)
	tmpDir := t.TempDir()
	icoPath := createTestICO(t)

	embedICOPath = icoPath
	embedMethod = "auto"
	embedOutputPath = filepath.Join(tmpDir, "resource.syso")
	embedArch = "arm64"

	output, err := executeEmbed(t)
	if err != nil {
		t.Fatalf("embed --arch arm64 failed: %v", err)
	}

	if !containsStr(output, "winres") {
		t.Errorf("expected 'winres' in output, got: %s", output)
	}
}
