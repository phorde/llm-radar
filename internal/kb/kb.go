// Package kb provides knowledge base management for model classification.
package kb

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"llm-radar/internal/models"
)

// ============================================================================
// CONFIGURATION STRUCTURES
// ============================================================================

// Config holds the knowledge base configuration.
type Config struct {
	FreeModels        map[string]ModelInfo    `json:"free_models"`
	FreeTierProviders map[string]ProviderInfo `json:"free_tier_providers"`
	SuccessRegex      string                  `json:"success_regex"`
	NotFoundRegex     string                  `json:"not_found_regex"`
	AuthRegex         string                  `json:"auth_regex"`
	QuotaRegex        string                  `json:"quota_regex"`
	RateLimitRegex    string                  `json:"rate_limit_regex"`
	TimeoutRegex      string                  `json:"timeout_regex"`
}

// ModelInfo describes a model in the knowledge base.
type ModelInfo struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Limits      string `json:"limits,omitempty"`
}

// ProviderInfo describes a provider in the knowledge base.
type ProviderInfo struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Limits      string `json:"limits"`
}

// Compiled holds both the config and compiled regex patterns.
type Compiled struct {
	Config      Config
	SuccessRe   *regexp.Regexp
	NotFoundRe  *regexp.Regexp
	AuthRe      *regexp.Regexp
	QuotaRe     *regexp.Regexp
	RateLimitRe *regexp.Regexp
	TimeoutRe   *regexp.Regexp
}

// ============================================================================
// DEFAULT CONFIGURATION
// ============================================================================

// DefaultConfig returns the built-in knowledge base configuration.
func DefaultConfig() Config {
	return Config{
		FreeModels: map[string]ModelInfo{
			"opencode/big-pickle": {
				Category:    models.CategoryFree,
				Description: "Zen - Big Pickle",
				Limits:      "limites não documentados",
			},
			"opencode/gpt-5-nano": {
				Category:    models.CategoryFree,
				Description: "Zen - GPT 5 Nano",
				Limits:      "limites não documentados",
			},
			"opencode/minimax-m2.1-free": {
				Category:    models.CategoryFree,
				Description: "Zen - MiniMax M2.1",
				Limits:      "limites não documentados",
			},
			"opencode/glm-4.7-free": {
				Category:    models.CategoryFree,
				Description: "Zen - GLM 4.7",
				Limits:      "limites não documentados",
			},
			"opencode/kimi-k2.5-free": {
				Category:    models.CategoryFree,
				Description: "Zen - Kimi K2.5",
				Limits:      "limites não documentados",
			},
			"opencode/trinity-large-preview-free": {
				Category:    models.CategoryFree,
				Description: "Zen - Trinity Large",
				Limits:      "limites não documentados",
			},
		},

		FreeTierProviders: map[string]ProviderInfo{
			"cerebras": {
				Category:    models.CategoryFreeLimited,
				Description: "Cerebras",
				Limits:      "1M tokens/dia (agregado)",
			},
			"deepseek": {
				Category:    models.CategoryFreeLimited,
				Description: "DeepSeek",
				Limits:      "5M tokens inicial + 50 RPM",
			},
			"groq": {
				Category:    models.CategoryFreeLimited,
				Description: "Groq",
				Limits:      "14.4K req/dia, 30 RPM",
			},
		},

		SuccessRegex:   `(?i)(^|\b)(2\s*,?\s*3\s*,?\s*5|prime|primos|OK)(\b|$)`,
		NotFoundRegex:  `(?i)(404|not\.found|entity.was.not.found|modelnotfounderror)`,
		AuthRegex:      `(?i)(auth|unauthoriz|api\.?key|invalid.*key|401|403)`,
		QuotaRegex:     `(?i)(insufficient.*quota|quota.*exceed|no.*credits?|billing.*limit)`,
		RateLimitRegex: `(?i)(rate.limit|too.many.*request|throttl|429)`,
		TimeoutRegex:   `(?i)(timeout|timed.out|deadline.exceeded)`,
	}
}

// ============================================================================
// LOADING AND COMPILATION
// ============================================================================

// LoadAndCompile loads configuration from file and compiles regex patterns.
// If path is empty or file doesn't exist, uses default configuration.
func LoadAndCompile(path string) (Compiled, error) {
	cfg := DefaultConfig()

	if path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			if err := json.Unmarshal(data, &cfg); err != nil {
				return Compiled{}, fmt.Errorf("erro ao parsear KB: %w", err)
			}
		}
	}

	return Compile(cfg)
}

// Compile compiles the regex patterns from a Config into a Compiled KB.
func Compile(cfg Config) (Compiled, error) {
	ckb := Compiled{Config: cfg}

	var err error
	ckb.SuccessRe, err = regexp.Compile(cfg.SuccessRegex)
	if err != nil {
		return ckb, fmt.Errorf("regex SuccessRegex inválida: %w", err)
	}

	ckb.NotFoundRe, err = regexp.Compile(cfg.NotFoundRegex)
	if err != nil {
		return ckb, fmt.Errorf("regex NotFoundRegex inválida: %w", err)
	}

	ckb.AuthRe, err = regexp.Compile(cfg.AuthRegex)
	if err != nil {
		return ckb, fmt.Errorf("regex AuthRegex inválida: %w", err)
	}

	ckb.QuotaRe, err = regexp.Compile(cfg.QuotaRegex)
	if err != nil {
		return ckb, fmt.Errorf("regex QuotaRegex inválida: %w", err)
	}

	ckb.RateLimitRe, err = regexp.Compile(cfg.RateLimitRegex)
	if err != nil {
		return ckb, fmt.Errorf("regex RateLimitRegex inválida: %w", err)
	}

	ckb.TimeoutRe, err = regexp.Compile(cfg.TimeoutRegex)
	if err != nil {
		return ckb, fmt.Errorf("regex TimeoutRegex inválida: %w", err)
	}

	return ckb, nil
}

// ============================================================================
// LOOKUP METHODS
// ============================================================================

// GetFreeModel returns model info if the model is in the free models list.
func (c *Compiled) GetFreeModel(model string) (ModelInfo, bool) {
	info, ok := c.Config.FreeModels[model]
	return info, ok
}

// GetFreeTierProvider returns provider info if the provider has a free tier.
func (c *Compiled) GetFreeTierProvider(provider string) (ProviderInfo, bool) {
	info, ok := c.Config.FreeTierProviders[provider]
	return info, ok
}
