/** Saúde agregada do serviço. */
export enum HealthState {
  Ok = 'ok',
  Degraded = 'degraded',
}

export type HealthStatus = {
  status: HealthState | `${HealthState}`;
  service: string;
  version: string;
  uptimeSec: number;
  timestamp: string;
};
