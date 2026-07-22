import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../services/sip_engine.dart';
import '../../state/softphone_controller.dart';

class DialerScreen extends StatelessWidget {
  const DialerScreen({super.key, required this.onOpenCall});

  final VoidCallback onOpenCall;

  static const _keys = [
    ['1', '2', '3'],
    ['4', '5', '6'],
    ['7', '8', '9'],
    ['*', '0', '#'],
  ];

  @override
  Widget build(BuildContext context) {
    final c = context.watch<SoftphoneController>();
    return Scaffold(
      appBar: AppBar(
        title: const Text('Discador'),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 12),
            child: Center(
              child: Text(
                c.isRegistered ? 'Online' : 'Offline',
                style: TextStyle(
                  color: c.isRegistered ? Colors.green.shade700 : Colors.orange,
                  fontWeight: FontWeight.w700,
                  fontSize: 13,
                ),
              ),
            ),
          ),
        ],
      ),
      body: SafeArea(
        child: Column(
          children: [
            const SizedBox(height: 24),
            Text(
              c.digits.isEmpty ? ' ' : c.digits,
              style: Theme.of(context).textTheme.headlineLarge?.copyWith(
                    letterSpacing: 2,
                    fontWeight: FontWeight.w600,
                  ),
            ),
            Text(
              'Ramal ${c.extension}',
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: const Color(0xFF5B6B7C),
                  ),
            ),
            const Spacer(),
            for (final row in _keys)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 6),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    for (final d in row)
                      _DialKey(label: d, onTap: () => c.tapDigit(d)),
                  ],
                ),
              ),
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                IconButton(
                  onPressed: c.clearDigits,
                  icon: const Icon(Icons.clear),
                  tooltip: 'Limpar',
                ),
                FilledButton(
                  onPressed: !c.isRegistered || c.digits.isEmpty
                      ? null
                      : () async {
                          await c.placeCall();
                          onOpenCall();
                        },
                  style: FilledButton.styleFrom(
                    shape: const CircleBorder(),
                    padding: const EdgeInsets.all(22),
                    backgroundColor: Colors.green.shade600,
                  ),
                  child: const Icon(Icons.call, size: 28),
                ),
                IconButton(
                  onPressed: c.backspace,
                  icon: const Icon(Icons.backspace_outlined),
                  tooltip: 'Apagar',
                ),
              ],
            ),
            const SizedBox(height: 24),
          ],
        ),
      ),
    );
  }
}

class _DialKey extends StatelessWidget {
  const _DialKey({required this.label, required this.onTap});

  final String label;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 78,
      height: 78,
      child: Material(
        color: Colors.white,
        shape: const CircleBorder(),
        elevation: 0.5,
        child: InkWell(
          customBorder: const CircleBorder(),
          onTap: onTap,
          child: Center(
            child: Text(
              label,
              style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                    fontWeight: FontWeight.w600,
                  ),
            ),
          ),
        ),
      ),
    );
  }
}

class CallScreen extends StatelessWidget {
  const CallScreen({super.key, required this.onEnded});

  final VoidCallback onEnded;

  @override
  Widget build(BuildContext context) {
    final c = context.watch<SoftphoneController>();
    if (c.callState == SipCallUiState.idle ||
        c.callState == SipCallUiState.ended) {
      WidgetsBinding.instance.addPostFrameCallback((_) => onEnded());
    }

    final label = switch (c.callState) {
      SipCallUiState.dialing => 'Discando…',
      SipCallUiState.ringing => 'Chamando…',
      SipCallUiState.active => 'Em chamada',
      SipCallUiState.ended => 'Encerrada',
      SipCallUiState.idle => 'Livre',
    };

    return Scaffold(
      backgroundColor: const Color(0xFF0B1B2B),
      body: SafeArea(
        child: Column(
          children: [
            const Spacer(),
            Text(
              c.peer ?? c.digits,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 36,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 8),
            Text(label, style: const TextStyle(color: Colors.white70)),
            const Spacer(),
            FilledButton(
              onPressed: () async {
                await c.hangup();
                onEnded();
              },
              style: FilledButton.styleFrom(
                backgroundColor: Colors.red.shade600,
                shape: const CircleBorder(),
                padding: const EdgeInsets.all(24),
              ),
              child: const Icon(Icons.call_end, size: 32),
            ),
            const SizedBox(height: 48),
          ],
        ),
      ),
    );
  }
}
