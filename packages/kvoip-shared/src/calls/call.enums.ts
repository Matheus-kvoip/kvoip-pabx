/** Direção da chamada. */
export enum CallDirection {
  Inbound = 'inbound',
  Outbound = 'outbound',
  Internal = 'internal',
}

/** Estado da sessão de chamada (alinhado ao PBX). */
export enum CallState {
  Idle = 'idle',
  Ringing = 'ringing',
  Answered = 'answered',
  Held = 'held',
  Ended = 'ended',
}
