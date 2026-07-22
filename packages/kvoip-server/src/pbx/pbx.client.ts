import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { CallRecord, CallState, Extension } from '@kvoip/shared';

type PbxRegistration = {
  aor: string;
  number: string;
  contact: string;
  expires: number;
  updatedAt: string;
};

type PbxCall = {
  id: string;
  direction: string;
  state: string;
  from: string;
  to: string;
  startedAt: string;
  answeredAt?: string;
  endedAt?: string;
  durationSec: number;
};

@Injectable()
export class PbxClient {
  private readonly logger = new Logger(PbxClient.name);
  private readonly baseUrl: string;
  private readonly enabled: boolean;

  constructor(private readonly config: ConfigService) {
    this.baseUrl = (
      this.config.get<string>('PBX_URL') ?? 'http://127.0.0.1:8080'
    ).replace(/\/$/, '');
    this.enabled =
      (this.config.get<string>('PBX_ENABLED') ?? 'true').toLowerCase() !==
      'false';
  }

  isEnabled() {
    return this.enabled;
  }

  async health(): Promise<boolean> {
    if (!this.enabled) return false;
    try {
      const res = await fetch(`${this.baseUrl}/health`, {
        signal: AbortSignal.timeout(1500),
      });
      return res.ok;
    } catch {
      return false;
    }
  }

  async getRegistrations(): Promise<PbxRegistration[]> {
    return this.getJSON<PbxRegistration[]>('/v1/registrations', []);
  }

  async getCalls(activeOnly = false): Promise<CallRecord[]> {
    const path = activeOnly ? '/v1/calls?active=true' : '/v1/calls';
    const rows = await this.getJSON<PbxCall[]>(path, []);
    return rows.map((row) => this.mapCall(row));
  }

  mapRegistrationsToExtensions(rows: PbxRegistration[]): Extension[] {
    return rows.map((row) => ({
      id: `sip-${row.number}`,
      number: row.number,
      displayName: row.number,
      status: 'online' as const,
      device: row.contact,
      createdAt: row.updatedAt,
    }));
  }

  private mapCall(row: PbxCall): CallRecord {
    const state = this.mapState(row.state);
    return {
      id: row.id,
      direction: 'internal',
      state,
      from: row.from,
      to: row.to,
      startedAt: row.startedAt,
      answeredAt: row.answeredAt,
      endedAt: row.endedAt,
      durationSec: row.durationSec,
    };
  }

  private mapState(state: string): CallState {
    switch (state) {
      case 'ringing':
        return 'ringing';
      case 'answered':
        return 'answered';
      case 'held':
        return 'held';
      case 'ended':
        return 'ended';
      default:
        return 'ringing';
    }
  }

  private async getJSON<T>(path: string, fallback: T): Promise<T> {
    if (!this.enabled) return fallback;
    try {
      const res = await fetch(`${this.baseUrl}${path}`, {
        signal: AbortSignal.timeout(2000),
      });
      if (!res.ok) {
        this.logger.warn(`PBX ${path} -> ${res.status}`);
        return fallback;
      }
      return (await res.json()) as T;
    } catch (err) {
      this.logger.warn(
        `PBX indisponível em ${this.baseUrl}${path}: ${
          err instanceof Error ? err.message : 'erro'
        }`,
      );
      return fallback;
    }
  }
}
