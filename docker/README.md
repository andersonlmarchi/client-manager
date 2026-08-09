# Docker (Client Manager)

## Pré-requisitos

- Docker Engine + Compose plugin
- Neste host, Docker costuma exigir `sudo`

## Subir só o Postgres

O banco sobe **vazio** (database/usuário via `POSTGRES_*`). Não há SQL de bootstrap criando schemas de domínio.

Na raiz do repositório:

```bash
cp .env.example .env
# edite POSTGRES_PASSWORD com um segredo forte

sudo docker compose -f docker/docker-compose.yml --env-file .env up -d
```

Healthcheck:

```bash
sudo docker compose -f docker/docker-compose.yml --env-file .env ps
```

Parar:

```bash
sudo docker compose -f docker/docker-compose.yml --env-file .env down
```

Apagar volume (apaga dados):

```bash
sudo docker compose -f docker/docker-compose.yml --env-file .env down -v
```

## Validar arquivo Compose (sem subir)

```bash
sudo docker compose -f docker/docker-compose.yml --env-file .env config
```

## Migrações (sob demanda, automáticas no container)

- Schemas/tabelas são criados **só quando o serviço precisar**, via migrações versionadas daquele serviço (fluxo ent/Atlas).
- Não criar a infra de todos os domínios no início.
- Cada serviço Go usa o entrypoint [`bin/service-entrypoint.sh`](bin/service-entrypoint.sh):
  1. espera o Postgres (`WAIT_HOST`/`WAIT_PORT`)
  2. roda `MIGRATE_CMD` (default: `/app/migrate up`) com `DATABASE_URL`
  3. inicia o binário do serviço (`exec`)

Exemplo futuro no Compose (quando o serviço existir):

```yaml
configuration:
  depends_on:
    postgres:
      condition: service_healthy
  environment:
    DATABASE_URL: postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable
  entrypoint: ["/entrypoint.sh"]
  command: ["/app/configuration"]
```

No Dockerfile do serviço: copiar `docker/bin/service-entrypoint.sh` como `/entrypoint.sh` e publicar um binário `/app/migrate` (ou ajustar `MIGRATE_CMD`).

`SKIP_MIGRATE=true` pula o passo 2 (útil em jobs pontuais).
