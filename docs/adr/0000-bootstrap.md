# ADR 0000: Bootstrap do monorepo

## Status

Aceito

## Contexto

Client Manager começa greenfield, em etapas miúdas, com microsserviços Go e packages compartilhados.

## Decisão

- Nome do software: **Client Manager** (slug `client-manager`).
- Monorepo com `go.work` apontando para `packages/*` (serviços entram nas etapas seguintes).
- Module path base: `github.com/t-code/client-manager/...`
- Licença: MIT (`LICENSE`).

## Consequências

Serviços e `sdk/` serão adicionados ao `go.work` quando ganharem `go.mod` próprio.
