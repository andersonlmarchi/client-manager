# ADR 0002: Pacote shared (primitivos cross-service)

## Status

Aceito

## Contexto

Serviços do Client Manager precisam de erros tipados, UUID, paginação, logging e respostas HTTP padronizadas sem regra de negócio.

## Decisão

- Pacote `packages/shared` no módulo único `github.com/andersonlmarchi/client-manager`.
- Erros com códigos estáveis (`invalid`, `not_found`, `conflict`, `forbidden`, `internal`).
- UUID v4 via stdlib (`crypto/rand`), sem dependência externa.
- Paginação offset com defaults e teto (`DefaultPageSize` / `MaxPageSize`).
- Logger JSON com `log/slog`.
- Respostas de erro em RFC 7807 (`application/problem+json`).

## Consequências

Handlers e application layers reutilizam esses primitivos. Nenhuma entidade de domínio entra em `shared`.
