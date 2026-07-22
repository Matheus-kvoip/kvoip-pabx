import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../config/app_config.dart';
import '../services/sip_engine.dart';

class SoftphoneController extends ChangeNotifier {
  SoftphoneController({
    required this.config,
    required SipEngine sip,
  }) : _sip = sip {
    _sub = _sip.events.listen(_onSipEvent);
  }

  final AppConfig config;
  final SipEngine _sip;
  late final StreamSubscription<SipEvent> _sub;

  String extension = '1001';
  String password = 'kvoip123';
  String digits = '';
  String? statusMessage;
  bool busy = false;

  SipRegistrationState registration = SipRegistrationState.idle;
  SipCallUiState callState = SipCallUiState.idle;
  String? peer;

  bool get isRegistered => registration == SipRegistrationState.registered;
  bool get inCall =>
      callState == SipCallUiState.dialing ||
      callState == SipCallUiState.ringing ||
      callState == SipCallUiState.active;

  Future<void> loadSaved() async {
    final prefs = await SharedPreferences.getInstance();
    extension = prefs.getString('sip_extension') ?? extension;
    password = prefs.getString('sip_password') ?? password;
    notifyListeners();
  }

  Future<void> saveCredentials() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('sip_extension', extension);
    await prefs.setString('sip_password', password);
  }

  Future<void> register() async {
    if (busy) return;
    busy = true;
    statusMessage = null;
    notifyListeners();
    try {
      await saveCredentials();
      await _sip.register(
        username: extension.trim(),
        password: password,
        host: config.sipHost,
        port: config.sipPort,
        realm: config.sipRealm,
      );
    } catch (e) {
      statusMessage = e.toString();
      registration = SipRegistrationState.failed;
    } finally {
      busy = false;
      notifyListeners();
    }
  }

  Future<void> unregister() async {
    await _sip.unregister();
  }

  void tapDigit(String d) {
    if (digits.length >= 20) return;
    digits += d;
    notifyListeners();
  }

  void backspace() {
    if (digits.isEmpty) return;
    digits = digits.substring(0, digits.length - 1);
    notifyListeners();
  }

  void clearDigits() {
    digits = '';
    notifyListeners();
  }

  Future<void> placeCall() async {
    final dest = digits.trim();
    if (dest.isEmpty || !isRegistered) return;
    await _sip.call(dest);
  }

  Future<void> hangup() async {
    await _sip.hangup();
  }

  void _onSipEvent(SipEvent event) {
    switch (event) {
      case SipRegistrationChanged(:final state, :final message):
        registration = state;
        if (message != null) statusMessage = message;
      case SipCallStateChanged(:final state, :final peer):
        callState = state;
        this.peer = peer ?? this.peer;
        if (state == SipCallUiState.idle) {
          this.peer = null;
        }
      case SipIncomingCall(:final from):
        peer = from;
        callState = SipCallUiState.ringing;
    }
    notifyListeners();
  }

  @override
  void dispose() {
    _sub.cancel();
    super.dispose();
  }
}
