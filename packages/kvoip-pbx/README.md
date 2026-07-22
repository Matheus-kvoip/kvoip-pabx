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
  media/           RTP bridge + SDP rewrite
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
| `SIP_AUTH_ENABLED` | `true` | Exige Digest no REGISTER |
| `SIP_AUTH_REALM` | `kvoip.local` | Realm Digest |
| `SIP_USERS` | `1001:kvoip123,1002:kvoip123` | user:senha |
| `MEDIA_ENABLED` | `true` | Relay RTP (reescrita SDP) |
| `MEDIA_BIND_HOST` | `0.0.0.0` | Bind dos sockets RTP |
| `MEDIA_ADVERTISE_HOST` | `SIP_ADVERTISED_HOST` | IP anunciado no SDP |
| `RTP_PORT_MIN` / `RTP_PORT_MAX` | `10000` / `20000` | Faixa de portas RTP |
| `CDR_WEBHOOK_URL` | Nest `/api/calls/cdr` | Envia CDR no hangup |
| `CDR_WEBHOOK_SECRET` | `kvoip-cdr-dev` | Header `X-CDR-Secret` |

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
- **Digest auth** no `REGISTER` (`401` + `WWW-Authenticate`)
- **RTP relay** com reescrita de SDP (`MEDIA_ENABLED`)
- **CDR webhook** para o Nest no hangup
- Destino inexistente → `404 Not Found`

## Testar ligação 1001 → 1002

1. Suba o PBX (`make run`)
2. Registre dois softphones (MicroSIP) com senha:
   - `1001` / `kvoip123`
   - `1002` / `kvoip123`
   - Domain/Realm: `kvoip.local` (ou `127.0.0.1`)
   - Porta local do softphone ≠ 5060
3. Do `1001`, disque `1002`
4. No painel (`npm run dev:server` + `dev:front`), veja ramais online e chamadas ao vivo

## API HTTP do PBX

| Rota | Descrição |
|------|-----------|
| `GET /health` | Saúde + contadores |
| `GET /v1/registrations` | Ramais registrados |
| `GET /v1/calls?active=true` | Chamadas |

## Próximos passos sugeridos

1. Persistência (Postgres) de ramais/CDR
2. Gravação de chamadas / DTMF
3. Digest também no INVITE (opcional)
4. TCP/TLS transports
