package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCLI_Help(t *testing.T) {
	cmd := exec.Command("./bin/lin_router", "-help")
	output, _ := cmd.CombinedOutput()

	outputStr := string(output)
	if !strings.Contains(outputStr, "Usage of") {
		t.Error("Expected usage information in help output")
	}
	if !strings.Contains(outputStr, "-i string") {
		t.Error("Expected -i flag in help output")
	}
	if !strings.Contains(outputStr, "-v") {
		t.Error("Expected -v flag in help output")
	}
}

func TestCLI_NoArguments(t *testing.T) {
	cmd := exec.Command("./bin/lin_router")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error when no arguments provided")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Error: -i flag is required") {
		t.Errorf("Expected error message about required -i flag, got: %s", outputStr)
	}
}

func TestCLI_WithInputFile(t *testing.T) {
	cmd := exec.Command("./bin/lin_router", "-i", "main.kicad_pcb")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with valid input file, got error: %v\nOutput: %s", err, string(output))
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "Hello, Lin Router!") {
		t.Error("Should not show greeting without -v flag")
	}
}

func TestCLI_WithVerboseFlag(t *testing.T) {
	cmd := exec.Command("./bin/lin_router", "-i", "main.kicad_pcb", "-v")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("Expected success with -v flag, got error: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Hello, Lin Router!") {
		t.Error("Expected greeting message with -v flag")
	}
}

func TestCLI_InvalidInputFile(t *testing.T) {
	cmd := exec.Command("./bin/lin_router", "-i", "nonexistent.kicad_pcb")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error with non-existent input file")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Error parsing PCB file") {
		t.Errorf("Expected error message about parsing, got: %s", outputStr)
	}
}