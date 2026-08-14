# ADR 0006: Configuration REST API

## Status

Aceito

## Contexto

Outros serviços e o admin precisam ler e atualizar políticas tipadas sem acesso direto ao schema `configuration`.

## Decisão

- Endpoints JSON em `/v1`:
  - `GET|PUT /v1/password-policy`
  - `GET|PUT /v1/oidc-settings`
  - `GET|PUT /v1/rate-limits`
  - `GET|PUT /v1/branding`
- Leitura aberta na rede interna; escrita exige `ADMIN_API_KEY` via `X-Admin-Key` ou `Authorization: Bearer`.
- Application service em `internal/application`; handlers em `internal/transport/http`.
- Corpo limitado (64 KiB), `DisallowUnknownFields`, erros em `application/problem+json`.
- `ADMIN_API_KEY` obrigatório no processo e no Compose.

## Consequências

Identity e demais serviços consomem GET tipado. Admin UI (portal) usará PUT autenticado. Sem JWT de usuário nesta etapa.
