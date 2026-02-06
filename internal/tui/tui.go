// Package tui provides the Bubble Tea UI model and rendering logic.
package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"llm-radar/internal/cache"
	"llm-radar/internal/kb"
	"llm-radar/internal/models"
	"llm-radar/internal/worker"
)

// ============================================================================
// ADAPTIVE COLORS
// ============================================================================

var (
	ColorSubtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	ColorHighlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	ColorSuccess   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	ColorWarning   = lipgloss.AdaptiveColor{Light: "#F2B824", Dark: "#F2B824"}
	ColorDanger    = lipgloss.AdaptiveColor{Light: "#F24A4A", Dark: "#F24A4A"}
	ColorInfo      = lipgloss.AdaptiveColor{Light: "#00A1E4", Dark: "#00A1E4"}
)

// ============================================================================
// STYLES
// ============================================================================

var (
	DocStyle     = lipgloss.NewStyle().Margin(1, 2)
	TitleStyle   = lipgloss.NewStyle().Bold(true).Foreground(ColorHighlight)
	StatusStyle  = lipgloss.NewStyle()
	SuccessStyle = StatusStyle.Foreground(ColorSuccess)
	WarningStyle = StatusStyle.Foreground(ColorWarning)
	DangerStyle  = StatusStyle.Foreground(ColorDanger)
	InfoStyle    = StatusStyle.Foreground(ColorInfo)
)

// ============================================================================
// MESSAGE TYPES
// ============================================================================

type (
	// ItemMsg represents a completed model test result
	ItemMsg models.ModelResult
	// WorkerStartMsg indicates a worker has started testing a model
	WorkerStartMsg struct {
		Model string
		Start time.Time
	}
	// DiscoveryMsg contains the list of discovered models
	DiscoveryMsg []string
	// TickMsg is sent periodically for UI updates
	TickMsg time.Time
	// ErrorMsg represents a fatal error
	ErrorMsg error
)

// ============================================================================
// APP MODEL
// ============================================================================

// AppModel is the main Bubble Tea model for the application.
type AppModel struct {
	results       []models.ModelResult
	activeJobs    map[string]time.Time
	models        []string
	total         int
	processed     int32
	discovering   bool
	quitting      bool
	done          bool
	err           error
	runCfg        models.RunConfig
	kb            kb.Compiled
	cache         *cache.ResultCache
	progress      progress.Model
	viewport      viewport.Model
	width         int
	height        int
	workerMsgChan chan tea.Msg
	mu            sync.RWMutex
	appName       string
	version       string
	cacheExpiry   time.Duration
}

// NewAppModel creates and initializes a new AppModel.
func NewAppModel(runCfg models.RunConfig, compiledKB kb.Compiled, appName, version string, cacheExpiry time.Duration) *AppModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	resultCache := cache.New(runCfg.CachePath, cacheExpiry)
	if runCfg.UseCache {
		resultCache.Load()
	}

	return &AppModel{
		runCfg:        runCfg,
		kb:            compiledKB,
		cache:         resultCache,
		progress:      p,
		workerMsgChan: make(chan tea.Msg, 100),
		activeJobs:    make(map[string]time.Time),
		discovering:   true,
		appName:       appName,
		version:       version,
		cacheExpiry:   cacheExpiry,
	}
}

// Init initializes the Bubble Tea model.
func (m *AppModel) Init() tea.Cmd {
	return tea.Batch(
		worker.DiscoverModelsCmd(),
		waitForWorkerMsg(m.workerMsgChan),
		tickCmd(),
	)
}

// Update handles Bubble Tea messages and updates the model.
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "s":
			if m.done {
				if err := m.saveResults(); err == nil {
					// Success
				}
			}
		}
		if m.done {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 20
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		m.resizeViewport()

	case TickMsg:
		if m.done || m.quitting {
			return m, nil
		}
		return m, tickCmd()

	// Handle []string from worker.DiscoverModelsCmd
	case []string:
		m.discovering = false
		m.models = msg
		m.total = len(msg)

		// Create map of free models for prioritization
		freeMap := make(map[string]bool)
		for k := range m.kb.Config.FreeModels {
			freeMap[k] = true
		}
		m.models = worker.PrioritizeModels(m.models, freeMap, "zai-coding-plan/")

		go worker.StartWorkers(m.models, m.runCfg, m.kb, m.cache, m.workerMsgChan, &m.processed)
		return m, nil

	case DiscoveryMsg:
		m.discovering = false
		m.models = msg
		m.total = len(msg)

		// Create map of free models for prioritization
		freeMap := make(map[string]bool)
		for k := range m.kb.Config.FreeModels {
			freeMap[k] = true
		}
		m.models = worker.PrioritizeModels(m.models, freeMap, "zai-coding-plan/")

		go worker.StartWorkers(m.models, m.runCfg, m.kb, m.cache, m.workerMsgChan, &m.processed)
		return m, nil

	// Handle generic worker start message from worker package
	case struct{ Model string; Start time.Time }:
		m.mu.Lock()
		m.activeJobs[msg.Model] = msg.Start
		m.mu.Unlock()
		return m, waitForWorkerMsg(m.workerMsgChan)

	case WorkerStartMsg:
		m.mu.Lock()
		m.activeJobs[msg.Model] = msg.Start
		m.mu.Unlock()
		return m, waitForWorkerMsg(m.workerMsgChan)

	// Handle ModelResult directly from worker package
	case models.ModelResult:
		m.mu.Lock()
		delete(m.activeJobs, msg.Model)
		m.results = append(m.results, msg)
		m.mu.Unlock()

		processed := atomic.LoadInt32(&m.processed)
		cmd := m.progress.SetPercent(float64(processed) / float64(m.total))
		m.viewport.SetContent(m.renderResultsList())
		m.viewport.GotoBottom()

		if int(processed) >= m.total {
			m.done = true
			m.mu.Lock()
			m.activeJobs = make(map[string]time.Time)
			m.mu.Unlock()
			m.resizeViewport()

			if m.runCfg.UseCache {
				m.cache.SaveResults(m.results)
			}

			return m, cmd
		}
		return m, tea.Batch(waitForWorkerMsg(m.workerMsgChan), cmd)

	case ItemMsg:
		res := models.ModelResult(msg)
		m.mu.Lock()
		delete(m.activeJobs, res.Model)
		m.results = append(m.results, res)
		m.mu.Unlock()

		processed := atomic.LoadInt32(&m.processed)
		cmd := m.progress.SetPercent(float64(processed) / float64(m.total))
		m.viewport.SetContent(m.renderResultsList())
		m.viewport.GotoBottom()

		if int(processed) >= m.total {
			m.done = true
			m.mu.Lock()
			m.activeJobs = make(map[string]time.Time)
			m.mu.Unlock()
			m.resizeViewport()

			if m.runCfg.UseCache {
				m.cache.SaveResults(m.results)
			}

			return m, cmd
		}
		return m, tea.Batch(waitForWorkerMsg(m.workerMsgChan), cmd)

	case ErrorMsg:
		m.err = msg
		return m, tea.Quit
	
	// Handle error from worker  package
	case error:
		m.err = msg
		return m, tea.Quit
	}

	return m, nil
}

// View renders the UI.
func (m *AppModel) View() string {
	if m.err != nil {
		return DangerStyle.Render(fmt.Sprintf("\n❌ Erro fatal: %v\n", m.err))
	}
	if m.discovering {
		title := TitleStyle.Render(fmt.Sprintf("🧪 %s v%s", m.appName, m.version))
		return fmt.Sprintf("\n %s\n\n 🔍 Descobrindo modelos no OpenCode...\n", title)
	}

	processed := int(atomic.LoadInt32(&m.processed))
	pad := strings.Repeat(" ", Padding(0, processed, m.total))
	status := fmt.Sprintf("%s %d/%d", pad, processed, m.total)
	prog := m.progress.View()

	title := TitleStyle.Render(fmt.Sprintf("🧪 %s v%s", m.appName, m.version))
	header := fmt.Sprintf("\n%s\n\n%s %s\n\n", title, prog, status)

	body := m.viewport.View()

	footer := ""
	if!m.done {
		footer = m.renderActiveJobs()
	} else {
		footer = m.renderFinalReport()
	}

	return DocStyle.Render(header + body + "\n" + footer)
}

// ============================================================================
// RENDERING HELPERS
// ============================================================================

func (m *AppModel) resizeViewport() {
	headerHeight := 8
	footerHeight := 6
	if m.done {
		footerHeight = 4
	}
	vpHeight := m.height - headerHeight - footerHeight
	if vpHeight < 5 {
		vpHeight = 5
	}
	m.viewport.Height = vpHeight
	m.viewport.Width = m.width
}

func (m *AppModel) renderActiveJobs() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.activeJobs) == 0 {
		return InfoStyle.Render("\n ⏳ Aguardando workers...")
	}

	var s strings.Builder
	s.WriteString(InfoStyle.Render("\n 🚧 Processando agora:\n"))

	var active []string
	for k := range m.activeJobs {
		active = append(active, k)
	}
	sort.Strings(active)

	maxShow := 5
	for i, model := range active {
		if i >= maxShow {
			remaining := len(active) - maxShow
			s.WriteString(WarningStyle.Render(fmt.Sprintf("    ... e mais %d modelos\n", remaining)))
			break
		}

		duration := time.Since(m.activeJobs[model]).Round(time.Second)

		style := lipgloss.NewStyle().Foreground(ColorSubtle)
		if duration > 10*time.Second {
			style = style.Foreground(ColorWarning)
		}
		if duration > 20*time.Second {
			style = style.Foreground(ColorDanger).Bold(true)
		}

		s.WriteString(style.Render(fmt.Sprintf("    ⟳ %-45s [%s]\n",
			worker.Truncate(model, 40), duration)))
	}
	return s.String()
}

func (m *AppModel) renderResultsList() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var s strings.Builder
	for _, r := range m.results {
		catStyle := GetStyleForCategory(r.Category)

		line := fmt.Sprintf("%s %-40s %s %8s",
			r.Icon,
			worker.Truncate(r.Model, 40),
			catStyle.Render(fmt.Sprintf("[%-16s]", r.Category)),
			r.Duration)

		if r.Category == models.CategoryFreeLimited && r.Reason != "" {
			line += " " + lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Render(fmt.Sprintf("(%s)", worker.Truncate(r.Reason, 25)))
		}

		s.WriteString(line + "\n")
	}
	return s.String()
}

func (m *AppModel) renderFinalReport() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]int)
	var totalDuration time.Duration

	for _, r := range m.results {
		stats[r.Category]++
		totalDuration += time.Duration(r.DurationMs) * time.Millisecond
	}

	var s strings.Builder
	s.WriteString(SuccessStyle.Render(fmt.Sprintf("\n🏁 Concluído - %d modelos testados em %s\n\n",
		m.total, totalDuration.Round(time.Second))))

	for _, cat := range models.AllCategories() {
		if count := stats[cat]; count > 0 {
			style := GetStyleForCategory(cat)
			icon := models.CategoryIcons[cat]
			s.WriteString(style.Render(fmt.Sprintf("  %s %-17s %3d  ", icon, cat, count)))
		}
	}

	s.WriteString("\n\n")

	usable := stats[models.CategoryFree] + stats[models.CategoryFreeLimited] +
		stats[models.CategoryPaid] + stats[models.CategoryAvailable]

	s.WriteString(SuccessStyle.Render(fmt.Sprintf("✨ %d modelos utilizáveis ", usable)))
	s.WriteString(lipgloss.NewStyle().Foreground(ColorSubtle).
		Render(fmt.Sprintf("(%d%%)\n", usable*100/m.total)))

	s.WriteString(lipgloss.NewStyle().Foreground(ColorSubtle).
		Render("\n(q: sair | s: salvar resultados)"))

	return s.String()
}

func (m *AppModel) saveResults() error {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("llm-radar-%s.json", timestamp)

	dir := filepath.Join(os.Getenv("HOME"), ".config", "opencode", "results")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(map[string]interface{}{
		"timestamp": timestamp,
		"version":   m.version,
		"total":     m.total,
		"results":   m.results,
	}, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ============================================================================
// HELPERS
// ============================================================================

// GetStyleForCategory returns the appropriate lipgloss style for a category.
func GetStyleForCategory(cat string) lipgloss.Style {
	switch cat {
	case models.CategoryFree, models.CategoryFreeLimited, models.CategoryPaid, models.CategoryAvailable:
		return SuccessStyle
	case models.CategoryTimeout, models.CategoryNotFound, models.CategoryRateLimited:
		return WarningStyle
	default:
		return DangerStyle
	}
}

// Padding calculates alignment padding for progress display.
func Padding(n, cur, total int) int {
	digits := len(intToString(total))
	return digits - len(intToString(cur)) + 1
}

// intToString is a simple helper (since we can't import strconv for this simple use).
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func waitForWorkerMsg(ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-ch
		if !ok {
			return nil
		}
		return msg
	}
}
