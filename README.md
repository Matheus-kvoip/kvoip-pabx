# Kvoip PABX

Monorepo da plataforma PABX da Kvoip.

## Estrutura

```
packages/
  kvoip-pbx/       Núcleo SIP (Go)
  kvoip-server/    API NestJS (controle)
  kvoip-front/     Painel Next.js
  kvoip-shared/    Tipos compartilhados
  kvoip-softphone/ Softphone (em breve)
```

## Pré-requisitos

- Node.js 20+
- npm 10+
- Go 1.22+ (para o núcleo SIP)

## Variáveis de ambiente

Copie os exemplos (já há `.env` / `.env.local` para desenvolvimento):

| Pacote | Arquivo | Principais chaves |
|--------|---------|-------------------|
| `kvoip-server` | `.env` | `PORT`, `CORS_ORIGIN`, `JWT_*`, `AUTH_*` |
| `kvoip-front` | `.env.local` | `PORT`, `NEXT_PUBLIC_*`, `AUTH_COOKIE` |
| `kvoip-pbx` | `.env` | `PORT_SERVER_SIP`, `SIP_BIND_HOST`, `LOG_LEVEL` |

Login demo (server):

- e-mail: `admin@kvoip.com.br`
- senha: `kvoip123`

## Desenvolvimento

```bash
npm install
npm run build -w @kvoip/shared
npm run dev:server
npm run dev:front
npm run dev:pbx
```

- Login: http://localhost:3000/login
- API: http://localhost:3001/api/health
- SIP UDP: `0.0.0.0:5060` (`packages/kvoip-pbx`)

Detalhes do PBX: `packages/kvoip-pbx/README.md`

## API

| Método | Rota | Auth | Descrição |
|--------|------|------|-----------|
| GET | `/api/health` | pública | Saúde |
| POST | `/api/auth/login` | pública | Login JWT |
| GET | `/api/auth/me` | JWT | Usuário atual |
| POST | `/api/auth/logout` | JWT | Logout |
| GET | `/api/dashboard` | JWT | Métricas |
| * | `/api/extensions` | JWT | Ramais |
| GET | `/api/trunks` | JWT | Troncos |
| GET | `/api/calls` | JWT | Chamadas |

## Docker

```bash
docker compose up --build
```

## Observação

Os dados da API ainda são **em memória** (seed de demonstração), sem banco. O núcleo SIP em Go continua separado e ainda não está integrado.
