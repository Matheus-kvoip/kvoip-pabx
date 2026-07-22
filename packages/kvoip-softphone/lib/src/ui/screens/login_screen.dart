import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../state/softphone_controller.dart';
import '../../services/sip_engine.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key, required this.onRegistered});

  final VoidCallback onRegistered;

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _extCtrl = TextEditingController();
  final _passCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      final c = context.read<SoftphoneController>();
      await c.loadSaved();
      if (!mounted) return;
      _extCtrl.text = c.extension;
      _passCtrl.text = c.password;
    });
  }

  @override
  void dispose() {
    _extCtrl.dispose();
    _passCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final c = context.watch<SoftphoneController>();
    return Scaffold(
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.fromLTRB(24, 32, 24, 24),
          children: [
            Center(
              child: Column(
                children: [
                  Image.asset('assets/logo-kvoip.png', width: 96, height: 96),
                  const SizedBox(height: 12),
                  Text(
                    'KVOIP',
                    style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                          fontWeight: FontWeight.w800,
                          color: const Color(0xFF0080FF),
                          letterSpacing: 1.2,
                        ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    'Softphone',
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: const Color(0xFF5B6B7C),
                        ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 36),
            TextField(
              controller: _extCtrl,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: 'Ramal',
                hintText: '1001',
              ),
              onChanged: (v) => c.extension = v,
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _passCtrl,
              obscureText: true,
              decoration: const InputDecoration(
                labelText: 'Senha SIP',
                hintText: 'kvoip123',
              ),
              onChanged: (v) => c.password = v,
            ),
            const SizedBox(height: 8),
            Text(
              'Servidor: ${c.config.sipHost}:${c.config.sipPort} · ${c.config.sipRealm}',
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: const Color(0xFF5B6B7C),
                  ),
            ),
            if (c.statusMessage != null) ...[
              const SizedBox(height: 12),
              Text(
                c.statusMessage!,
                style: TextStyle(
                  color: c.registration == SipRegistrationState.failed
                      ? Colors.red.shade700
                      : const Color(0xFF5B6B7C),
                ),
              ),
            ],
            const SizedBox(height: 24),
            FilledButton(
              onPressed: c.busy
                  ? null
                  : () async {
                      await c.register();
                      if (c.isRegistered && context.mounted) {
                        widget.onRegistered();
                      }
                    },
              child: Text(c.busy ? 'Registrando…' : 'Registrar no PBX'),
            ),
            const SizedBox(height: 12),
            Text(
              'MVP: registro SIP ainda é mock. Em seguida entra PJSIP/Linphone.',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: const Color(0xFF5B6B7C),
                  ),
            ),
          ],
        ),
      ),
    );
  }
}
