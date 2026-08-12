# Docker (Client Manager)

## Pré-requisitos

- Docker Engine + Compose plugin
- Neste host, Docker costuma exigir `sudo`

## Subir stack

Na raiz do repositório:

```bash
cp .env.example .env
# edite POSTGRES_PASSWORD com um segredo forte

sudo docker compose -f docker/docker-compose.yml --env-file .env up -d --build
```

Serviços atuais: `postgres`, `configuration` (`GET http://localhost:8081/health`).

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

## Migrações

- Schemas/tabelas por serviço, sob demanda (ent/Atlas).
- Entrypoint: [`bin/service-entrypoint.sh`](bin/service-entrypoint.sh).
- `configuration` usa `SKIP_MIGRATE=true` até a etapa de migrate do serviço.
