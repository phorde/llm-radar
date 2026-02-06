# Arquitetura - OpenCode Check

## Diagrama de Componentes

```mermaid
graph TB
    subgraph CLI["CLI Layer"]
        Main["main()"]
        Flags["Flag Parser"]
    end

    subgraph TUI["TUI Layer (Bubble Tea)"]
        Model["appModel"]
        View["View Renderer"]
        Progress["Progress Bar"]
        Viewport["Viewport"]
    end

    subgraph Core["Core Logic"]
        Discovery["Model Discovery"]
        Workers["Worker Pool"]
        Classifier["Classifier"]
        KB["Knowledge Base"]
    end

    subgraph External["External Dependencies"]
        OpenCode["OpenCode CLI"]
        Cache["File Cache"]
        Config["KB Config JSON"]
    end

    Main --> Flags
    Main --> Model
    Model --> View
    Model --> Progress
    Model --> Viewport
    
    Model --> Discovery
    Model --> Workers
    Workers --> Classifier
    Classifier --> KB
    
    Discovery --> OpenCode
    Workers --> OpenCode
    KB --> Config
    Model --> Cache
```

## Fluxo de Dados

```mermaid
sequenceDiagram
    participant User
    participant TUI
    participant Discovery
    participant Workers
    participant OpenCode
    participant Classifier
    participant Cache

    User->>TUI: Inicia aplicação
    TUI->>Discovery: discoverModelsCmd()
    Discovery->>OpenCode: opencode models
    OpenCode-->>Discovery: Lista de modelos
    Discovery-->>TUI: discoveryMsg[]
    
    TUI->>Workers: startWorkers()
    
    loop Para cada modelo
        Workers->>Cache: Verificar cache
        alt Cache válido
            Cache-->>Workers: Resultado cached
        else Cache miss
            Workers->>OpenCode: opencode run --model X
            OpenCode-->>Workers: Output + exit code
            Workers->>Classifier: classify()
            Classifier-->>Workers: Categoria + Reason
            Workers->>Cache: Salvar resultado
        end
        Workers-->>TUI: itemMsg
        TUI->>TUI: Atualizar UI
    end
    
    TUI-->>User: Relatório final
```

## Estrutura de Dados

### ModelResult
```go
type ModelResult struct {
    Model      string  // "provider/model-name"
    Provider   string  // "opencode", "groq", etc.
    Category   string  // FREE, TIMEOUT, etc.
    Reason     string  // Descrição legível
    Duration   string  // "1.234s"
    DurationMs int64   // 1234
    Output     string  // Resposta truncada
    ExitCode   int     // 0 = sucesso
    Icon       string  // Emoji da categoria
    Timestamp  string  // RFC3339
}
```

### Categorias

```
┌────────────────┬──────┬─────────────────────────────────┐
│ Categoria      │ Icon │ Significado                     │
├────────────────┼──────┼─────────────────────────────────┤
│ FREE           │ 🆓   │ Gratuito sem limites            │
│ FREE_LIMITED   │ 📊   │ Gratuito com quotas             │
│ PAID_ZAI_OK    │ 💎   │ Pago ZAI ativo                  │
│ AVAILABLE      │ ✅   │ Disponível (outros)             │
│ NOT_FOUND      │ ❓   │ Modelo não existe               │
│ TIMEOUT        │ ⏰   │ Timeout (20s)                   │
│ AUTH_FAILED    │ 🔒   │ API key inválida                │
│ NO_QUOTA       │ ❌   │ Sem créditos                    │
│ RATE_LIMITED   │ ⏱️   │ Rate limit atingido             │
│ ERROR          │ ⚠️   │ Erro desconhecido               │
└────────────────┴──────┴─────────────────────────────────┘
```

## Dependências

| Pacote | Versão | Uso |
|--------|--------|-----|
| bubbletea | v1.3.10 | Framework TUI |
| bubbles | v0.21.1 | Componentes (progress, viewport) |
| lipgloss | v1.1.0 | Estilização |

## Pontos de Extensão

1. **Knowledge Base** - Arquivo JSON externo para customizar classificações
2. **Regex Patterns** - Configuráveis para detectar erros específicos
3. **Cache Strategy** - Expiração e local configuráveis
4. **Workers** - Número ajustável via CLI
