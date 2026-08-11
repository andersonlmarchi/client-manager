# ADR 0003: Pacote events (envelope e outbox)

## Status

Aceito

## Contexto

Microsserviços não compartilham bus in-process. Efeitos colaterais (audit, webhooks, etc.) precisam de um contrato comum de evento e de registro durável (outbox).

## Decisão

- Pacote `packages/events` com `DomainEvent` (envelope JSON) e `OutboxRecord`.
- Produtores serializam o evento e gravam outbox; consumidores leem outbox e reconstroem o envelope.
- Status do outbox: `pending`, `processed`, `failed` (com contagem de tentativas e último erro).
- Persistência concreta (tabelas/schema) fica para o serviço que publicar ou consumir; este pacote só define o contrato e helpers de marshal.

## Consequências

Audit/integrations/worker dependem deste contrato. Schema `events` e migrations surgem quando o primeiro publisher/consumer for implementado.
