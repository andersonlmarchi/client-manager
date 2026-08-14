# Client Manager

Identidade, organizações, produtos, licenciamento, billing com medição de uso, auditoria e integrações.

Arquitetura em microsserviços Go, PostgreSQL, portal React (Vite) atrás de nginx, orquestrados com Docker Compose.

## Estrutura

Hoje no repositório:

```text
packages/       # shared, events
services/       # configuration
docker/         # Compose + Postgres
docs/adr/       # decisões de arquitetura
```

Pastas como `web/` e `sdk/` entram quando a etapa correspondente criar código real (sem diretórios vazios).

## Desenvolvimento

- Go 1.25+ (um módulo na raiz: `github.com/andersonlmarchi/client-manager`)
- Docker Compose sobe a stack localmente (neste ambiente use `sudo` se necessário)
- Commits e PRs são feitos manualmente ao fechar cada etapa

```bash
export PATH=/usr/local/go/bin:$PATH
go test ./...
```

## Docker Compose

Instruções detalhadas: [`docker/README.md`](docker/README.md).

```bash
cp .env.example .env
# defina POSTGRES_PASSWORD

sudo docker compose -f docker/docker-compose.yml --env-file .env up -d --build
sudo docker compose -f docker/docker-compose.yml --env-file .env down
```

Health do configuration: `http://localhost:8081/health`

API (leitura): `GET http://localhost:8081/v1/password-policy` (também `oidc-settings`, `rate-limits`, `branding`).

Escrita admin: `PUT` nos mesmos paths com header `X-Admin-Key` (valor de `ADMIN_API_KEY`).

Apagar volumes (apaga dados do Postgres):

```bash
sudo docker compose -f docker/docker-compose.yml --env-file .env down -v
```

## Documentação de integração

Guia para integrar outros sistemas: [`docs/integration/`](docs/integration/) (preenchido nas etapas de docs).
