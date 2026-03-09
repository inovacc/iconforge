package detect

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Framework represents a detected frontend framework.
type Framework int

const (
	FrameworkNone Framework = iota
	FrameworkTauri
	FrameworkElectron
	FrameworkWails
	FrameworkFyne
)

func (f Framework) String() string {
	switch f {
	case FrameworkTauri:
		return "Tauri"
	case FrameworkElectron:
		return "Electron"
	case FrameworkWails:
		return "Wails"
	case FrameworkFyne:
		return "Fyne"
	default:
		return "None"
	}
}

// DetectFramework scans the project directory for known framework markers.
func DetectFramework(projectDir string) Framework {
	// Tauri: tauri.conf.json or src-tauri/
	if fileExists(filepath.Join(projectDir, "tauri.conf.json")) ||
		dirExists(filepath.Join(projectDir, "src-tauri")) {
		return FrameworkTauri
	}

	// Wails: wails.json or build/appicon.png
	if fileExists(filepath.Join(projectDir, "wails.json")) ||
		fileExists(filepath.Join(projectDir, "build", "appicon.png")) {
		return FrameworkWails
	}

	// Fyne: check go.mod for fyne import
	if hasGoModDep(projectDir, "fyne.io/fyne") {
		return FrameworkFyne
	}

	// Electron: check package.json for electron dependency
	if hasPackageJSONDep(projectDir, "electron") {
		return FrameworkElectron
	}

	return FrameworkNone
}

// TauriIconSizes returns the required icon sizes for Tauri.
func TauriIconSizes() []int {
	return []int{32, 128, 256, 512, 1024}
}

// ElectronIconSizes returns the required icon sizes for Electron.
func ElectronIconSizes() []int {
	return []int{16, 32, 48, 64, 128, 256, 512, 1024}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func hasGoModDep(projectDir, dep string) bool {
	data, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return false
	}
	return containsString(string(data), dep)
}

func hasPackageJSONDep(projectDir, dep string) bool {
	data, err := os.ReadFile(filepath.Join(projectDir, "package.json"))
	if err != nil {
		return false
	}

	var pkg map[string]any
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}

	for _, key := range []string{"dependencies", "devDependencies"} {
		if deps, ok := pkg[key].(map[string]any); ok {
			if _, exists := deps[dep]; exists {
				return true
			}
		}
	}

	return false
}

func containsString(haystack, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 &&
		indexOf(haystack, needle) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
