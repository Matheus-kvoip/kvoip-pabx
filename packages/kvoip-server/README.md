# Database (kvoip-server)

## Local (padrão)

`DB_TYPE=sqlite` — arquivo em `data/kvoip.sqlite`, `DB_SYNC=true`.

## Postgres (produção / compose)

```env
DB_TYPE=postgres
DATABASE_URL=postgres://kvoip:kvoip@127.0.0.1:5432/kvoip
DB_SYNC=false
DB_MIGRATIONS_RUN=true
```

Subir banco:

```bash
# com Docker Desktop
docker compose up -d db
```

Rodar migrations manualmente:

```bash
cd packages/kvoip-server
npm run migration:run
```

## CDR

O PBX envia `POST /api/calls/cdr` no hangup com header `X-CDR-Secret`.
Segredo compartilhado: `CDR_WEBHOOK_SECRET` (server + pbx).
