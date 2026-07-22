/// Configuração do softphone KVOIP.
class AppConfig {
  const AppConfig({
    required this.apiUrl,
    required this.sipHost,
    required this.sipPort,
    required this.sipRealm,
  });

  final String apiUrl;
  final String sipHost;
  final int sipPort;
  final String sipRealm;

  /// Defaults de desenvolvimento (máquina local).
  static const AppConfig development = AppConfig(
    apiUrl: String.fromEnvironment(
      'API_URL',
      defaultValue: 'http://10.0.2.2:3001/api', // emulador Android → host
    ),
    sipHost: String.fromEnvironment('SIP_HOST', defaultValue: '10.0.2.2'),
    sipPort: int.fromEnvironment('SIP_PORT', defaultValue: 5060),
    sipRealm: String.fromEnvironment('SIP_REALM', defaultValue: 'kvoip.local'),
  );
}
