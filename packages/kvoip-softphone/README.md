# KVOIP Softphone (Flutter)

App mobile do softphone Kvoip — **Android / iOS**.

## Estado atual (início)

- UI: login SIP (ramal/senha) → discador → tela de chamada
- Tema Kvoip (azul `#0080FF` + logo)
- `MockSipEngine` para validar fluxo **sem** SDK SIP ainda
- Próximo passo: trocar o mock por **PJSIP / Linphone SDK**

## Pré-requisito

1. Instale o [Flutter SDK](https://docs.flutter.dev/get-started/install/windows)
2. No PATH: `flutter` e `dart`
3. Aceite licenses Android: `flutter doctor --android-licenses`

## Gerar pastas nativas (primeira vez)

Na pasta deste pacote:

```bash
cd packages/kvoip-softphone
flutter create . --project-name kvoip_softphone --org com.kvoip
flutter pub get
```

> Isso cria `android/`, `ios/`, etc. sem apagar o `lib/` já pronto.

## Rodar

```bash
flutter run
# ou
flutter run -d windows   # smoke na desktop
flutter run -d chrome    # UI only
```

### Emulador Android → PBX na sua máquina

Defaults em `AppConfig.development`:

| Chave | Default | Nota |
|-------|---------|------|
| `SIP_HOST` | `10.0.2.2` | host do PC visto do emulador |
| `SIP_PORT` | `5060` | UDP SIP do `kvoip-pbx` |
| `SIP_REALM` | `kvoip.local` | Digest |
| `API_URL` | `http://10.0.2.2:3001/api` | Nest |

No device físico, use o IP LAN do PC, por exemplo:

```bash
flutter run --dart-define=SIP_HOST=192.168.0.10 --dart-define=API_URL=http://192.168.0.10:3001/api
```

Credenciais seed: ramal `1001` / senha `kvoip123`.

## Estrutura

```
lib/
  main.dart
  src/
    config/app_config.dart
    services/kvoip_api.dart      # Nest (login painel)
    services/sip_engine.dart     # contrato + MockSipEngine
    state/softphone_controller.dart
    ui/theme/kvoip_theme.dart
    ui/screens/login_screen.dart
    ui/screens/dialer_screen.dart
assets/logo-kvoip.png
```

## Roadmap SIP real

1. Integrar Linphone SDK / flutter_pjsip (plugin)
2. REGISTER Digest no PBX Kvoip
3. INVITE + áudio (RTP)
4. Push (FCM/APNs) para chamada entrante em background
