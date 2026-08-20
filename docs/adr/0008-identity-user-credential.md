# ADR 0008: Identity User e Credential (Argon2id)

## Status

Aceito

## Contexto

Identity precisa persistir usuários e credenciais com hash seguro antes dos fluxos de sessão/login.

## Decisão

- Entidades ent: `User` (email único, status) e `Credential` (hash sensível, algorithm).
- Hash **Argon2id** (PHC string) em `domain.HashPassword` / `VerifyPassword`; senha em claro nunca é persistida.
- `UserRepository.CreateUserWithPassword` cria user+credential em transação; `Authenticate` valida status e hash.
- Schema Postgres `identity` via migrate no boot (`cmd/migrate`); DSN com `search_path=identity`.
- Transport HTTP permanece com `Server` + `/health` (login/sessão na próxima etapa).

## Consequências

E10 adiciona Session e endpoints de login/logout reutilizando o repositório e o `Server` existentes.
