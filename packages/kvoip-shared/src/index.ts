export type ExtensionStatus = 'online' | 'offline' | 'busy' | 'ringing';

export type Extension = {
  id: string;
  number: string;
  displayName: string;
  email?: string;
  status: ExtensionStatus;
  device?: string;
  createdAt: string;
};

export type TrunkStatus = 'up' | 'down' | 'degraded';

export type Trunk = {
  id: string;
  name: string;
  host: string;
  port: number;
  protocol: 'udp' | 'tcp' | 'tls';
  status: TrunkStatus;
  concurrentCalls: number;
  maxChannels: number;
};

export type CallDirection = 'inbound' | 'outbound' | 'internal';
export type CallState = 'ringing' | 'answered' | 'held' | 'ended';

export type CallRecord = {
  id: string;
  direction: CallDirection;
  state: CallState;
  from: string;
  to: string;
  startedAt: string;
  answeredAt?: string;
  endedAt?: string;
  durationSec: number;
};

export type DashboardStats = {
  activeCalls: number;
  extensionsOnline: number;
  extensionsTotal: number;
  trunksUp: number;
  trunksTotal: number;
  callsToday: number;
  avgAnswerSec: number;
  pbxOnline?: boolean;
};

export type HealthStatus = {
  status: 'ok' | 'degraded';
  service: string;
  version: string;
  uptimeSec: number;
  timestamp: string;
};

export type CreateExtensionInput = {
  number: string;
  displayName: string;
  email?: string;
  device?: string;
};

export type AuthUser = {
  id: string;
  email: string;
  name: string;
};

export type LoginInput = {
  email: string;
  password: string;
};

export type LoginResponse = {
  accessToken: string;
  tokenType: 'Bearer';
  expiresIn: string;
  user: AuthUser;
};
