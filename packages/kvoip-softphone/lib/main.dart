import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import 'src/config/app_config.dart';
import 'src/services/sip_engine.dart';
import 'src/state/softphone_controller.dart';
import 'src/ui/screens/dialer_screen.dart';
import 'src/ui/screens/login_screen.dart';
import 'src/ui/theme/kvoip_theme.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const KvoipSoftphoneApp());
}

class KvoipSoftphoneApp extends StatefulWidget {
  const KvoipSoftphoneApp({super.key});

  @override
  State<KvoipSoftphoneApp> createState() => _KvoipSoftphoneAppState();
}

class _KvoipSoftphoneAppState extends State<KvoipSoftphoneApp> {
  late final SoftphoneController _controller;
  late final MockSipEngine _sip;
  var _route = _AppRoute.login;

  @override
  void initState() {
    super.initState();
    _sip = MockSipEngine();
    _controller = SoftphoneController(
      config: AppConfig.development,
      sip: _sip,
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    _sip.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider.value(
      value: _controller,
      child: MaterialApp(
        title: 'KVOIP Softphone',
        debugShowCheckedModeBanner: false,
        theme: KvoipTheme.light(),
        home: switch (_route) {
          _AppRoute.login => LoginScreen(
              onRegistered: () => setState(() => _route = _AppRoute.dialer),
            ),
          _AppRoute.dialer => DialerScreen(
              onOpenCall: () => setState(() => _route = _AppRoute.call),
            ),
          _AppRoute.call => CallScreen(
              onEnded: () => setState(() => _route = _AppRoute.dialer),
            ),
        },
      ),
    );
  }
}

enum _AppRoute { login, dialer, call }
