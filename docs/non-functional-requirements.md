# Requisitos Não-Funcionais - OpenCode Check

## Performance

| ID | Requisito | Meta | Método de Validação |
|----|-----------|------|---------------------|
| RNF-001 | Tempo total de execução | < 60s para 50 modelos | Benchmark automatizado |
| RNF-002 | Latência p95 por modelo | < 15s | Logs de timing |
| RNF-003 | Uso de memória | < 100MB | Profiling com pprof |
| RNF-004 | Taxa de refresh da TUI | ≥ 30 FPS | Observação visual |
| RNF-005 | Startup time | < 2s | Medição cronometrada |

---

## Segurança

| ID | Requisito | Implementação |
|----|-----------|---------------|
| RNF-006 | Isolamento de processos | `Setpgid: true` em processos filhos |
| RNF-007 | Cleanup de processos | Kill group no timeout |
| RNF-008 | Sem exposição de API keys | Keys não logadas em output |
| RNF-009 | Permissões de arquivos | Cache/results com modo 0644 |

---

## Escalabilidade

| ID | Requisito | Suporte Atual |
|----|-----------|---------------|
| RNF-010 | Workers paralelos | Configurável via `-c` (1-N) |
| RNF-011 | Número de modelos | Testado com 50+, sem limite hard |
| RNF-012 | Tamanho de output | Truncamento inteligente (64KB default) |
| RNF-013 | Crescimento de provedores | KB extensível via JSON |

---

## Confiabilidade

| ID | Requisito | Meta |
|----|-----------|------|
| RNF-014 | Disponibilidade | N/A (ferramenta CLI) |
| RNF-015 | Recuperação de erros | Continuar após falha individual |
| RNF-016 | Retry automático | 1 retry em rate limit |
| RNF-017 | Graceful shutdown | `q` ou `Ctrl+C` limpa recursos |

---

## Usabilidade

| ID | Requisito | Implementação |
|----|-----------|---------------|
| RNF-018 | Feedback visual | Ícones coloridos por categoria |
| RNF-019 | Cores adaptativas | Light/Dark mode automático |
| RNF-020 | Atalhos de teclado | `q` sair, `s` salvar |
| RNF-021 | Responsividade | Ajuste automático ao terminal |

---

## Compatibilidade

| ID | Requisito | Status |
|----|-----------|--------|
| RNF-022 | Linux x86_64 | ✅ Suportado |
| RNF-023 | Linux arm64 | ✅ Suportado |
| RNF-024 | macOS | ✅ Suportado (syscall adaptado) |
| RNF-025 | Windows | ⚠️ Parcial (sem Setpgid) |
| RNF-026 | Go 1.20+ | ✅ Requer 1.20+ |

---

## Manutenibilidade

| ID | Requisito | Implementação |
|----|-----------|---------------|
| RNF-027 | Modularidade | Separação clara: TUI, workers, cache |
| RNF-028 | Logs estruturados | JSON output via flag |
| RNF-029 | Versionamento | Semver via constante `Version` |
| RNF-030 | Configuração externa | KB via arquivo JSON |
