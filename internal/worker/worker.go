// Package worker provides helpers for CLI execution and result processing.
package worker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"

	"llm-radar/internal/cache"
	"llm-radar/internal/classifier"
	"llm-radar/internal/kb"
	"llm-radar/internal/models"
)

// ============================================================================
// WORKER ORCHESTRATION
// ============================================================================

// StartWorkers spawns concurrent workers to test models.
// The msgChan receives generic tea.Msg values that should be understood by the TUI layer.
func StartWorkers(
	modelList []string,
	cfg models.RunConfig,
	compiledKB kb.Compiled,
	resCache *cache.ResultCache,
	msgChan chan tea.Msg,
	processed *int32,
) {
	jobs := make(chan string, len(modelList))
	var wg sync.WaitGroup

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		initialDelay := time.Duration(i*200) * time.Millisecond

		go func(delay time.Duration) {
			defer wg.Done()
			time.Sleep(delay)

			for model := range jobs {
				// Send worker start notification as a generic message
				// The TUI layer will handle the actual message type
				startMsg := struct {
					Model string
					Start time.Time
				}{model, time.Now()}
				msgChan <- startMsg

				var res models.ModelResult
				if cfg.UseCache {
					if cached, ok := resCache.Get(model); ok {
						res = cached
						res.Reason += " (cached)"
					}
				}

				if res.Model == "" {
					res = TestModel(model, cfg, compiledKB)
					if cfg.UseCache {
						resCache.Set(model, res)
					}
				}

				atomic.AddInt32(processed, 1)
				msgChan <- res
			}
		}(initialDelay)
	}

	for _, m := range modelList {
		jobs <- m
	}
	close(jobs)

	wg.Wait()
	close(msgChan)
}

// TestModel tests a single model and returns the result.
func TestModel(modelName string, cfg models.RunConfig, compiledKB kb.Compiled) models.ModelResult {
	provider := ExtractProvider(modelName)

	var lastOut string
	var exitCode int
	var duration time.Duration

	for attempt := 0; attempt <= cfg.Retries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		start := time.Now()

		args := []string{"run", "--model", modelName}
		args = append(args, cfg.Prompt)

		out, code, err := ExecuteCommandSecure(ctx, "opencode", args...)
		duration = time.Since(start)
		cancel()

		lastOut = out
		exitCode = code

		if err == context.DeadlineExceeded {
			exitCode = 124
		}

		if exitCode == 0 && compiledKB.SuccessRe.MatchString(lastOut) {
			break
		}

		if compiledKB.NotFoundRe.MatchString(lastOut) ||
			compiledKB.AuthRe.MatchString(lastOut) ||
			compiledKB.QuotaRe.MatchString(lastOut) {
			break
		}

		if compiledKB.RateLimitRe.MatchString(lastOut) && attempt < cfg.Retries {
			backoff := time.Duration((attempt+1)*500) * time.Millisecond
			time.Sleep(backoff)
			continue
		}

		if exitCode == 124 {
			break
		}
	}

	outTrimmed := SmartTrim(lastOut, cfg.MaxOutputKB)
	result := classifier.Classify(modelName, exitCode, outTrimmed, compiledKB)

	return models.ModelResult{
		Model:      modelName,
		Provider:   provider,
		Category:   result.Category,
		Reason:     result.Reason,
		Icon:       result.Icon,
		Duration:   duration.Round(time.Millisecond).String(),
		DurationMs: duration.Milliseconds(),
		Output:     outTrimmed,
		ExitCode:   exitCode,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
}

// ExecuteCommandSecure executes a command with timeout and proper cleanup.
func ExecuteCommandSecure(ctx context.Context, name string, args ...string) (string, int, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return "", 1, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		if runtime.GOOS != "windows" && cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		} else if cmd.Process != nil {
			cmd.Process.Kill()
		}
		<-done
		return buf.String(), 124, ctx.Err()

	case err := <-done:
		exitCode := 0
		if err != nil {
			var ee *exec.ExitError
			if errors.As(err, &ee) {
				exitCode = ee.ExitCode()
			} else {
				exitCode = 1
			}
		}
		return buf.String(), exitCode, nil
	}
}

// DiscoverModelsCmd returns a Bubble Tea command that discovers available models.
func DiscoverModelsCmd() tea.Cmd {
	return func() tea.Msg {
		out, err := exec.Command("opencode", "models").Output()
		if err != nil {
			// Return error directly
			return fmt.Errorf("falha ao descobrir modelos: %w", err)
		}

		lines := strings.Split(string(out), "\n")
		var models []string
		re := regexp.MustCompile(`^[A-Za-z0-9_-]+/[A-Za-z0-9._-]+$`)

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if re.MatchString(line) {
				models = append(models, line)
			}
		}

		sort.Strings(models)
		// Return the slice directly - TUI layer will convert to DiscoveryMsg
		return models
	}
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// PrioritizeModels reorders models to test free ones first.
func PrioritizeModels(models []string, freeModels map[string]bool, zaiPrefix string) []string {
	var free, zai, other []string

	for _, m := range models {
		if freeModels[m] {
			free = append(free, m)
		} else if strings.HasPrefix(m, zaiPrefix) {
			zai = append(zai, m)
		} else {
			other = append(other, m)
		}
	}

	result := make([]string, 0, len(models))
	result = append(result, free...)
	result = append(result, zai...)
	result = append(result, other...)

	return result
}

// SmartTrim truncates long output strings while preserving start and end.
func SmartTrim(s string, maxKB int) string {
	maxBytes := maxKB * 1024
	if len(s) <= maxBytes {
		return s
	}

	headSize := maxBytes * 2 / 5
	tailSize := maxBytes * 2 / 5

	head := s[:headSize]
	tail := s[len(s)-tailSize:]

	return head + "\n\n...[TRUNCATED]...\n\n" + tail
}

// Truncate shortens a string to a maximum length with an ellipsis.
func Truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

// Padding calculates the leading space for alignment in progress status.
func Padding(n, cur, total int) int {
	digits := len(fmt.Sprint(total))
	return digits - len(fmt.Sprint(cur)) + 1
}

// ExtractProvider extracts the provider prefix from a model name.
func ExtractProvider(model string) string {
	parts := strings.Split(model, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

