package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
)

func resetCmdtreeFlags(t *testing.T) {
	t.Helper()
	cmdtreeVerbose = true
	cmdtreeBrief = false
	cmdtreeCommand = ""
	cmdtreeJSON = false
}

func executeCmdtree(t *testing.T) (string, error) {
	t.Helper()
	buf := new(bytes.Buffer)
	cmdtreeCmd.SetOut(buf)
	cmdtreeCmd.SetErr(buf)
	err := cmdtreeCmd.RunE(cmdtreeCmd, nil)
	return buf.String(), err
}

func TestCmdtree_Default(t *testing.T) {
	resetCmdtreeFlags(t)

	output, err := executeCmdtree(t)
	if err != nil {
		t.Fatalf("cmdtree default failed: %v", err)
	}

	if !containsStr(output, "iconforge") {
		t.Errorf("expected 'iconforge' in output, got: %s", output)
	}
	if !containsStr(output, "forge") {
		t.Errorf("expected 'forge' in output, got: %s", output)
	}
}

func TestCmdtree_Brief(t *testing.T) {
	resetCmdtreeFlags(t)
	cmdtreeBrief = true

	output, err := executeCmdtree(t)
	if err != nil {
		t.Fatalf("cmdtree --brief failed: %v", err)
	}

	if !containsStr(output, "forge") {
		t.Errorf("expected 'forge' in output, got: %s", output)
	}
}

func TestCmdtree_JSON(t *testing.T) {
	resetCmdtreeFlags(t)
	cmdtreeJSON = true

	output, err := executeCmdtree(t)
	if err != nil {
		t.Fatalf("cmdtree --json failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("cmdtree --json returned invalid JSON: %v\noutput: %s", err, output)
	}

	if _, ok := result["name"]; !ok {
		t.Error("expected 'name' field in JSON output")
	}
}

func TestCmdtree_SingleCommand(t *testing.T) {
	resetCmdtreeFlags(t)
	cmdtreeCommand = "forge"

	output, err := executeCmdtree(t)
	if err != nil {
		t.Fatalf("cmdtree --command forge failed: %v", err)
	}

	if !containsStr(output, "forge") {
		t.Errorf("expected 'forge' in output, got: %s", output)
	}
}

func TestCmdtree_SingleCommandNotFound(t *testing.T) {
	resetCmdtreeFlags(t)
	cmdtreeCommand = "nonexistent"

	_, err := executeCmdtree(t)
	if err == nil {
		t.Fatal("expected error for nonexistent command, got nil")
	}
	if !containsStr(err.Error(), "command not found") {
		t.Errorf("expected 'command not found' error, got: %v", err)
	}
}
