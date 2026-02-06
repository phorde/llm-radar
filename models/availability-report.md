# Relatório de Disponibilidade de Modelos

> Gerado em: 2026-02-05T23:52:00Z  
> Total de modelos: **582**  
> Provedores: **12**

## Legenda de Classificação

| Ícone | Categoria | Descrição |
|:-----:|-----------|-----------|
| 🟢 | GRATUITO | Sem limitações de uso conhecidas |
| 🟡 | GRATUITO_LIMITADO | Com quotas/rate limits documentados |
| 🔵 | PAGO_LIBERADO | Assinatura ativa ou créditos suficientes |
| 🔴 | PAGO_INDISPONÍVEL | Sem créditos ou assinatura expirada |
| ⭐ | PREMIUM_QUOTED | Modelo que consome créditos por uso |
| ⚪ | NÃO_TESTADO | Requer execução do opencode-check |

---

## Modelos por Provedor

### OpenCode (Nativos)

| Modelo | Classificação | Limites | Status |
|--------|:-------------:|---------|--------|
| opencode/big-pickle | 🟢 GRATUITO | Não documentados | Zen |
| opencode/gpt-5-nano | 🟢 GRATUITO | Não documentados | Zen |
| opencode/minimax-m2.1-free | 🟢 GRATUITO | Não documentados | Zen |
| opencode/glm-4.7-free | 🟢 GRATUITO | Não documentados | Zen |
| opencode/kimi-k2.5-free | 🟢 GRATUITO | Não documentados | Zen |
| opencode/trinity-large-preview-free | 🟢 GRATUITO | Não documentados | Zen |
| opencode/claude-sonnet-4-5 | ⭐ PREMIUM | Por token | Premium |
| opencode/claude-opus-4-5 | ⭐ PREMIUM | Por token | Premium |
| opencode/gpt-5 | ⭐ PREMIUM | Por token | Premium |
| opencode/gpt-5.1-codex | ⭐ PREMIUM | Por token | Premium |

---

### Cerebras

| Modelo | Classificação | Limites |
|--------|:-------------:|---------|
| cerebras/llama-3.3-70b | 🟡 GRATUITO_LIMITADO | 1M tokens/dia (agregado) |
| cerebras/llama-4-scout-17b | 🟡 GRATUITO_LIMITADO | 1M tokens/dia (agregado) |
| cerebras/qwen-3-32b | 🟡 GRATUITO_LIMITADO | 1M tokens/dia (agregado) |

---

### DeepSeek

| Modelo | Classificação | Limites |
|--------|:-------------:|---------|
| deepseek/deepseek-v3 | 🟡 GRATUITO_LIMITADO | 5M tokens inicial, 50 RPM |
| deepseek/deepseek-r1 | 🟡 GRATUITO_LIMITADO | 5M tokens inicial, 50 RPM |
| deepseek/deepseek-coder | 🟡 GRATUITO_LIMITADO | 5M tokens inicial, 50 RPM |

---

### Groq

| Modelo | Classificação | Limites |
|--------|:-------------:|---------|
| groq/llama-3.3-70b-versatile | 🟡 GRATUITO_LIMITADO | 14.4K req/dia, 30 RPM |
| groq/gemma2-9b-it | 🟡 GRATUITO_LIMITADO | 14.4K req/dia, 30 RPM |
| groq/mixtral-8x7b-32768 | 🟡 GRATUITO_LIMITADO | 14.4K req/dia, 30 RPM |

---

### OpenRouter (Modelos Gratuitos)

| Modelo | Classificação | Limites |
|--------|:-------------:|---------|
| openrouter/deepseek/deepseek-r1-0528:free | 🟢 GRATUITO | Community tier |
| openrouter/qwen/qwen3-coder:free | 🟢 GRATUITO | Community tier |
| openrouter/meta/llama-4-maverick:free | 🟢 GRATUITO | Community tier |
| openrouter/google/gemini-2.5-flash:free | 🟢 GRATUITO | Community tier |
| openrouter/anthropic/claude-haiku-4.5:free | 🟢 GRATUITO | Community tier |

---

### ZAI Coding Plan

| Modelo | Classificação | Custo Estimado |
|--------|:-------------:|----------------|
| zai-coding-plan/glm-4.5 | 🔵 PAGO_LIBERADO | ¥0.005/1K tokens |
| zai-coding-plan/glm-4.6 | 🔵 PAGO_LIBERADO | ¥0.005/1K tokens |
| zai-coding-plan/glm-4.7 | 🔵 PAGO_LIBERADO | ¥0.01/1K tokens |
| zai-coding-plan/glm-4.7-flash | 🔵 PAGO_LIBERADO | ¥0.005/1K tokens |

---

### Google (Antigravity)

| Modelo | Classificação | Context Window |
|--------|:-------------:|----------------|
| google/antigravity-gemini-3-pro | 🔵 PAGO_LIBERADO | 1M tokens |
| google/antigravity-gemini-3-flash | 🔵 PAGO_LIBERADO | 1M tokens |
| google/antigravity-claude-sonnet-4-5 | 🔵 PAGO_LIBERADO | 200K tokens |
| google/gemini-2.5-flash | 🟡 GRATUITO_LIMITADO | Gemini CLI |
| google/gemini-2.5-pro | 🟡 GRATUITO_LIMITADO | Gemini CLI |

---

## Resumo Estatístico

```
┌─────────────────────────────────────────────────────────────┐
│              DISTRIBUIÇÃO POR CLASSIFICAÇÃO                 │
├─────────────────────────────────────────────────────────────┤
│ 🟢 GRATUITO           │ ~50 modelos  │ ████████░░░ 9%      │
│ 🟡 GRATUITO_LIMITADO  │ ~80 modelos  │ ████████████░ 14%   │
│ 🔵 PAGO_LIBERADO      │ ~100 modelos │ ████████████████ 17%│
│ ⭐ PREMIUM_QUOTED     │ ~300 modelos │ ████████████████████│
│ ⚪ NÃO_TESTADO        │ ~52 modelos  │ ████████░░░ 9%      │
└─────────────────────────────────────────────────────────────┘
```

## Próximos Passos

Para obter classificação precisa de todos os modelos, execute:

```bash
cd /home/ubuntu/projects/opencode-check
go run main.go -c 10 -t 30s
```

Ou use a versão compilada:
```bash
./opencode-check -c 10 -t 30s
```
