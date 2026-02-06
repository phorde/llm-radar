package classifier

import (
	"testing"

	"llm-radar/internal/kb"
	"llm-radar/internal/models"
)

func getTestKB(t *testing.T) kb.Compiled {
	compiled, err := kb.Compile(kb.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to compile KB: %v", err)
	}
	return compiled
}

func TestClassifyNotFound(t *testing.T) {
	compiled := getTestKB(t)

	tests := []struct {
		output string
	}{
		{"404 not found"},
		{"ModelNotFoundError: model does not exist"},
		{"entity was not found in registry"},
	}

	for _, tt := range tests {
		result := Classify("any/model", 1, tt.output, compiled)
		if result.Category != models.CategoryNotFound {
			t.Errorf("Expected NOT_FOUND for %q, got %s", tt.output, result.Category)
		}
	}
}

func TestClassifyTimeout(t *testing.T) {
	compiled := getTestKB(t)

	// Exit code 124 indicates timeout
	result := Classify("any/model", 124, "", compiled)
	if result.Category != models.CategoryTimeout {
		t.Errorf("Expected TIMEOUT for exit code 124, got %s", result.Category)
	}

	// Timeout in output
	result = Classify("any/model", 1, "connection timed out", compiled)
	if result.Category != models.CategoryTimeout {
		t.Errorf("Expected TIMEOUT for 'timed out' output, got %s", result.Category)
	}
}

func TestClassifyFreeModel(t *testing.T) {
	compiled := getTestKB(t)

	// Successful free model
	result := Classify("opencode/big-pickle", 0, "2, 3, 5", compiled)
	if result.Category != models.CategoryFree {
		t.Errorf("Expected FREE, got %s", result.Category)
	}

	// Failed free model
	result = Classify("opencode/big-pickle", 1, "error occurred", compiled)
	if result.Category != models.CategoryFreeError {
		t.Errorf("Expected FREE_ERROR, got %s", result.Category)
	}
}

func TestClassifyFreeSuffix(t *testing.T) {
	compiled := getTestKB(t)

	// Model with -free suffix that succeeds
	result := Classify("unknown/model-free", 0, "primos", compiled)
	if result.Category != models.CategoryFree {
		t.Errorf("Expected FREE for -free suffix, got %s", result.Category)
	}

	// Model with -free suffix that fails
	result = Classify("unknown/model-free", 1, "error", compiled)
	if result.Category != models.CategoryFreeError {
		t.Errorf("Expected FREE_ERROR for failed -free model, got %s", result.Category)
	}
}


func TestClassifyFreeTierProvider(t *testing.T) {
	compiled := getTestKB(t)

	// Successful model from groq (free tier provider)
	result := Classify("groq/llama-3.1", 0, "2, 3, 5", compiled)
	if result.Category != models.CategoryFreeLimited {
		t.Errorf("Expected FREE_LIMITED for groq, got %s", result.Category)
	}

	// Successful model from cerebras (free tier provider)
	result = Classify("cerebras/llama-3.1", 0, "primos", compiled)
	if result.Category != models.CategoryFreeLimited {
		t.Errorf("Expected FREE_LIMITED for cerebras, got %s", result.Category)
	}
}

func TestClassifyAvailable(t *testing.T) {
	compiled := getTestKB(t)

	// Generic successful model
	result := Classify("unknown/model", 0, "2, 3, 5", compiled)
	if result.Category != models.CategoryAvailable {
		t.Errorf("Expected AVAILABLE, got %s", result.Category)
	}
}

func TestClassifyAuthFailed(t *testing.T) {
	compiled := getTestKB(t)

	tests := []struct {
		output string
	}{
		{"authentication failed"},
		{"unauthorized access"},
		{"invalid api key"},
		{"401 Unauthorized"},
	}

	for _, tt := range tests {
		result := Classify("unknown/model", 1, tt.output, compiled)
		if result.Category != models.CategoryAuthFailed {
			t.Errorf("Expected AUTH_FAILED for %q, got %s", tt.output, result.Category)
		}
	}
}

func TestClassifyNoQuota(t *testing.T) {
	compiled := getTestKB(t)

	tests := []struct {
		output string
	}{
		{"insufficient quota"},
		{"quota exceeded"},
		{"no credits remaining"},
	}

	for _, tt := range tests {
		result := Classify("unknown/model", 1, tt.output, compiled)
		if result.Category != models.CategoryNoQuota {
			t.Errorf("Expected NO_QUOTA for %q, got %s", tt.output, result.Category)
		}
	}
}

func TestClassifyRateLimited(t *testing.T) {
	compiled := getTestKB(t)

	tests := []struct {
		output string
	}{
		{"rate limit exceeded"},
		{"too many requests"},
		{"429 Too Many Requests"},
	}

	for _, tt := range tests {
		result := Classify("unknown/model", 1, tt.output, compiled)
		if result.Category != models.CategoryRateLimited {
			t.Errorf("Expected RATE_LIMITED for %q, got %s", tt.output, result.Category)
		}
	}
}

func TestClassifyError(t *testing.T) {
	compiled := getTestKB(t)

	// Unknown error
	result := Classify("unknown/model", 1, "something completely random", compiled)
	if result.Category != models.CategoryError {
		t.Errorf("Expected ERROR, got %s", result.Category)
	}
}

func TestClassifyResultHasIcon(t *testing.T) {
	compiled := getTestKB(t)

	result := Classify("opencode/big-pickle", 0, "2, 3, 5", compiled)
	if result.Icon == "" {
		t.Error("Expected non-empty icon")
	}
}

func TestClassifyPriorityNotFoundOverTimeout(t *testing.T) {
	compiled := getTestKB(t)

	// Not found should take priority even with timeout output
	result := Classify("any/model", 124, "404 not found timeout", compiled)
	if result.Category != models.CategoryNotFound {
		t.Errorf("Expected NOT_FOUND over TIMEOUT, got %s", result.Category)
	}
}
