package test

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"llm-radar/internal/classifier"
	"llm-radar/internal/kb"
	"llm-radar/internal/models"
	"llm-radar/internal/worker"
)

// TestIntegration_WorkerExecuteCommand tests command execution with timeout
func TestIntegration_WorkerExecuteCommand(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		args        []string
		timeout     time.Duration
		expectError bool
		expectCode  int
	}{
		{
			name:        "Simple echo command",
			cmd:         "echo",
			args:        []string{"hello"},
			timeout:     2 * time.Second,
			expectError: false,
			expectCode:  0,
		},
		{
			name:        "Command timeout",
			cmd:         "sleep",
			args:        []string{"10"},
			timeout:     100 * time.Millisecond,
			expectError: false,
			expectCode:  124, // timeout exit code
		},
		{
			name:        "Non-existent command",
			cmd:         "nonexistentcommand12345",
			args:        []string{},
			timeout:     2 * time.Second,
			expectError: true,
			expectCode:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			_, code, err := worker.ExecuteCommandSecure(ctx, tt.cmd, tt.args...)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil && err != context.DeadlineExceeded {
				t.Errorf("Unexpected error: %v", err)
			}

			if code != tt.expectCode {
				t.Errorf("Expected exit code %d, got %d", tt.expectCode, code)
			}
		})
	}
}

// TestIntegration_ClassificationPipeline tests the end-to-end classification
func TestIntegration_ClassificationPipeline(t *testing.T) {
	// Load default KB
	compiledKB, err := kb.LoadAndCompile("")
	if err != nil {
		t.Fatalf("Failed to load KB: %v", err)
	}

	tests := []struct {
		name         string
		modelName    string
		exitCode     int
		output       string
		expectCat    string
	}{
		{
			name:      "Successful response",
			modelName: "test/model",
			exitCode:  0,
			output:    "2, 3, 5",
			expectCat: models.CategoryAvailable,
		},
		{
			name:      "Model not found",
			modelName: "test/nonexistent",
			exitCode:  1,
			output:    "Error: model 'test/nonexistent' not found - 404",
			expectCat: models.CategoryNotFound,
		},
		{
			name:      "Timeout",
			modelName: "test/slow",
			exitCode:  124,
			output:    "",
			expectCat: models.CategoryTimeout,
		},
		{
			name:      "Authentication failed",
			modelName: "test/auth",
			exitCode:  1,
			output:    "Error: unauthorized",
			expectCat: models.CategoryAuthFailed,
		},
		{
			name:      "No quota",
			modelName: "test/quota",
			exitCode:  1,
			output:    "Error: insufficient_quota",
			expectCat: models.CategoryNoQuota,
		},
		{
			name:      "Rate limited",
			modelName: "test/rate",
			exitCode:  1,
			output:    "Error: rate_limit_exceeded",
			expectCat: models.CategoryRateLimited,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := classifier.Classify(tt.modelName, tt.exitCode, tt.output, compiledKB)

			if result.Category != tt.expectCat {
				t.Errorf("Expected category %s, got %s", tt.expectCat, result.Category)
			}

			if result.Icon == "" {
				t.Errorf("Icon should not be empty")
			}
		})
	}
}

// TestIntegration_ModelPrioritization tests model prioritization logic
func TestIntegration_ModelPrioritization(t *testing.T) {
	models := []string{
		"provider1/model-a",
		"google/gemini-flash",
		"zai-coding-plan/test",
		"provider2/model-b",
		"groq/llama3-8b-free",
	}

	freeModels := map[string]bool{
		"google/gemini-flash":  true,
		"groq/llama3-8b-free": true,
	}

	prioritized := worker.PrioritizeModels(models, freeModels, "zai-coding-plan/")

	// Verify free models come first
	if prioritized[0] != "google/gemini-flash" && prioritized[0] != "groq/llama3-8b-free" {
		t.Errorf("Expected free model first, got %s", prioritized[0])
	}

	// Verify zai models come before regular models (but after free)
	foundZai := false
	foundRegular := false
	for _, m := range prioritized {
		if m == "zai-coding-plan/test" {
			foundZai = true
			if foundRegular {
				t.Errorf("ZAI model should come before regular models")
			}
		}
		if m == "provider1/model-a" || m == "provider2/model-b" {
			foundRegular = true
		}
	}

	if !foundZai {
		t.Errorf("ZAI model not found in prioritized list")
	}
}

// TestIntegration_DiscoveryCommand tests model discovery (requires opencode CLI)
func TestIntegration_DiscoveryCommand(t *testing.T) {
	// Skip if opencode is not available
	if _, err := exec.LookPath("opencode"); err != nil {
		t.Skip("Skipping integration test: opencode CLI not found")
	}

	// Create a temporary mock script for opencode
	tmpDir := t.TempDir()
	mockScript := tmpDir + "/opencode"
	
	mockContent := `#!/bin/bash
if [ "$1" = "models" ]; then
  echo "google/gemini-flash"
  echo "anthropic/claude-3-5-sonnet"
  echo "groq/llama3-8b"
  echo "invalid-line-without-slash"
fi
`

	if err := os.WriteFile(mockScript, []byte(mockContent), 0755); err != nil {
		t.Fatalf("Failed to create mock script: %v", err)
	}

	// Temporarily modify PATH
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tmpDir+":"+oldPath)
	defer os.Setenv("PATH", oldPath)

	// Test discovery
	cmd := worker.DiscoverModelsCmd()
	msg := cmd()

	switch result := msg.(type) {
	case []string:
		if len(result) != 3 {
			t.Errorf("Expected 3 valid models, got %d", len(result))
		}

		// Verify invalid lines were filtered out
		for _, model := range result {
			if model == "invalid-line-without-slash" {
				t.Errorf("Invalid model should have been filtered out")
			}
		}
	case error:
		t.Errorf("Discovery returned error: %v", result)
	default:
		t.Errorf("Unexpected message type: %T", msg)
	}
}

// TestIntegration_EndToEnd tests a complete workflow with mock data
func TestIntegration_EndToEnd(t *testing.T) {
	// This test simulates a complete workflow without actually calling opencode
	
	// 1. Load KB
	compiledKB, err := kb.LoadAndCompile("")
	if err != nil {
		t.Fatalf("Failed to load KB: %v", err)
	}

	// 2. Simulate model list
	modelList := []string{
		"google/gemini-flash",
		"anthropic/claude-3-5-sonnet",
		"groq/llama3-8b",
	}

	// 3. Prioritize models
	freeModels := map[string]bool{
		"google/gemini-flash": true,
		"groq/llama3-8b":     true,
	}
	prioritized := worker.PrioritizeModels(modelList, freeModels, "zai-")

	if len(prioritized) != len(modelList) {
		t.Errorf("Prioritization changed model count")
	}

	// 4. Simulate testing models (without actual execution)
	for _, model := range prioritized {
		// Simulate a simple classification
		result := classifier.Classify(model, 0, "2, 3, 5", compiledKB)

		if result.Category == "" {
			t.Errorf("Model %s classification failed", model)
		}

		if result.Icon == "" {
			t.Errorf("Model %s has no icon", model)
		}
	}

	// Test passed if we got here
	t.Logf("End-to-end workflow completed successfully")
}
