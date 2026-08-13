# Docker (Client Manager)

## Pré-requisitos

- Docker Engine + Compose plugin
- Neste host, Docker costuma exigir `sudo`

## Up / Down

Na raiz do repositório:

```bash
cp .env.example .env
# edite POSTGRES_PASSWORD com um segredo forte

sudo docker compose -f docker/docker-compose.yml --env-file .env up -d --build
sudo docker compose -f docker/docker-compose.yml --env-file .env down
```

Serviços atuais: `postgres`, `configuration` (`GET http://localhost:8081/health`).

Status:

```bash
sudo docker compose -f docker/docker-compose.yml --env-file .env ps
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

- O entrypoint [`bin/service-entrypoint.sh`](bin/service-entrypoint.sh) espera o Postgres e roda `/app/migrate up`.
- O migrate do `configuration` cria o schema `configuration`, aplica tabelas ent e grava defaults das políticas.
