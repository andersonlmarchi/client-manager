# ADR 0004: Serviço configuration (skeleton)

## Status

Aceito

## Contexto

Configuration é o primeiro microsserviço de domínio. A E05 só precisa provar build, health HTTP e integração no Compose; persistência entra na etapa seguinte.

## Decisão

- Serviço em `services/configuration` (transport HTTP + `cmd/configuration`).
- Endpoint `GET /health` retorna JSON `{"status":"ok"}` sem dados sensíveis.
- Imagem multi-stage (Go build + Alpine), usuário não-root, entrypoint compartilhado `docker/bin/service-entrypoint.sh`.
- Compose: serviço `configuration` depende do Postgres healthy; `SKIP_MIGRATE=true` até existir binário/migrate (E06+).
- Porta host padrão `8081` → container `8080`.

## Consequências

Próximas etapas adicionam schema ent, migrate no boot e API de políticas. O contrato de health permanece estável para healthchecks.
