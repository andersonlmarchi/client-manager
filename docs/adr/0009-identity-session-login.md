# ADR 0009: Identity Session e login/logout

## Status

Aceito

## Contexto

O portal e clientes precisam de sessão autenticada após validar email/senha, com revogação no logout.

## Decisão

- Entidade ent `Session`: `token_hash` (SHA-256 do token opaco), `expires_at`, `revoked_at`, IP/UA opcionais.
- Token opaco só viaja no cookie `cm_session` (HttpOnly, SameSite=Lax) e/ou no JSON de login / Bearer.
- `AuthService`: Login, Logout, CurrentUser; falhas de credencial sempre como `unauthorized` (sem enumerar usuário).
- Endpoints: `POST /v1/login`, `POST /v1/logout`, `GET /v1/me`.
- TTL configurável por `SESSION_TTL_HOURS` (default 168h); `COOKIE_SECURE` para HTTPS.

## Consequências

Device/TrustedDevice (próxima etapa) associa sessões a dispositivos sem mudar o contrato de login.
