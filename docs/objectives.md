# Objetivos SMART - OpenCode Check

## Objetivo Principal

Fornecer aos usuários uma visão consolidada da disponibilidade de modelos LLM, com uma meta de performance alvo de < 60 segundos para ~50 modelos em ambientes padrão.

---

## Objetivos Específicos

### 1. Descoberta de Modelos

| Componente | Meta |
|------------|------|
| **Específico** | Descobrir automaticamente modelos registrados no OpenCode via `models.dev` e `ai-sdk` |
| **Mensurável** | Cobertura total de modelos retornados pelo backend do OpenCode |
| **Atingível** | Via comando `opencode models` |
| **Relevante** | Sincronismo com o ecossistema Vercel AI SDK e providers suportados |
| **Temporal** | Em menos de 5 segundos |

**KPI:** Taxa de descoberta ≥ 99%

---

### 2. Teste de Disponibilidade (Escalabilidade)

| Componente | Meta |
|------------|------|
| **Específico** | Testar cada modelo com prompt mínimo ("2, 3, 5") |
| **Mensurável** | Latência p95 por modelo < 15s |
| **Atingível** | Via workers paralelos (padrão: dinâmico baseado em CPU/Hardware) |
| **Relevante** | Identificar "cold starts" e disponibilidade real |
| **Temporal** | Tempo total proporcional à concorrência (ex: 50 modelos / 5 workers * 5s med. = 50s) |

> [!NOTE]
> O número de workers é dinamicamente definido pelo setup do usuário para evitar lagging na TUI e otimizar throughput.

---

### 3. Classificação e Taxonomia

| Componente | Meta |
|------------|------|
| **Específico** | Categorizar resultados em 11 classes fundamentais (veja abaixo) |
| **Mensurável** | Acurácia de mapeamento de erro ≥ 90% |
| **Atingível** | Via análise de output (Regex como ferramenta de compatibilidade universal) |
| **Relevante** | Decisão de uso baseada em custo e cota |
| **Temporal** | Classificação instantânea pós-teste |

#### Definição de Categorias
1.  **FREE**: Modelos sem custo e sem limites documentados (ex: Zen models).
2.  **FREE_LIMITED**: Gratuitos, mas com quotas diárias ou limites de RPM estritos.
3.  **PAID**: Modelos pagos que exigem saldo ou assinatura ativa.
4.  **AVAILABLE**: Sucesso no teste, mas sem metadados específicos na KB.
5.  **NOT_FOUND**: Modelo inexistente ou removido do provider.
6.  **TIMEOUT**: Estouro do limite de tempo configurado (default 20s).
7.  **AUTH_FAILED**: Erro de credencial (API Key inválida ou expirada).
8.  **NO_QUOTA**: Saldo insuficiente ou limite de faturamento atingido.
9.  **RATE_LIMITED**: Throttling ativo por limite de volume.
10. **FREE_ERROR**: Falha técnica em modelo conhecido como gratuito.
11. **ERROR**: Outros erros (rede, protocolo, falha interna do provider).

> [!TIP]
> **Por que Regex?** Dada a diversidade de providers (Groq, Anthropic, Google, etc.) cujas saídas de erro no OpenCode CLI nem sempre são estruturadas (JSON), o Regex oferece a maior flexibilidade de adaptação sem exigir mudanças nos backends dos providers.

---

### 4. Experiência do Usuário (UX)

| Componente | Meta |
|------------|------|
| **Específico** | TUI Responsiva com feedback visual imediato para cada ação |
| **Mensurável** | FPS constante (~30), zero bloqueio de input |
| **Atingível** | Via Bubble Tea (Event-driven Architecture) |
| **Relevante** | Transparência total sobre o estado do processo |
| **Temporal** | Feedback instantâneo em < 100ms para inputs do usuário |

**Critério:** Todo erro ou ação (salvamento, cancelamento) deve disparar uma animação ou mudança de estado visual clara.

---

## Métricas de Sucesso

```
┌─────────────────────────────────────────────────────────────┐
│                    DASHBOARD DE MÉTRICAS                    │
├─────────────────────┬───────────────┬───────────────────────┤
│ Métrica             │ Meta          │ Frequência            │
├─────────────────────┼───────────────┼───────────────────────┤
│ Taxa de descoberta  │ ≥ 99%         │ Por execução          │
│ Tempo Per Modelo    │ < 15s (p95)   │ Por execução          │
│ Precisão Classif.   │ ≥ 90%         │ Semanal (validação)   │
│ Estabilidade TUI    │ 0 Lags        │ Por execução          │
│ Feedback Visual     │ 100% ações    │ Por execução          │
└─────────────────────┴───────────────┴───────────────────────┘
```
