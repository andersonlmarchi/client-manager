# ADR 0005: Configuration policies com ent

## Status

Aceito

## Contexto

Configuration precisa persistir políticas globais (password, OIDC, rate limits, branding/SMTP refs) com ORM tipado e migrate automático no boot do container.

## Decisão

- ORM **ent** em `services/configuration/ent` com entidades: `PasswordPolicy`, `OIDCSettings`, `RateLimitSettings`, `BrandingSettings`.
- Domínio em `internal/domain` com validação; repositório em `internal/infrastructure`.
- Schema Postgres `configuration` criado no migrate; DSN usa `search_path=configuration`.
- Binário `cmd/migrate` (`migrate up`) cria schema/tabelas e faz upsert dos defaults (`id=default`).
- Testes de repositório com SQLite em memória (`modernc.org/sqlite`), sem depender de Docker no `go test`.
- SMTP armazena apenas **refs** de secret (não credenciais em claro).

## Consequências

API REST (E07) lê/escreve via repositório. Outros serviços consultam configuration por HTTP, nunca o schema diretamente.
