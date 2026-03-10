package cmd

import (
	"encoding/json"
	"testing"
)

func TestGetVersionInfo(t *testing.T) {
	info := GetVersionInfo()

	if info == nil {
		t.Fatal("GetVersionInfo returned nil")
	}

	if info.Version == "" {
		t.Error("Version should not be empty")
	}
	if info.GitHash == "" {
		t.Error("GitHash should not be empty")
	}
	if info.BuildTime == "" {
		t.Error("BuildTime should not be empty")
	}
	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
}

func TestGetVersionJSON(t *testing.T) {
	jsonStr := GetVersionJSON()

	if jsonStr == "" {
		t.Fatal("GetVersionJSON returned empty string")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("GetVersionJSON returned invalid JSON: %v", err)
	}

	// Verify expected fields exist
	expectedFields := []string{"version", "git_hash", "build_time", "build_hash", "go_version", "goos", "goarch"}
	for _, field := range expectedFields {
		if _, ok := result[field]; !ok {
			t.Errorf("expected field %q in version JSON", field)
		}
	}
}

func TestGetVersionInfo_DefaultValues(t *testing.T) {
	info := GetVersionInfo()

	// The default values set in the var block
	if info.Version != "dev" {
		t.Errorf("expected default Version 'dev', got %q", info.Version)
	}
	if info.GitHash != "none" {
		t.Errorf("expected default GitHash 'none', got %q", info.GitHash)
	}
	if info.BuildHash != "none" {
		t.Errorf("expected default BuildHash 'none', got %q", info.BuildHash)
	}
}
