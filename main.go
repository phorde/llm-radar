package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"llm-radar/internal/kb"
	"llm-radar/internal/models"
	"llm-radar/internal/tui"
)

// ============================================================================
// VERSION AND METADATA
// ============================================================================

const (
	Version     = "0.1.0"
	AppName     = "LLM Radar"
	CacheExpiry = 24 * time.Hour
)

// Build variables (injected via ldflags)
var BuildTime = "dev"

// ============================================================================
// MAIN
// ============================================================================

func main() {
	parallel := flag.Int("c", 0, "Número de workers paralelos (0 = automático)")
	timeoutFlag := flag.Duration("t", 20*time.Second, "Timeout por modelo")
	refresh := flag.Bool("refresh", false, "Atualizar lista de modelos")
	kbFile := flag.String("kb", "", "Arquivo JSON com KB customizada")
	useCache := flag.Bool("cache", false, "Usar cache (válido 24h)")
	version := flag.Bool("version", false, "Mostrar versão")

	flag.Parse()

	if *version {
		fmt.Printf("%s v%s\n", AppName, Version)
		os.Exit(0)
	}

	concurrency := *parallel
	if concurrency <= 0 {
		concurrency = runtime.NumCPU()
		if concurrency > 8 { // Cap to avoid overwhelming
			concurrency = 8
		}
		if concurrency < 2 {
			concurrency = 2
		}
	}

	compiledKB, err := kb.LoadAndCompile(*kbFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro ao carregar KB: %v\n", err)
		os.Exit(1)
	}

	homeDir, _ := os.UserHomeDir()
	cachePath := filepath.Join(homeDir, ".config", "opencode", "cache", "results.json")

	runCfg := models.RunConfig{
		Prompt:      "Escreva apenas: 2, 3, 5",
		Timeout:     *timeoutFlag,
		Concurrency: concurrency,
		Retries:     1,
		MaxOutputKB: 64,
		UseCache:    *useCache,
		CachePath:   cachePath,
	}

	if *refresh {
		fmt.Println("🔄 Atualizando lista de modelos...")
		if err := exec.Command("opencode", "models", "--refresh").Run(); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Aviso: falha ao atualizar: %v\n", err)
		}
	}

	model := tui.NewAppModel(runCfg, compiledKB, AppName, Version, CacheExpiry)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro TUI: %v\n", err)
		os.Exit(1)
	}
}
