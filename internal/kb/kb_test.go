package kb

import (
	"os"
	"path/filepath"
	"testing"

	"llm-radar/internal/models"
)

func TestDefaultConfigHasFreeModels(t *testing.T) {
	cfg := DefaultConfig()

	if len(cfg.FreeModels) == 0 {
		t.Error("DefaultConfig should have free models")
	}

	// Check specific model
	if _, ok := cfg.FreeModels["opencode/big-pickle"]; !ok {
		t.Error("Expected opencode/big-pickle in free models")
	}
}

func TestDefaultConfigHasFreeTierProviders(t *testing.T) {
	cfg := DefaultConfig()

	expected := []string{"cerebras", "deepseek", "groq"}
	for _, provider := range expected {
		if _, ok := cfg.FreeTierProviders[provider]; !ok {
			t.Errorf("Expected %q in free tier providers", provider)
		}
	}
}

func TestCompileDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	compiled, err := Compile(cfg)

	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if compiled.SuccessRe == nil {
		t.Error("SuccessRe should not be nil")
	}
	if compiled.NotFoundRe == nil {
		t.Error("NotFoundRe should not be nil")
	}
	if compiled.AuthRe == nil {
		t.Error("AuthRe should not be nil")
	}
}

func TestSuccessRegexMatches(t *testing.T) {
	compiled, err := Compile(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"2, 3, 5", true},
		{"2,3,5", true},
		{"2 3 5", true},
		{"The prime numbers are 2, 3, 5", true},
		{"primos", true},
		{"OK", true},
		{"random text", false},
		{"", false},
	}

	for _, tt := range tests {
		result := compiled.SuccessRe.MatchString(tt.input)
		if result != tt.expected {
			t.Errorf("SuccessRe.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestNotFoundRegexMatches(t *testing.T) {
	compiled, err := Compile(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"404 not found", true},
		{"ModelNotFoundError", true},
		{"entity was not found", true},
		{"model works fine", false},
	}

	for _, tt := range tests {
		result := compiled.NotFoundRe.MatchString(tt.input)
		if result != tt.expected {
			t.Errorf("NotFoundRe.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestAuthRegexMatches(t *testing.T) {
	compiled, err := Compile(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"authentication failed", true},
		{"unauthorized", true},
		{"invalid api key", true},
		{"401 Unauthorized", true},
		{"403 Forbidden", true},
		{"model works fine", false},
	}

	for _, tt := range tests {
		result := compiled.AuthRe.MatchString(tt.input)
		if result != tt.expected {
			t.Errorf("AuthRe.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestQuotaRegexMatches(t *testing.T) {
	compiled, err := Compile(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"insufficient quota", true},
		{"quota exceeded", true},
		{"no credits remaining", true},
		{"billing limit reached", true},
		{"model works fine", false},
	}

	for _, tt := range tests {
		result := compiled.QuotaRe.MatchString(tt.input)
		if result != tt.expected {
			t.Errorf("QuotaRe.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestRateLimitRegexMatches(t *testing.T) {
	compiled, err := Compile(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"rate limit exceeded", true},
		{"too many requests", true},
		{"throttled", true},
		{"429 Too Many Requests", true},
		{"model works fine", false},
	}

	for _, tt := range tests {
		result := compiled.RateLimitRe.MatchString(tt.input)
		if result != tt.expected {
			t.Errorf("RateLimitRe.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestTimeoutRegexMatches(t *testing.T) {
	compiled, err := Compile(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"timeout", true},
		{"timed out", true},
		{"deadline exceeded", true},
		{"context deadline exceeded", true},
		{"model works fine", false},
	}

	for _, tt := range tests {
		result := compiled.TimeoutRe.MatchString(tt.input)
		if result != tt.expected {
			t.Errorf("TimeoutRe.MatchString(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestLoadAndCompileFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "kb.json")

	// Write custom KB file
	content := `{
		"free_models": {
			"custom/model": {
				"category": "FREE",
				"description": "Custom Model"
			}
		},
		"success_regex": "(?i)success"
	}`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	compiled, err := LoadAndCompile(path)
	if err != nil {
		t.Fatalf("LoadAndCompile failed: %v", err)
	}

	// Check custom model was loaded
	info, ok := compiled.GetFreeModel("custom/model")
	if !ok {
		t.Error("Expected custom/model in free models")
	}
	if info.Description != "Custom Model" {
		t.Errorf("Description mismatch: got %q", info.Description)
	}
}

func TestLoadAndCompileWithEmptyPath(t *testing.T) {
	compiled, err := LoadAndCompile("")
	if err != nil {
		t.Fatalf("LoadAndCompile with empty path failed: %v", err)
	}

	// Should use defaults
	if len(compiled.Config.FreeModels) == 0 {
		t.Error("Expected default free models")
	}
}

func TestGetFreeModel(t *testing.T) {
	compiled, _ := Compile(DefaultConfig())

	info, ok := compiled.GetFreeModel("opencode/big-pickle")
	if !ok {
		t.Fatal("Expected to find opencode/big-pickle")
	}

	if info.Category != models.CategoryFree {
		t.Errorf("Category mismatch: got %q", info.Category)
	}
}


func TestGetFreeTierProvider(t *testing.T) {
	compiled, _ := Compile(DefaultConfig())

	info, ok := compiled.GetFreeTierProvider("groq")
	if !ok {
		t.Fatal("Expected to find groq provider")
	}

	if info.Category != models.CategoryFreeLimited {
		t.Errorf("Category mismatch: got %q", info.Category)
	}
}

func TestInvalidRegexReturnsError(t *testing.T) {
	cfg := DefaultConfig()
	cfg.SuccessRegex = "[invalid" // Invalid regex

	_, err := Compile(cfg)
	if err == nil {
		t.Error("Expected error for invalid regex")
	}
}
