# 🧪 OpenCode Check

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Nota**: Este projeto **não** é oficialmente construído ou afiliado à equipe do OpenCode (Anomaly). É uma ferramenta da comunidade para diagnóstico de disponibilidade de modelos LLM.

Uma ferramenta CLI com TUI interativa que testa e classifica automaticamente a disponibilidade de modelos LLM através do OpenCode CLI. Obtenha uma visão completa de quais modelos estão acessíveis, gratuitos, com limite de taxa ou que requerem autenticação em menos de 60 segundos.

**🇬🇧 English Version:** [README.md](./README.md)

## ✨ Funcionalidades

- 🔍 **Descoberta Automática** - Encontra todos os modelos configurados via `opencode models`
- ⚡ **Teste Paralelo** - Testa múltiplos modelos simultaneamente com workers configuráveis
- 📊 **Classificação Inteligente** - Categoriza modelos em 12 estados distintos (FREE, PAID, TIMEOUT, etc.)
- 🎨 **TUI em Tempo Real** - Interface de terminal com barras de progresso e atualizações ao vivo
- 💾 **Cache Inteligente** - Reutiliza resultados por 24 horas para acelerar execuções subsequentes
- 🧠 **Base de Conhecimento Extensível** - Personalize classificações de modelos via config JSON

## 📋 Pré-requisitos

- **OpenCode CLI** instalado e configurado ([opencode.ai/docs](https://opencode.ai/docs))
- **Go 1.24+** (para build a partir do código fonte)
- Um emulador de terminal moderno (WezTerm, Alacritty, Ghostty, Kitty, etc.)

## 🚀 Instalação

### Opção 1: Download do Binário

Baixe a versão mais recente em [Releases](https://github.com/your-username/opencode-check/releases) e adicione ao seu PATH.

### Opção 2: Build do Código Fonte

```bash
git clone https://github.com/your-username/opencode-check.git
cd opencode-check
go build -o opencode-check
sudo mv opencode-check /usr/local/bin/
```

### Opção 3: Instalar com Go

```bash
go install github.com/your-username/opencode-check@latest
```

## 📖 Uso

### Uso Básico

```bash
# Testar todos os modelos disponíveis
opencode-check

# Usar cache para acelerar execuções subsequentes
opencode-check --cache

# Atualizar lista de modelos antes de testar
opencode-check --refresh
```

### Opções Avançadas

```bash
# Personalizar número de workers paralelos (padrão: 5)
opencode-check -c 10

# Ajustar timeout por modelo (padrão: 20s)
opencode-check -t 30s

# Usar base de conhecimento customizada
opencode-check --kb custom-kb.json

# Mostrar versão
opencode-check --version
```

### Referência de Flags

| Flag | Padrão | Descrição |
|------|--------|-----------|
| `-c` | `5` | Número de workers paralelos |
| `-t` | `20s` | Timeout por modelo |
| `--cache` | `false` | Usar resultados em cache (válido por 24h) |
| `--refresh` | `false` | Atualizar lista de modelos antes de testar |
| `--kb` | `""` | Caminho para JSON de base de conhecimento customizada |
| `--version` | - | Mostrar informação de versão |

## 📊 Categorias de Modelos

Os resultados são classificados em 12 categorias:

| Ícone | Categoria | Significado |
|-------|-----------|-------------|
| 🆓 | `FREE` | Modelos gratuitos sem limites conhecidos |
| 📊 | `FREE_LIMITED` | Gratuitos com quotas (ex: Cerebras, DeepSeek, Groq) |
| 💰 | `PAID` | Modelos ZAI pagos com créditos ativos |
| ✅ | `AVAILABLE` | Disponível (geral) |
| ❓ | `NOT_FOUND` | Modelo não existe |
| ⏰ | `TIMEOUT` | Timeout (padrão 20s) |
| 🔒 | `AUTH_FAILED` | Chave de API inválida |
| ❌ | `NO_QUOTA` | Sem créditos restantes |
| ⏱️ | `RATE_LIMITED` | Limite de taxa atingido |
| ⚠️ | `ERROR` | Erro desconhecido |

## 🔧 Configuração

### Base de Conhecimento Customizada

Crie um arquivo JSON para sobrescrever as classificações padrão:

```json
{
  "free_models": {
    "opencode/meu-modelo": {
      "category": "FREE",
      "description": "Meu Modelo",
      "limits": "sem limites documentados"
    }
  },
  "free_tier_providers": {
    "meuprovider": {
      "category": "FREE_LIMITED",
      "description": "Meu Provider",
      "limits": "1M tokens/dia"
    }
  }
}
```

Use com:
```bash
opencode-check --kb custom-kb.json
```

## 📁 Arquivos de Saída

Resultados são salvos em:
- **Cache**: `~/.config/opencode/cache/results.json` (ao usar `--cache`)
- **Relatórios**: `~/.config/opencode/results/opencode-check-YYYYMMDD-HHMMSS.json` (pressione `s` para salvar)

## 🤝 Compatibilidade com Plugins OpenCode

Esta ferramenta funciona junto com plugins populares do OpenCode:

- **[opencode-antigravity-auth](https://github.com/NoeFabris/opencode-antigravity-auth)** - Testa modelos OAuth Antigravity
- **[oh-my-opencode](https://github.com/code-yeongyu/oh-my-opencode)** - Compatível com recursos de agente aprimorados

## 🐛 Troubleshooting

### "falha ao descobrir modelos"
- Certifique-se de que o OpenCode CLI está instalado: `opencode --version`
- Verifique se você configurou pelo menos um provedor: `opencode models`

### Modelos aparecem como "NOT_FOUND"
- Execute `opencode models --refresh` para atualizar a lista
- Alguns modelos podem ter sido descontinuados ou renomeados

### Alta taxa de timeouts
- Aumente o timeout: `opencode-check -t 30s`
- Verifique sua conexão com a internet
- Alguns provedores podem estar temporariamente indisponíveis

## 📜 Licença

Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- Construído para a comunidade [OpenCode](https://opencode.ai)
- Usa [Bubble Tea](https://github.com/charmbracelet/bubbletea) para TUI
- Inspirado pelo incrível ecossistema de plugins do OpenCode

---

**Aviso Legal**: Esta é uma ferramenta independente da comunidade e não é oficialmente mantida ou afiliada à Anomaly (criadores do OpenCode).
