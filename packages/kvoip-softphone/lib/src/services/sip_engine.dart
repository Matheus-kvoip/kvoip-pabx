import 'dart:async';

/// Contrato SIP — implementação real virá com PJSIP/Linphone SDK.
abstract class SipEngine {
  Future<void> register({
    required String username,
    required String password,
    required String host,
    required int port,
    required String realm,
  });

  Future<void> unregister();

  Future<void> call(String destination);

  Future<void> hangup();

  Future<void> answer();

  Stream<SipEvent> get events;

  SipRegistrationState get registrationState;
}

enum SipRegistrationState { idle, registering, registered, failed }

sealed class SipEvent {
  const SipEvent();
}

class SipRegistrationChanged extends SipEvent {
  const SipRegistrationChanged(this.state, {this.message});
  final SipRegistrationState state;
  final String? message;
}

class SipIncomingCall extends SipEvent {
  const SipIncomingCall(this.from);
  final String from;
}

class SipCallStateChanged extends SipEvent {
  const SipCallStateChanged(this.state, {this.peer});
  final SipCallUiState state;
  final String? peer;
}

enum SipCallUiState { idle, dialing, ringing, active, ended }

/// Stub para desenvolver UI sem SDK nativo ainda.
class MockSipEngine implements SipEngine {
  MockSipEngine() {
    _controller = StreamController<SipEvent>.broadcast();
  }

  late final StreamController<SipEvent> _controller;
  SipRegistrationState _registration = SipRegistrationState.idle;

  @override
  SipRegistrationState get registrationState => _registration;

  @override
  Stream<SipEvent> get events => _controller.stream;

  @override
  Future<void> register({
    required String username,
    required String password,
    required String host,
    required int port,
    required String realm,
  }) async {
    _registration = SipRegistrationState.registering;
    _controller.add(SipRegistrationChanged(_registration));
    await Future<void>.delayed(const Duration(milliseconds: 600));
    _registration = SipRegistrationState.registered;
    _controller.add(
      SipRegistrationChanged(
        _registration,
        message: 'Mock: $username@$realm → $host:$port',
      ),
    );
  }

  @override
  Future<void> unregister() async {
    _registration = SipRegistrationState.idle;
    _controller.add(SipRegistrationChanged(_registration));
  }

  @override
  Future<void> call(String destination) async {
    _controller.add(
      SipCallStateChanged(SipCallUiState.dialing, peer: destination),
    );
    await Future<void>.delayed(const Duration(milliseconds: 400));
    _controller.add(
      SipCallStateChanged(SipCallUiState.ringing, peer: destination),
    );
    await Future<void>.delayed(const Duration(milliseconds: 800));
    _controller.add(
      SipCallStateChanged(SipCallUiState.active, peer: destination),
    );
  }

  @override
  Future<void> hangup() async {
    _controller.add(const SipCallStateChanged(SipCallUiState.ended));
    await Future<void>.delayed(const Duration(milliseconds: 200));
    _controller.add(const SipCallStateChanged(SipCallUiState.idle));
  }

  @override
  Future<void> answer() async {
    _controller.add(const SipCallStateChanged(SipCallUiState.active));
  }

  Future<void> dispose() async {
    await _controller.close();
  }
}
