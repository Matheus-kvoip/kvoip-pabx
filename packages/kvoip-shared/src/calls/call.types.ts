import type { CallDirection, CallState } from './call.enums';

export type CallRecord = {
  id: string;
  direction: CallDirection | `${CallDirection}`;
  state: CallState | `${CallState}`;
  from: string;
  to: string;
  startedAt: string;
  answeredAt?: string;
  endedAt?: string;
  durationSec: number;
};
