# ADR 0000: Bootstrap do monorepo

## Status

Aceito

## Contexto

Client Manager começa greenfield, em etapas miúdas, com microsserviços Go e packages compartilhados no mesmo repositório.

## Decisão

- Nome do software: **Client Manager** (slug `client-manager`).
- Um único módulo Go na raiz: `github.com/andersonlmarchi/client-manager` (remoto https://github.com/andersonlmarchi/client-manager).
- Monorepo permanente: packages e services são pastas do mesmo módulo, sem multi-module/`go.work`.
- Licença: MIT (`LICENSE`).

## Consequências

Novos serviços entram como `services/<nome>/...` no mesmo `go.mod`. Imports usam o path do módulo raiz.
