// Package models defines the core data structures and constants for the application.
package models

import (
	"time"
)

// ModelResult represents the outcome of availability test for a single model.
type ModelResult struct {
	Model      string `json:"model"`
	Provider   string `json:"provider"`
	Category   string `json:"category"`
	Reason     string `json:"reason"`
	Duration   string `json:"duration"`
	DurationMs int64  `json:"duration_ms"`
	Output     string `json:"output,omitempty"`
	ExitCode   int    `json:"exit_code"`
	Icon       string `json:"icon"`
	Timestamp  string `json:"timestamp"`
}

// ModelInfo contains metadata about a specific model.
type ModelInfo struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Limits      string `json:"limits,omitempty"`
}

// ProviderInfo contains metadata about a specific provider.
type ProviderInfo struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Limits      string `json:"limits"`
}

// RunConfig holds the configuration for a test execution.
type RunConfig struct {
	Prompt      string
	Timeout     time.Duration
	Concurrency int
	Retries     int
	MaxOutputKB int
	UseCache    bool
	CachePath   string
}

// CachedResult wraps a ModelResult with metadata for caching purposes.
type CachedResult struct {
	Result    ModelResult `json:"result"`
	CachedAt  time.Time   `json:"cached_at"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// IsExpired checks if the cached result has exceeded its lifespan.
func (cr *CachedResult) IsExpired() bool {
	return time.Now().After(cr.ExpiresAt)
}

// Model categories constants.
const (
	CategoryFree        = "FREE"
	CategoryFreeLimited = "FREE_LIMITED"
	CategoryPaid        = "PAID"
	CategoryAvailable   = "AVAILABLE"
	CategoryNotFound    = "NOT_FOUND"
	CategoryTimeout     = "TIMEOUT"
	CategoryAuthFailed  = "AUTH_FAILED"
	CategoryNoQuota     = "NO_QUOTA"
	CategoryRateLimited = "RATE_LIMITED"
	CategoryFreeError   = "FREE_ERROR"
	CategoryError       = "ERROR"
)

// CategoryIcons provides a visual representation for each category.
var CategoryIcons = map[string]string{
	CategoryFree:        "🆓",
	CategoryFreeLimited: "📊",
	CategoryPaid:        "💰",
	CategoryAvailable:   "✅",
	CategoryNotFound:    "❓",
	CategoryTimeout:     "⏰",
	CategoryAuthFailed:  "🔒",
	CategoryNoQuota:     "❌",
	CategoryRateLimited: "⏱️",
	CategoryFreeError:   "⚠️",
	CategoryError:       "⚠️",
}

// AllCategories returns a list of all supported categories.
func AllCategories() []string {
	return []string{
		CategoryFree,
		CategoryFreeLimited,
		CategoryPaid,
		CategoryAvailable,
		CategoryNotFound,
		CategoryTimeout,
		CategoryAuthFailed,
		CategoryNoQuota,
		CategoryRateLimited,
		CategoryFreeError,
		CategoryError,
	}
}
