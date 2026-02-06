# Escopo do Projeto - OpenCode Check

## Visão Geral

**OpenCode Check** é uma ferramenta CLI escrita em Go que automatiza o teste de disponibilidade de modelos LLM através do OpenCode CLI. A aplicação fornece uma interface TUI interativa que permite aos usuários visualizar em tempo real quais modelos estão acessíveis e suas respectivas classificações de disponibilidade.

## Contexto de Negócio

### Problema

Usuários do OpenCode precisam identificar rapidamente:
- Quais modelos LLM estão disponíveis para uso
- Quais modelos são gratuitos vs. pagos
- Quais modelos têm limitações de quota/rate limit
- Status de autenticação e créditos por provedor

### Solução

Uma ferramenta que:
1. Descobre automaticamente todos os modelos configurados
2. Testa cada modelo com uma requisição mínima
3. Classifica os resultados em categorias claras
4. Exibe progresso em tempo real via TUI
5. Persiste resultados em cache para consultas rápidas

## Stakeholders

| Stakeholder | Interesse |
|-------------|-----------|
| Desenvolvedores | Escolher modelos disponíveis para projetos |
| Equipe DevOps | Monitorar status de APIs e créditos |
| Product Managers | Entender custos de LLM disponíveis |

## Escopo

### Incluído

- ✅ Descoberta automática de modelos via `opencode models`
- ✅ Teste de conectividade com timeout configurável
- ✅ Classificação automática (FREE, FREE_LIMITED, PAID, ERROR, etc.)
- ✅ Interface TUI com barra de progresso
- ✅ Sistema de cache com expiração de 24h
- ✅ Base de conhecimento extensível via JSON
- ✅ Suporte a múltiplos workers paralelos

### Excluído

- ❌ Gerenciamento de API keys
- ❌ Billing e cobrança
- ❌ Métricas de performance de modelos
- ❌ Comparação de qualidade de respostas

## Restrições

- Requer OpenCode CLI instalado e configurado
- Dependência de conectividade de rede
- Limitado a modelos registrados no OpenCode
