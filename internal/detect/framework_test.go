package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFrameworkString(t *testing.T) {
	tests := []struct {
		name      string
		framework Framework
		want      string
	}{
		{name: "None", framework: FrameworkNone, want: "None"},
		{name: "Tauri", framework: FrameworkTauri, want: "Tauri"},
		{name: "Electron", framework: FrameworkElectron, want: "Electron"},
		{name: "Wails", framework: FrameworkWails, want: "Wails"},
		{name: "Fyne", framework: FrameworkFyne, want: "Fyne"},
		{name: "Unknown", framework: Framework(99), want: "None"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.framework.String(); got != tt.want {
				t.Errorf("Framework.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDetectFramework(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, dir string)
		want   Framework
	}{
		{
			name: "Tauri via tauri.conf.json",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "tauri.conf.json"), `{}`)
			},
			want: FrameworkTauri,
		},
		{
			name: "Tauri via src-tauri directory",
			setup: func(t *testing.T, dir string) {
				if err := os.MkdirAll(filepath.Join(dir, "src-tauri"), 0o755); err != nil {
					t.Fatal(err)
				}
			},
			want: FrameworkTauri,
		},
		{
			name: "Wails via wails.json",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "wails.json"), `{}`)
			},
			want: FrameworkWails,
		},
		{
			name: "Wails via build/appicon.png",
			setup: func(t *testing.T, dir string) {
				if err := os.MkdirAll(filepath.Join(dir, "build"), 0o755); err != nil {
					t.Fatal(err)
				}
				writeFile(t, filepath.Join(dir, "build", "appicon.png"), "fake png")
			},
			want: FrameworkWails,
		},
		{
			name: "Fyne via go.mod with fyne.io/fyne",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "go.mod"), `module example.com/app

go 1.21

require fyne.io/fyne/v2 v2.4.0
`)
			},
			want: FrameworkFyne,
		},
		{
			name: "Electron via package.json devDependencies",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "package.json"), `{
  "name": "my-app",
  "devDependencies": {
    "electron": "^28.0.0"
  }
}`)
			},
			want: FrameworkElectron,
		},
		{
			name: "Electron via package.json dependencies",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "package.json"), `{
  "name": "my-app",
  "dependencies": {
    "electron": "^28.0.0"
  }
}`)
			},
			want: FrameworkElectron,
		},
		{
			name: "FrameworkNone when no markers",
			setup: func(t *testing.T, dir string) {
				// empty directory
			},
			want: FrameworkNone,
		},
		{
			name: "FrameworkNone with unrelated package.json",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "package.json"), `{
  "name": "my-app",
  "dependencies": {
    "express": "^4.18.0"
  }
}`)
			},
			want: FrameworkNone,
		},
		{
			name: "FrameworkNone with go.mod without fyne",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "go.mod"), `module example.com/app

go 1.21

require github.com/gin-gonic/gin v1.9.0
`)
			},
			want: FrameworkNone,
		},
		{
			name: "Priority: Tauri before Wails when both present",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "tauri.conf.json"), `{}`)
				writeFile(t, filepath.Join(dir, "wails.json"), `{}`)
			},
			want: FrameworkTauri,
		},
		{
			name: "Priority: Wails before Fyne when both present",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "wails.json"), `{}`)
				writeFile(t, filepath.Join(dir, "go.mod"), `module example.com/app

go 1.21

require fyne.io/fyne/v2 v2.4.0
`)
			},
			want: FrameworkWails,
		},
		{
			name: "Priority: Fyne before Electron when both present",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "go.mod"), `module example.com/app

go 1.21

require fyne.io/fyne/v2 v2.4.0
`)
				writeFile(t, filepath.Join(dir, "package.json"), `{
  "devDependencies": { "electron": "^28.0.0" }
}`)
			},
			want: FrameworkFyne,
		},
		{
			name: "Invalid package.json is not Electron",
			setup: func(t *testing.T, dir string) {
				writeFile(t, filepath.Join(dir, "package.json"), `not valid json`)
			},
			want: FrameworkNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(t, dir)

			got := DetectFramework(dir)
			if got != tt.want {
				t.Errorf("DetectFramework() = %v (%s), want %v (%s)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestTauriIconSizes(t *testing.T) {
	sizes := TauriIconSizes()
	expected := []int{32, 128, 256, 512, 1024}
	if len(sizes) != len(expected) {
		t.Fatalf("TauriIconSizes() returned %d sizes, want %d", len(sizes), len(expected))
	}
	for i, s := range sizes {
		if s != expected[i] {
			t.Errorf("TauriIconSizes()[%d] = %d, want %d", i, s, expected[i])
		}
	}
}

func TestElectronIconSizes(t *testing.T) {
	sizes := ElectronIconSizes()
	expected := []int{16, 32, 48, 64, 128, 256, 512, 1024}
	if len(sizes) != len(expected) {
		t.Fatalf("ElectronIconSizes() returned %d sizes, want %d", len(sizes), len(expected))
	}
	for i, s := range sizes {
		if s != expected[i] {
			t.Errorf("ElectronIconSizes()[%d] = %d, want %d", i, s, expected[i])
		}
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, dir string) string
		want  bool
	}{
		{
			name: "existing file returns true",
			setup: func(t *testing.T, dir string) string {
				p := filepath.Join(dir, "test.txt")
				writeFile(t, p, "hello")
				return p
			},
			want: true,
		},
		{
			name: "non-existent file returns false",
			setup: func(t *testing.T, dir string) string {
				return filepath.Join(dir, "nope.txt")
			},
			want: false,
		},
		{
			name: "directory returns false",
			setup: func(t *testing.T, dir string) string {
				return dir
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := tt.setup(t, dir)
			if got := fileExists(path); got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", path, got, tt.want)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, dir string) string
		want  bool
	}{
		{
			name: "existing directory returns true",
			setup: func(t *testing.T, dir string) string {
				p := filepath.Join(dir, "subdir")
				if err := os.Mkdir(p, 0o755); err != nil {
					t.Fatal(err)
				}
				return p
			},
			want: true,
		},
		{
			name: "non-existent directory returns false",
			setup: func(t *testing.T, dir string) string {
				return filepath.Join(dir, "nope")
			},
			want: false,
		},
		{
			name: "file returns false",
			setup: func(t *testing.T, dir string) string {
				p := filepath.Join(dir, "file.txt")
				writeFile(t, p, "data")
				return p
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := tt.setup(t, dir)
			if got := dirExists(path); got != tt.want {
				t.Errorf("dirExists(%q) = %v, want %v", path, got, tt.want)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		name     string
		haystack string
		needle   string
		want     bool
	}{
		{name: "found at start", haystack: "hello world", needle: "hello", want: true},
		{name: "found at end", haystack: "hello world", needle: "world", want: true},
		{name: "found in middle", haystack: "hello world", needle: "lo wo", want: true},
		{name: "not found", haystack: "hello world", needle: "xyz", want: false},
		{name: "empty haystack", haystack: "", needle: "abc", want: false},
		{name: "empty needle", haystack: "hello", needle: "", want: false},
		{name: "both empty", haystack: "", needle: "", want: false},
		{name: "exact match", haystack: "abc", needle: "abc", want: true},
		{name: "needle longer than haystack", haystack: "ab", needle: "abc", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsString(tt.haystack, tt.needle); got != tt.want {
				t.Errorf("containsString(%q, %q) = %v, want %v", tt.haystack, tt.needle, got, tt.want)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   int
	}{
		{name: "found at start", s: "hello", substr: "hel", want: 0},
		{name: "found at end", s: "hello", substr: "llo", want: 2},
		{name: "found in middle", s: "hello", substr: "ell", want: 1},
		{name: "not found", s: "hello", substr: "xyz", want: -1},
		{name: "empty substr", s: "hello", substr: "", want: 0},
		{name: "substr longer than s", s: "hi", substr: "hello", want: -1},
		{name: "exact match", s: "abc", substr: "abc", want: 0},
		{name: "single char found", s: "abc", substr: "b", want: 1},
		{name: "single char not found", s: "abc", substr: "z", want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexOf(tt.s, tt.substr); got != tt.want {
				t.Errorf("indexOf(%q, %q) = %d, want %d", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
