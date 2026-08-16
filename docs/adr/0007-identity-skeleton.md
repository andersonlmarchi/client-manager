# ADR 0007: Serviço identity (skeleton)

## Status

Aceito

## Contexto

Identity é o serviço de autenticação. A E08 só precisa provar build, health HTTP e entrada no Compose; usuários/credenciais entram na etapa seguinte.

## Decisão

- Serviço em `services/identity` (`cmd/identity` + transport HTTP).
- Endpoint `GET /health` retorna `{"status":"ok"}`.
- Imagem multi-stage (Go build + Alpine), usuário não-root, entrypoint compartilhado.
- Compose: serviço `identity` na porta host `8082`; `SKIP_MIGRATE=true` até existir migrate (E09+).
- `DATABASE_URL` já aponta `search_path=identity` para as próximas etapas.

## Consequências

E09 adiciona User/Credential (Argon2id), ent e migrate no boot.
