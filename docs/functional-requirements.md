# Requisitos Funcionais - OpenCode Check

Priorização usando método **MoSCoW** (Must/Should/Could/Won't).

---

## Must Have (Essencial)

### RF-001: Descoberta Automática de Modelos
- **Descrição:** Listar todos os modelos disponíveis via `opencode models`
- **Critério de Aceite:** Retornar lista com formato `provider/model`
- **Prioridade:** P0

### RF-002: Teste de Conectividade
- **Descrição:** Executar prompt mínimo em cada modelo descoberto
- **Critério de Aceite:** Capturar output, exit code e tempo de resposta
- **Prioridade:** P0

### RF-003: Classificação de Modelos
- **Descrição:** Categorizar resultados em classes predefinidas
- **Classes:**
  - `FREE` - Gratuito sem limites conhecidos
  - `FREE_LIMITED` - Gratuito com quotas
  - `PAID_ZAI_OK` - Pago com créditos disponíveis
  - `AVAILABLE` - Disponível (outros)
  - `NOT_FOUND` - Modelo não encontrado
  - `TIMEOUT` - Timeout na requisição
  - `AUTH_FAILED` - Falha de autenticação
  - `NO_QUOTA` - Sem créditos
  - `RATE_LIMITED` - Rate limit atingido
  - `ERROR` - Erro genérico
- **Prioridade:** P0

### RF-004: Interface TUI
- **Descrição:** Exibir progresso em tempo real com barra e lista
- **Critério de Aceite:** Atualização visual a cada resultado
- **Prioridade:** P0

---

## Should Have (Importante)

### RF-005: Workers Paralelos
- **Descrição:** Executar testes em paralelo via flag `-c`
- **Padrão:** 5 workers
- **Critério de Aceite:** Reduzir tempo total proporcionalmente
- **Prioridade:** P1

### RF-006: Timeout Configurável
- **Descrição:** Limitar tempo de espera por modelo via flag `-t`
- **Padrão:** 20 segundos
- **Prioridade:** P1

### RF-007: Base de Conhecimento Extensível
- **Descrição:** Carregar classificações customizadas via JSON (`-kb`)
- **Critério de Aceite:** Merge com defaults internos
- **Prioridade:** P1

### RF-008: Persistência de Resultados
- **Descrição:** Salvar resultados em JSON ao pressionar `s`
- **Local:** `~/.config/opencode/results/`
- **Prioridade:** P1

---

## Could Have (Desejável)

### RF-009: Sistema de Cache
- **Descrição:** Reutilizar resultados recentes via flag `--cache`
- **Expiração:** 24 horas
- **Local:** `~/.config/opencode/cache/results.json`
- **Prioridade:** P2

### RF-010: Atualização de Lista
- **Descrição:** Forçar refresh via `--refresh`
- **Critério de Aceite:** Executar `opencode models --refresh`
- **Prioridade:** P2

### RF-011: Priorização Inteligente
- **Descrição:** Testar modelos gratuitos primeiro
- **Ordem:** FREE → ZAI → Outros
- **Prioridade:** P2

---

## Won't Have (Fora do Escopo)

### RF-012: Interface Web
- Fora do escopo atual - apenas CLI/TUI

### RF-013: Benchmark de Performance
- Não mede qualidade ou velocidade de inferência

### RF-014: Gerenciamento de API Keys
- Delega ao OpenCode CLI
