import type { CallDirection, CallState } from '../calls/call.enums';

/** Registro SIP exposto pela API HTTP do PBX. */
export type PbxRegistration = {
  aor: string;
  number: string;
  contact: string;
  expires: number;
  updatedAt: string;
};

/** Chamada exposta pela API HTTP do PBX. */
export type PbxCall = {
  id: string;
  direction: CallDirection | `${CallDirection}` | string;
  state: CallState | `${CallState}` | string;
  from: string;
  to: string;
  startedAt: string;
  answeredAt?: string;
  endedAt?: string;
  durationSec: number;
};

/** Payload de health do PBX (`GET /health`). */
export type PbxHealth = {
  status: string;
  service: string;
  version: string;
  uptimeSec: number;
  timestamp: string;
  sip: string;
  bindings: number;
  activeCalls: number;
};
