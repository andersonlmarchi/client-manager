# Client Manager

Identidade, organizações, produtos, licenciamento, billing com medição de uso, auditoria e integrações.

Arquitetura em microsserviços Go, PostgreSQL, portal React (Vite) atrás de nginx, orquestrados com Docker Compose.

## Estrutura

Hoje no repositório:

```text
packages/       # shared, events
docker/         # Compose + Postgres
docs/adr/       # decisões de arquitetura
```

Pastas como `services/`, `web/` e `sdk/` entram quando a etapa correspondente criar código real (sem diretórios vazios).

## Desenvolvimento

- Go 1.22+ (um módulo na raiz: `github.com/andersonlmarchi/client-manager`)
- Docker Compose sobe a stack localmente (neste ambiente use `sudo` se necessário)
- Commits e PRs são feitos manualmente ao fechar cada etapa

```bash
export PATH=/usr/local/go/bin:$PATH
go test ./...
```

## Postgres (Docker)

Instruções completas: [`docker/README.md`](docker/README.md).

```bash
cp .env.example .env
# defina POSTGRES_PASSWORD

sudo docker compose -f docker/docker-compose.yml --env-file .env up -d
```

## Documentação de integração

Guia para integrar outros sistemas: [`docs/integration/`](docs/integration/) (preenchido nas etapas de docs).
