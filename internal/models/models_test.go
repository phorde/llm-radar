package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCategoryIconsComplete(t *testing.T) {
	categories := AllCategories()

	for _, cat := range categories {
		if _, ok := CategoryIcons[cat]; !ok {
			t.Errorf("Category %q is missing an icon", cat)
		}
	}
}

func TestAllCategoriesCount(t *testing.T) {
	categories := AllCategories()
	expected := 11

	if len(categories) != expected {
		t.Errorf("Expected %d categories, got %d", expected, len(categories))
	}
}

func TestModelResultJSON(t *testing.T) {
	result := ModelResult{
		Model:      "test/model",
		Provider:   "test",
		Category:   CategoryFree,
		Reason:     "Test reason",
		Duration:   "1.5s",
		DurationMs: 1500,
		Output:     "2, 3, 5",
		ExitCode:   0,
		Icon:       CategoryIcons[CategoryFree],
		Timestamp:  "2026-02-06T00:00:00Z",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ModelResult: %v", err)
	}

	var unmarshaled ModelResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ModelResult: %v", err)
	}

	if unmarshaled.Model != result.Model {
		t.Errorf("Model mismatch: got %q, want %q", unmarshaled.Model, result.Model)
	}

	if unmarshaled.Category != result.Category {
		t.Errorf("Category mismatch: got %q, want %q", unmarshaled.Category, result.Category)
	}
}

func TestCachedResultExpiration(t *testing.T) {
	// Create an already expired result
	expired := CachedResult{
		Result:    ModelResult{Model: "test/expired"},
		CachedAt:  time.Now().Add(-2 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if !expired.IsExpired() {
		t.Error("Expected expired result to be marked as expired")
	}

	// Create a valid (not expired) result
	valid := CachedResult{
		Result:    ModelResult{Model: "test/valid"},
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if valid.IsExpired() {
		t.Error("Expected valid result to NOT be marked as expired")
	}
}

func TestRunConfigDefaults(t *testing.T) {
	cfg := RunConfig{
		Prompt:      "test prompt",
		Timeout:     20 * time.Second,
		Concurrency: 5,
		Retries:     1,
		MaxOutputKB: 64,
		UseCache:    false,
		CachePath:   "/tmp/cache.json",
	}

	if cfg.Concurrency < 1 {
		t.Error("Concurrency should be at least 1")
	}

	if cfg.Timeout < time.Second {
		t.Error("Timeout should be at least 1 second")
	}
}
