# ADR 0001: Docker Compose + Postgres + migrações sob demanda

## Status

Aceito

## Contexto

Client Manager roda em VPS com Docker. O banco não deve nascer com todos os schemas de domínio. Migrações precisam ser automatizáveis na subida dos containers.

## Decisão

- Orquestração: Docker Compose (`docker/docker-compose.yml`).
- PostgreSQL 16 Alpine sobe com database/usuário vazios (`POSTGRES_*`); **sem** `docker-entrypoint-initdb.d` de domínio.
- Credenciais só por variáveis de ambiente (`.env` local; `.env.example` no repo).
- `POSTGRES_PASSWORD` obrigatório (Compose falha se ausente).
- Schema/tabela por serviço, **quando a etapa do serviço precisar**, via migrações versionadas (ent + Atlas).
- Automação: `docker/bin/service-entrypoint.sh` espera o Postgres, aplica `MIGRATE_CMD`, depois `exec` do processo.
- Agente de desenvolvimento não sobe containers; o operador usa `sudo` quando necessário.

## Consequências

- Volume `postgres_data` persiste dados; não há init SQL de schemas no first boot.
- Cada serviço é responsável pelo próprio schema e por falhar o boot se migrate falhar.
- nginx e demais serviços entram no mesmo Compose nas etapas seguintes.
