# Recomendações de Modelos por Tarefa

Sugestões de modelos otimizados para diferentes casos de uso baseadas em custo/desempenho.

---

## 🏆 Melhores Escolhas por Categoria

### Codificação Geral

| Prioridade | Modelo | Motivo |
|:----------:|--------|--------|
| 1° | `openrouter/qwen/qwen3-coder:free` | Gratuito, especializado em código |
| 2° | `deepseek/deepseek-v3` | Excelente custo-benefício |
| 3° | `cerebras/llama-3.3-70b` | Rápido, 1M tokens/dia grátis |

### Velocidade Máxima

| Prioridade | Modelo | Latência Esperada |
|:----------:|--------|-------------------|
| 1° | `groq/llama-3.3-70b-versatile` | ~100ms |
| 2° | `cerebras/*` | ~200ms |
| 3° | `opencode/gpt-5-nano` | ~500ms |

### Contexto Longo (>100K tokens)

| Prioridade | Modelo | Context Window |
|:----------:|--------|----------------|
| 1° | `google/antigravity-gemini-3-pro` | 1M tokens |
| 2° | `google/gemini-2.5-pro` | 1M tokens |
| 3° | `anthropic/claude-sonnet-4-5` | 200K tokens |

### Gratuitos Ilimitados (Zen)

| Modelo | Uso Recomendado |
|--------|-----------------|
| `opencode/big-pickle` | Tarefas gerais |
| `opencode/gpt-5-nano` | Testes rápidos |
| `opencode/trinity-large-preview-free` | Raciocínio complexo |

---

## 💰 Otimização de Custos

### Estratégia de Fallback

```
Nível 1 (Gratuito)
    ↓ se limite atingido
Nível 2 (Free Tier)
    ↓ se limite atingido  
Nível 3 (Pago Barato)
    ↓ se necessário
Nível 4 (Premium)
```

### Implementação Sugerida

```go
// Ordem de prioridade por custo
var modelPriority = []string{
    // Gratuitos ilimitados
    "opencode/big-pickle",
    "opencode/gpt-5-nano",
    
    // Free tier com limites
    "cerebras/llama-3.3-70b",
    "deepseek/deepseek-v3",
    "groq/llama-3.3-70b-versatile",
    
    // Pagos baratos
    "zai-coding-plan/glm-4.7-flash",
    
    // Premium (último recurso)
    "opencode/claude-sonnet-4-5",
}
```

---

## 🎯 Recomendações por Projeto

### Desenvolvimento Local / Testes
- **Primário:** `opencode/big-pickle`, `opencode/gpt-5-nano`
- **Backup:** `cerebras/*`, `groq/*`
- **Custo:** $0

### Produção com Volume Moderado
- **Primário:** `deepseek/deepseek-v3`
- **Backup:** `zai-coding-plan/glm-4.7`
- **Custo:** ~$10-50/mês

### Produção Crítica
- **Primário:** `google/antigravity-gemini-3-pro`
- **Backup:** `opencode/claude-sonnet-4-5`
- **Custo:** ~$100-500/mês

---

## ⚠️ Notas Importantes

1. **Rate Limits:** Cerebras, Groq e provedores free tier têm limites rígidos
2. **Latência:** Modelos via OpenRouter podem ter latência adicional
3. **Disponibilidade:** Modelos `:free` do OpenRouter podem ficar indisponíveis
4. **Custos:** Preços podem variar; verifique documentação oficial

---

## 📊 Atualização

Para atualizar este relatório com dados reais:

```bash
./opencode-check --refresh
```
