# kvoip-pbx

Núcleo SIP / PABX da KVOIP (Go).

## Layout

```
cmd/kvoip-pbx/     entrypoint
configs/           config.yaml de referência
internal/
  config/          env + defaults
  server/          listeners (UDP primeiro)
  sip/             parse mínimo SIP
  handlers/        OPTIONS/REGISTER/INVITE/...
  dialog/          estado de diálogo
  session/         sessões de chamada
  proxy/           location / router
  routing/         dialplan
  transport/       udp/tcp/tls types
  storage/         persistência (placeholder)
  events/          bus interno
pkg/version/       versão do binário
docs/              notas de arquitetura
```

## Variáveis (`.env`)

| Chave | Default | Descrição |
|-------|---------|-----------|
| `PORT_SERVER_SIP` | `5060` | Porta SIP |
| `SIP_BIND_HOST` | `0.0.0.0` | Bind |
| `SIP_ADVERTISED_HOST` | `127.0.0.1` | Host no Via (proxy) |
| `SIP_BUFFER_SIZE` | `8192` | Buffer UDP |
| `LOG_LEVEL` | `info` | debug/info/warn/error |
| `SERVICE_NAME` | `kvoip-pbx` | Nome no log |

## Pré-requisito

- Go 1.22+

## Comandos

Na pasta `packages/kvoip-pbx`:

```bash
go mod tidy
make run
# ou
go run ./cmd/kvoip-pbx
```

Build:

```bash
make build
./bin/kvoip-pbx
```

## Estado atual

- Listener **UDP** com reply e **SendTo** (proxy)
- `OPTIONS` / `REGISTER` / **INVITE proxy**
- **API HTTP** em `:8080` para o Nest (`/health`, `/v1/registrations`, `/v1/calls`)
- Destino inexistente → `404 Not Found`
- Sem auth Digest / sem RTP ainda (SDP é só retransmitido)

## Testar ligação 1001 → 1002

1. Suba o PBX (`make run`)
2. Registre dois softphones no domínio `kvoip.local` (ou IP do PBX), ramais `1001` e `1002`
3. Do `1001`, disque `1002`
4. No painel (`npm run dev:server` + `dev:front`), veja ramais online e chamadas ao vivo

## API HTTP do PBX

| Rota | Descrição |
|------|-----------|
| `GET /health` | Saúde + contadores |
| `GET /v1/registrations` | Ramais registrados |
| `GET /v1/calls?active=true` | Chamadas |

## Próximos passos sugeridos

1. Auth Digest no REGISTER
2. Persistência (Postgres) de ramais/CDR
3. TCP/TLS transports
