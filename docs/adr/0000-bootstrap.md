# ADR 0000: Bootstrap do monorepo

## Status

Aceito

## Contexto

Platform Core começa greenfield, em etapas miúdas, com microsserviços Go e packages compartilhados.

## Decisão

- Monorepo com `go.work` apontando para `packages/*` (serviços entram nas etapas seguintes).
- Module path base: `github.com/t-code/platform-core/...`
- Licença: MIT (`LICENSE`).

## Consequências

Serviços e `sdk/` serão adicionados ao `go.work` quando ganharem `go.mod` próprio.
