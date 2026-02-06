package worker

import (
	"strings"
	"testing"
)

func TestPrioritizeModels(t *testing.T) {
	freeModels := map[string]bool{
		"opencode/free-model": true,
		"opencode/big-pickle": true,
	}

	models := []string{
		"anthropic/claude",
		"opencode/free-model",
		"zai-coding-plan/glm-4.5",
		"groq/llama-3.1",
		"opencode/big-pickle",
		"zai-coding-plan/glm-4.7",
	}

	result := PrioritizeModels(models, freeModels, "zai-coding-plan/")

	// Free models should be first
	if result[0] != "opencode/free-model" && result[0] != "opencode/big-pickle" {
		t.Errorf("Expected free model first, got %s", result[0])
	}

	// Check ordering: free, then zai, then others
	freeCount := 0
	zaiCount := 0
	otherStart := -1

	for i, m := range result {
		if freeModels[m] {
			freeCount++
		} else if strings.HasPrefix(m, "zai-coding-plan/") {
			zaiCount++
			if otherStart != -1 && otherStart < i {
				t.Error("ZAI model appeared after other model")
			}
		} else {
			if otherStart == -1 {
				otherStart = i
			}
		}
	}

	if freeCount != 2 {
		t.Errorf("Expected 2 free models, got %d", freeCount)
	}
	if zaiCount != 2 {
		t.Errorf("Expected 2 ZAI models, got %d", zaiCount)
	}
}

func TestSmartTrimShortString(t *testing.T) {
	short := "This is a short string"
	result := SmartTrim(short, 64)

	if result != short {
		t.Error("Short strings should not be trimmed")
	}
}

func TestSmartTrimLongString(t *testing.T) {
	// Create a string longer than 1KB
	long := strings.Repeat("x", 2048)
	result := SmartTrim(long, 1) // 1KB max

	if len(result) >= len(long) {
		t.Error("Long string should be trimmed")
	}

	if !strings.Contains(result, "[TRUNCATED]") {
		t.Error("Trimmed string should contain truncation marker")
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exact", 5, "exact"},
	}

	for _, tt := range tests {
		result := Truncate(tt.input, tt.max)
		if result != tt.expected {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.max, result, tt.expected)
		}
	}
}

func TestExtractProvider(t *testing.T) {
	tests := []struct {
		model    string
		expected string
	}{
		{"opencode/big-pickle", "opencode"},
		{"groq/llama-3.1-70b", "groq"},
		{"zai-coding-plan/glm-4.5", "zai-coding-plan"},
		{"single", "single"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ExtractProvider(tt.model)
		if result != tt.expected {
			t.Errorf("ExtractProvider(%q) = %q, want %q", tt.model, result, tt.expected)
		}
	}
}

func TestPrioritizeModelsEmptyInput(t *testing.T) {
	result := PrioritizeModels([]string{}, map[string]bool{}, "zai-")

	if len(result) != 0 {
		t.Error("Empty input should produce empty output")
	}
}

func TestPrioritizeModelsNoFreeModels(t *testing.T) {
	models := []string{"anthropic/claude", "openai/gpt-4"}
	result := PrioritizeModels(models, map[string]bool{}, "zai-coding-plan/")

	if len(result) != 2 {
		t.Error("Should preserve all models")
	}
}
