# Arquitetura do kvoip-pbx

Fluxo alvo:

```
UA / Softphone
    │ SIP (UDP/TCP/TLS)
    ▼
transport → server → sip.Parse → handlers
                              ├─ dialog
                              ├─ session
                              ├─ proxy/location
                              └─ routing/dialplan
                                    │
                                    ▼
                              events → (futuro) kvoip-server
```

MVP: UDP + parse + stubs de handler.
