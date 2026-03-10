package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func resetAIContextFlags(t *testing.T) {
	t.Helper()
	aicontextJSON = false
	aicontextCompact = false
}

func executeAIContext(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	aicontextCmd.SetOut(buf)
	aicontextCmd.SetErr(buf)
	err := aicontextCmd.RunE(aicontextCmd, nil)
	return buf.String(), err
}

func TestAIContext_Markdown(t *testing.T) {
	resetAIContextFlags(t)

	output, err := executeAIContext(t)
	if err != nil {
		t.Fatalf("aicontext markdown failed: %v", err)
	}

	if !containsStr(output, "iconforge") {
		t.Errorf("expected 'iconforge' in output, got: %s", output)
	}
	if !containsStr(output, "Overview") {
		t.Errorf("expected 'Overview' in output, got: %s", output)
	}
	if !containsStr(output, "Commands") {
		t.Errorf("expected 'Commands' in output, got: %s", output)
	}
}

func TestAIContext_Compact(t *testing.T) {
	resetAIContextFlags(t)
	aicontextCompact = true

	output, err := executeAIContext(t)
	if err != nil {
		t.Fatalf("aicontext --compact failed: %v", err)
	}

	if !containsStr(output, "iconforge") {
		t.Errorf("expected 'iconforge' in output, got: %s", output)
	}
}

func TestAIContext_JSON(t *testing.T) {
	resetAIContextFlags(t)
	aicontextJSON = true

	output, err := executeAIContext(t)
	if err != nil {
		t.Fatalf("aicontext --json failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("aicontext --json returned invalid JSON: %v\noutput: %s", err, output)
	}

	if _, ok := result["tool"]; !ok {
		t.Error("expected 'tool' field in JSON output")
	}
	if _, ok := result["commands"]; !ok {
		t.Error("expected 'commands' field in JSON output")
	}
}
