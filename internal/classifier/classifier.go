// Package classifier provides logic for categorizing model test results.
package classifier

import (
	"strings"

	"llm-radar/internal/kb"
	"llm-radar/internal/models"
)

// Result holds the classification details.
type Result struct {
	Category string
	Reason   string
	Icon     string
}

// Classify determines the category, reason, and icon for a model's test result.
func Classify(model string, exitCode int, output string, compiledKB kb.Compiled) Result {
	// Check for not found first (takes priority)
	if compiledKB.NotFoundRe.MatchString(output) {
		return Result{
			Category: models.CategoryNotFound,
			Reason:   "Modelo não disponível no OpenCode",
			Icon:     models.CategoryIcons[models.CategoryNotFound],
		}
	}

	// Check for timeout
	if exitCode == 124 || compiledKB.TimeoutRe.MatchString(output) {
		return Result{
			Category: models.CategoryTimeout,
			Reason:   "Timeout (20s)",
			Icon:     models.CategoryIcons[models.CategoryTimeout],
		}
	}

	// Check Knowledge Base for known free models
	if info, exists := compiledKB.Config.FreeModels[model]; exists {
		if exitCode == 0 && compiledKB.SuccessRe.MatchString(output) {
			return Result{
				Category: info.Category,
				Reason:   info.Description,
				Icon:     models.CategoryIcons[info.Category],
			}
		}
		return Result{
			Category: models.CategoryFreeError,
			Reason:   info.Description + " (falhou no teste)",
			Icon:     models.CategoryIcons[models.CategoryFreeError],
		}
	}

	// Check for -free suffix
	if strings.HasSuffix(model, "-free") {
		if exitCode == 0 && compiledKB.SuccessRe.MatchString(output) {
			return Result{
				Category: models.CategoryFree,
				Reason:   "Sufixo -free detectado",
				Icon:     models.CategoryIcons[models.CategoryFree],
			}
		}
		return Result{
			Category: models.CategoryFreeError,
			Reason:   "Sufixo -free (falhou)",
			Icon:     models.CategoryIcons[models.CategoryFreeError],
		}
	}

	// Check for successful generic model
	if exitCode == 0 && compiledKB.SuccessRe.MatchString(output) {
		provider := strings.Split(model, "/")[0]
		// Check if provider has free tier
		if info, exists := compiledKB.Config.FreeTierProviders[provider]; exists {
			return Result{
				Category: info.Category,
				Reason:   info.Limits,
				Icon:     models.CategoryIcons[info.Category],
			}
		}
		// Fallback to generic available
		return Result{
			Category: models.CategoryAvailable,
			Reason:   "Modelo disponível",
			Icon:     models.CategoryIcons[models.CategoryAvailable],
		}
	}

	// Check for specific error patterns
	if compiledKB.QuotaRe.MatchString(output) {
		return Result{
			Category: models.CategoryNoQuota,
			Reason:   "Sem créditos",
			Icon:     models.CategoryIcons[models.CategoryNoQuota],
		}
	}
	if compiledKB.AuthRe.MatchString(output) {
		return Result{
			Category: models.CategoryAuthFailed,
			Reason:   "API key inválida",
			Icon:     models.CategoryIcons[models.CategoryAuthFailed],
		}
	}
	if compiledKB.RateLimitRe.MatchString(output) {
		return Result{
			Category: models.CategoryRateLimited,
			Reason:   "Rate limit",
			Icon:     models.CategoryIcons[models.CategoryRateLimited],
		}
	}

	// Default error
	return Result{
		Category: models.CategoryError,
		Reason:   "Erro desconhecido",
		Icon:     models.CategoryIcons[models.CategoryError],
	}
}
