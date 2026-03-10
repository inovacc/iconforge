package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// resetDetectFlags resets all detect global flag variables to their defaults.
func resetDetectFlags(t *testing.T) {
	t.Helper()
	detectDir = ""
}

// executeDetect sets up the detectCmd output buffers and calls its RunE directly.
func executeDetect(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	detectCmd.SetOut(buf)
	detectCmd.SetErr(buf)
	err := detectCmd.RunE(detectCmd, nil)
	return buf.String(), err
}

func TestDetect_NoFramework(t *testing.T) {
	resetDetectFlags(t)
	tmpDir := t.TempDir()
	detectDir = tmpDir

	output, err := executeDetect(t)
	if err != nil {
		t.Fatalf("detect with empty dir failed: %v", err)
	}

	if !containsStr(output, "No framework detected") {
		t.Errorf("expected 'No framework detected' in output, got: %s", output)
	}
	if !containsStr(output, "Supported frameworks") {
		t.Errorf("expected 'Supported frameworks' in output, got: %s", output)
	}
}

func TestDetect_Tauri(t *testing.T) {
	resetDetectFlags(t)
	tmpDir := t.TempDir()

	// Create tauri marker file
	if err := os.WriteFile(filepath.Join(tmpDir, "tauri.conf.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write tauri.conf.json: %v", err)
	}

	detectDir = tmpDir

	output, err := executeDetect(t)
	if err != nil {
		t.Fatalf("detect tauri failed: %v", err)
	}

	if !containsStr(output, "Detected framework: Tauri") {
		t.Errorf("expected 'Detected framework: Tauri' in output, got: %s", output)
	}
	if !containsStr(output, "icon.ico") {
		t.Errorf("expected 'icon.ico' in required assets, got: %s", output)
	}
	if !containsStr(output, "icon.icns") {
		t.Errorf("expected 'icon.icns' in required assets, got: %s", output)
	}
}

func TestDetect_Electron(t *testing.T) {
	resetDetectFlags(t)
	tmpDir := t.TempDir()

	// Create package.json with electron dependency
	pkgJSON := `{"devDependencies": {"electron": "^28.0.0"}}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0o644); err != nil {
		t.Fatalf("write package.json: %v", err)
	}

	detectDir = tmpDir

	output, err := executeDetect(t)
	if err != nil {
		t.Fatalf("detect electron failed: %v", err)
	}

	if !containsStr(output, "Detected framework: Electron") {
		t.Errorf("expected 'Detected framework: Electron' in output, got: %s", output)
	}
	if !containsStr(output, "icon.png") {
		t.Errorf("expected 'icon.png' in required assets, got: %s", output)
	}
}

func TestDetect_Wails(t *testing.T) {
	resetDetectFlags(t)
	tmpDir := t.TempDir()

	// Create wails marker file
	if err := os.WriteFile(filepath.Join(tmpDir, "wails.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatalf("write wails.json: %v", err)
	}

	detectDir = tmpDir

	output, err := executeDetect(t)
	if err != nil {
		t.Fatalf("detect wails failed: %v", err)
	}

	if !containsStr(output, "Detected framework: Wails") {
		t.Errorf("expected 'Detected framework: Wails' in output, got: %s", output)
	}
	if !containsStr(output, "appicon.png") {
		t.Errorf("expected 'appicon.png' in required assets, got: %s", output)
	}
}

func TestDetect_Fyne(t *testing.T) {
	resetDetectFlags(t)
	tmpDir := t.TempDir()

	// Create go.mod with fyne dependency
	goMod := `module myapp

go 1.21

require fyne.io/fyne/v2 v2.4.0
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	detectDir = tmpDir

	output, err := executeDetect(t)
	if err != nil {
		t.Fatalf("detect fyne failed: %v", err)
	}

	if !containsStr(output, "Detected framework: Fyne") {
		t.Errorf("expected 'Detected framework: Fyne' in output, got: %s", output)
	}
	if !containsStr(output, "Icon.png") {
		t.Errorf("expected 'Icon.png' in required assets, got: %s", output)
	}
}

func TestDetect_DefaultDir(t *testing.T) {
	resetDetectFlags(t)

	// detectDir is empty, so it should use current working directory
	output, err := executeDetect(t)
	if err != nil {
		t.Fatalf("detect with default dir failed: %v", err)
	}

	// Should produce some output (either detected or not detected)
	if output == "" {
		t.Error("expected non-empty output from detect command")
	}
}
