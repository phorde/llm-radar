# 📡 LLM Radar

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Note**: This project is **not** officially built by or affiliated with the OpenCode team (Anomaly). It is a community tool for discovering and analyzing LLM model availability across providers.

A CLI tool with an interactive TUI that automatically tests and classifies the availability of LLM models through the OpenCode CLI. Get a complete view of which models are accessible, free, rate-limited, or require authentication in under 60 seconds.

> [!IMPORTANT]
> **Security**: This tool executes the OpenCode CLI with your configured credentials. Never run with untrusted model configurations or knowledge base files. API keys are managed via OpenCode's credential system (`~/.config/opencode/opencode.json`)—review with `opencode auth` before use. See [`.env.example`](./.env.example) for optional environment overrides.

**🇵🇹 Versão em Português:** [README.pt-BR.md](./README.pt-BR.md)

## ✨ Features

- 🔍 **Automatic Discovery** - Finds all configured models via `opencode models`
- ⚡ **Parallel Testing** - Tests multiple models simultaneously with configurable workers
- 📊 **Smart Classification** - Categorizes models into 12 distinct states (FREE, PAID, TIMEOUT, etc.)
- 🎨 **Real-time TUI** - Beautiful terminal interface with progress bars and live updates
- 💾 **Intelligent Caching** - Reuses results for 24 hours to speed up subsequent runs
- 🧠 **Extensible Knowledge Base** - Customize model classifications via JSON config

## 📋 Prerequisites

- **OpenCode CLI** installed and configured ([opencode.ai/docs](https://opencode.ai/docs))
- **Go 1.24+** (for building from source)
- A modern terminal emulator (WezTerm, Alacritty, Ghostty, Kitty, etc.)

## 🚀 Installation

### Option 1: Download Pre-built Binary

Download the latest release from [Releases](https://github.com/USERNAME/llm-radar/releases) and add to your PATH.

```bash
# Linux/macOS
curl -L https://github.com/USERNAME/llm-radar/releases/latest/download/llm-radar-linux-amd64 -o llm-radar
chmod +x llm-radar
sudo mv llm-radar /usr/local/bin/

# Verify installation
llm-radar --version
```

### Option 2: Build from Source

```bash
git clone https://github.com/USERNAME/llm-radar.git
cd llm-radar
go build -o llm-radar
sudo mv llm-radar /usr/local/bin/
```

### Option 3: Install with Go

```bash
go install github.com/USERNAME/llm-radar@latest
```

## 📖 Usage

### Basic Usage

```bash
# Test all available models
llm-radar

# Use cache to speed up subsequent runs
llm-radar --cache

# Refresh model list before testing
llm-radar --refresh
```

### Advanced Options

```bash
# Customize number of parallel workers (default: 5)
llm-radar -c 10

# Adjust timeout per model (default: 20s)
llm-radar -t 30s

# Use custom knowledge base
llm-radar --kb custom-kb.json

# Show version
llm-radar --version
```

### Flags Reference

| Flag | Default | Description |
|------|---------|-------------|
| `-c` | `5` | Number of parallel workers |
| `-t` | `20s` | Timeout per model |
| `--cache` | `false` | Use cached results (valid for 24h) |
| `--refresh` | `false` | Refresh model list before testing |
| `--kb` | `""` | Path to custom knowledge base JSON |
| `--version` | - | Show version information |

## 📊 Model Categories

Results are classified into 12 categories:

| Icon | Category | Meaning |
|------|----------|---------|
| 🆓 | `FREE` | Free models without known limits |
| 📊 | `FREE_LIMITED` | Free with quotas (e.g., Cerebras, DeepSeek, Groq) |
| 💰 | `PAID` | Paid ZAI models with active credits |
| ✅ | `AVAILABLE` | Available (general) |
| ❓ | `NOT_FOUND` | Model doesn't exist |
| ⏰ | `TIMEOUT` | Timeout (20s default) |
| 🔒 | `AUTH_FAILED` | Invalid API key |
| ❌ | `NO_QUOTA` | No credits remaining |
| ⏱️ | `RATE_LIMITED` | Rate limit reached |
| ⚠️ | `ERROR` | Unknown error |

## 🎨 Screenshots

### Running Test
```
📡 LLM Radar v0.1.0

████████████████████████████░░░░  75% 30/40

🆓 opencode/big-pickle           [FREE           ]   1.234s
📊 cerebras/llama3.3-70b         [FREE_LIMITED   ]   2.567s
💰 zai-coding-plan/glm-4.7       [PAID           ]   1.890s

🚧 Processing now:
    ⟳ groq/llama-3.3-70b-versatile        [5s]
    ⟳ deepseek/deepseek-chat              [3s]
```

### Final Report
```
🏁 Completed - 40 models tested in 45s

  🆓 FREE              6    📊 FREE_LIMITED     8    💰 PAID      3  
  ❓ NOT_FOUND        12    ⏰ TIMEOUT          5    🔒 AUTH_FAILED      6  

✨ 17 usable models (42%)

(q: quit | s: save results)
```

## 🔧 Configuration

### Custom Knowledge Base

Create a JSON file to override default classifications:

```json
{
  "free_models": {
    "opencode/my-custom-model": {
      "category": "FREE",
      "description": "My Custom Model",
      "limits": "no documented limits"
    }
  },
  "free_tier_providers": {
    "myprovider": {
      "category": "FREE_LIMITED",
      "description": "My Provider",
      "limits": "1M tokens/day"
    }
  },
  "zai_models": {
    "zai-coding-plan/custom": {
      "category": "PAID",
      "description": "Custom ZAI Model"
    }
  }
}
```

Use with:
```bash
opencode-check --kb custom-kb.json
```

## 📁 Output Files

Results are saved to:
- **Cache**: `~/.config/opencode/cache/results.json` (when using `--cache`)
- **Reports**: `~/.config/opencode/results/opencode-check-YYYYMMDD-HHMMSS.json` (press `s` to save)

## 🤝 Compatibility with OpenCode Plugins

This tool works alongside popular OpenCode plugins:

- **[opencode-antigravity-auth](https://github.com/NoeFabris/opencode-antigravity-auth)** - Tests Antigravity OAuth models
- **[oh-my-opencode](https://github.com/code-yeongyu/oh-my-opencode)** - Compatible with enhanced agent features

## 🗺️ Roadmap & Next Steps

This project is in active development. Here are planned enhancements aligned with the vision of providing comprehensive, automated LLM model discovery and classification.

### Phase 2 — Automated Discovery (v0.2.0)

#### 🌐 Provider API Discovery
**Goal**: Eliminate manual knowledge base maintenance by auto-discovering provider metadata for ~50% of OpenCode providers, prioritizing free tier offerings.

**Target Providers** (focus on free/freemium tiers):
- **OpenCode Zen Models**: `opencode/*-free` endpoints (already well-documented)
- **OpenRouter**: Query `/api/v1/models` → filter by `pricing.prompt === 0`
- **Vercel AI**: Parse model catalog for free tier indicators
- **GitHub Copilot**: Models accessible with GitHub subscription
- **Google AI Studio**: Gemini free tier limits (15 RPM, 1M tokens/day)
- **AIHubMix**: Community-curated free model aggregator
- **Ollama Cloud**: Parse public model registry
- **Hugging Face Inference**: Free tier detection via HF Hub API
- **Additional coverage**: Groq, Cerebras, DeepSeek, Moonshot (if free tiers available)

**Implementation Strategy**:
1. **Provider APIs**: Query structured endpoints (`/v1/models`, `/api/models`)
2. **Public Documentation**: Scrape official docs for rate limits
3. **Community Sources**: Leverage OpenCode community KB
4. **Fallback**: Cached KB for offline/degraded scenarios

**Benefits**: Covers majority of free LLM access points, reduces manual updates, adapts to provider changes automatically.

#### 🏆 Tier Classification (S/A/B/C)
**Goal**: Rank models by performance using benchmark data, not just availability.

**Metrics** (weighted composite score):
- **Latency** (time-to-first-token): 25% weight
- **Throughput** (tokens/second): 20% weight  
- **Reliability** (uptime over 7 days): 20% weight
- **Cost Efficiency** (quality/$ for paid, quality/quota for free): 25% weight
- **Success Rate** (% of requests completed): 10% weight

**Data Sources**:
1. **Integrated Benchmarks**: Lightweight probe prompts (e.g., "Solve: 2+2", "Explain: quantum entanglement")
2. **External Benchmark Imports**: 
   - Artificial Analysis (latency/pricing data)
   - LMSYS Chatbot Arena (ELO scores)
   - OpenLLM Leaderboard (task-specific performance)
3. **Historical Data**: LLM Radar's own tracking database

**Scoring Tiers**:
- **S-Tier** (top 10%): Best-in-class performance + cost efficiency
- **A-Tier** (next 15%): Excellent performance, minor trade-offs
- **B-Tier** (next 25%): Solid general-purpose models
- **C-Tier** (remaining 50%): Functional but with notable limitations

**UI Update**: Display tiers with colored badges (🔴 S, 🟠 A, 🟡 B, ⚪ C) in results table.

#### 🎯 Specialization Tags (Cascade Classification)
**Goal**: Help users find the right model for their task using multi-signal classification.

**Categories**: `coding`, `reasoning`, `multimodal`, `conversation`, `speed`, `math`

**Detection Cascade** (priority order):
1. **Benchmark Performance** (highest priority):
   - HumanEval scores → `coding` (e.g., Claude 3.5 Sonnet despite no "code" in name)
   - MATH dataset scores → `math`
   - MMLU scores → `reasoning`
   - Vision benchmarks → `multimodal`

2. **Provider Metadata**:
   - Official model cards (GitHub Copilot → `coding`)
   - Provider classifications (Vercel AI, OpenRouter tags)

3. **Model Name Heuristics** (lowest priority):
   - `code-`, `-coder`, `codestral` → `coding`
   - `vision`, `-v`, `gpt-4o` → `multimodal`
   - `turbo`, `flash`, `instant` → `speed`

4. **Lightweight Probe Prompts** (fallback):
   - Coding: "Write a function to reverse a string"
   - Reasoning: "Explain the Ship of Theseus paradox"
   - Math: "Solve: ∫(x² + 3x)dx"

**Rationale**: Benchmark data is most reliable (e.g., Claude Opus 3.5 excels at coding despite generic name). Name-based detection is last resort to avoid false negatives.

**Display**: Icon badges in TUI with confidence scores (🖥️ coding:95%, 🧠 reasoning:78%, etc.)

### Phase 3 — Ecosystem Integration (v0.3.0)

#### 🔌 OpenCode Plugin
**Goal**: Official integration as `opencode models --check`.

- **Package**: Publish to OpenCode plugin registry
- **API**: Expose check results as structured JSON for other plugins
- **Benefits**: Seamless UX for existing OpenCode users
- **Maintenance**: Coordinate with Anomaly team for API stability

#### 📊 Export Formats
**Goal**: Enable CI/CD automation.

- **Formats**: JSON, CSV, Markdown, HTML report
- **Use cases**:
  - Alert when free models drop below threshold (CI script)
  - Generate weekly availability reports (cron)
  - Export to spreadsheet for cost analysis
- **Implementation**: `--export json --output report.json` flag

#### 📈 Historical Tracking
**Goal**: Identify availability trends.

- **Storage**: Local SQLite database (`~/.config/opencode/check.db`)
- **Schema**: `(timestamp, model, category, duration, success)`
- **Queries**: 
  - "Which models had uptime > 95% last month?"
  - "When did Groq rate limits increase?"
- **Visualization**: ASCII charts in TUI or web dashboard

### Explicitly Out of Scope

To maintain focus, these are **NOT** planned:
- ❌ Model fine-tuning or training (not a dev tool concern)
- ❌ Direct LLM API calls bypassing OpenCode (defeats the abstraction)
- ❌ Provider authentication management (use `opencode auth`)
- ❌ Response quality benchmarking (use dedicated eval frameworks like HELM/BBQ)

### Contributing

Interested in implementing these features? See [CONTRIBUTING.md](./CONTRIBUTING.md) for development guidelines. Priority areas:
1. Provider API scrapers (start with Groq/Cerebras)
2. Benchmark harness for tier classification
3. SQLite schema design for historical tracking

---

## 🐛 Troubleshooting

### "Failed to discover models"
- Ensure OpenCode CLI is installed: `opencode --version`
- Check that you've configured at least one provider: `opencode models`

### Models showing as "NOT_FOUND"
- Run `opencode models --refresh` to update the model list
- Some models may have been deprecated or renamed

### High timeout rates
- Increase timeout: `opencode-check -t 30s`
- Check your internet connection
- Some providers may be temporarily unavailable

### Cache issues
- Clear cache: `rm ~/.config/opencode/cache/results.json`
- Disable cache: don't use `--cache` flag

## 🏗️ Local Development

### Prerequisites

1. **Go 1.24 or later**
   ```bash
   go version  # Should show 1.24 or higher
   ```

2. **OpenCode CLI** (for integration testing)
   ```bash
   # Install OpenCode if you haven't already
   # Follow instructions at https://opencode.ai/docs
   opencode --version
   ```

### Setup

```bash
# 1. Clone the repository
git clone https://github.com/your-username/opencode-check.git
cd opencode-check

# 2. Install dependencies
go mod download

# 3. Verify setup
make check
```

### Building

```bash
# Build the binary
make build

# The binary will be created at: bin/opencode-check
```

### Running Locally

```bash
# Run directly with go
go run . --help

# Or use the built binary
./bin/opencode-check

# Run with specific options
./bin/opencode-check -c 3 -t 10s --cache
```

### Testing

```bash
# Run all unit tests
make test

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -v ./test/...

# Run specific package tests
go test -v ./internal/classifier/

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Project Structure

```
opencode-check/
├── main.go                    # Entry point (~94 lines)
├── internal/
│   ├── cache/                 # Result caching system
│   │   ├── cache.go
│   │   └── cache_test.go
│   ├── classifier/            # Model classification logic
│   │   ├── classifier.go
│   │   └── classifier_test.go
│   ├── kb/                    # Knowledge base handling
│   │   ├── kb.go
│   │   └── kb_test.go
│   ├── models/                # Data structures
│   │   ├── models.go
│   │   └── models_test.go
│   ├── tui/                   # Bubble Tea UI
│   │   └── tui.go
│   └── worker/                # Parallel execution engine
│       ├── worker.go
│       └── worker_test.go
├── test/                      # Integration tests
│   └── integration_test.go
├── docs/                      # Documentation
│   ├── architecture.md
│   ├── objectives.md
│   └── functional-requirements.md
└── Makefile                   # Build automation
```

### Development Workflow

1. **Make changes** to the code
2. **Run tests** to ensure nothing breaks:
   ```bash
   make test
   ```
3. **Build** the binary:
   ```bash
   make build
   ```
4. **Test manually** with your changes:
   ```bash
   ./bin/opencode-check
   ```
5. **Verify code quality**:
   ```bash
   go vet ./...
   go fmt ./...
   ```

### Makefile Commands

```bash
# Build the project
make build

# Run all tests
make test

# Clean build artifacts
make clean

# Run all checks (vet + test)
make check
```

### Debugging

Enable verbose output by modifying the code or adding debug prints:

```go
// In internal/worker/worker.go
fmt.Printf("DEBUG: Testing model %s\n", modelName)
```

Run with additional flags for troubleshooting:
```bash
# Test with a single worker to see sequential output
./bin/opencode-check -c 1

# Increase timeout for slow connections
./bin/opencode-check -t 60s
```

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built for the [OpenCode](https://opencode.ai) community
- Uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) for TUI
- Inspired by the amazing OpenCode plugin ecosystem

## 🔗 Related Projects

- [OpenCode Official Repo](https://github.com/anomalyco/opencode)
- [Awesome OpenCode](https://github.com/awesome-opencode/awesome-opencode)
- [OpenCode Documentation](https://opencode.ai/docs)

---

**Disclaimer**: This is an independent community tool and is not officially maintained by or affiliated with Anomaly (the creators of OpenCode).
