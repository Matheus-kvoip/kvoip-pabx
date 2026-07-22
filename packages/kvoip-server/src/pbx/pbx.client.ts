import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import {
  CallDirection,
  CallState,
  type CallRecord,
  type Extension,
  type PbxCall,
  type PbxRegistration,
} from '@kvoip/shared';


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
    return {
      id: row.id,
      direction: this.mapDirection(row.direction),
      state: this.mapState(row.state),
      from: row.from,
      to: row.to,
      startedAt: row.startedAt,
      answeredAt: row.answeredAt,
      endedAt: row.endedAt,
      durationSec: row.durationSec,
    };
  }

  private mapDirection(direction: string): CallDirection {
    switch (direction) {
      case CallDirection.Inbound:
        return CallDirection.Inbound;
      case CallDirection.Outbound:
        return CallDirection.Outbound;
      case CallDirection.Internal:
        return CallDirection.Internal;
      default:
        return CallDirection.Internal;
    }
  }

  private mapState(state: string): CallState {
    switch (state) {
      case CallState.Idle:
        return CallState.Idle;
      case CallState.Ringing:
        return CallState.Ringing;
      case CallState.Answered:
        return CallState.Answered;
      case CallState.Held:
        return CallState.Held;
      case CallState.Ended:
        return CallState.Ended;
      default:
        return CallState.Ringing;
    }
  }

  async syncSipUsers(users: Record<string, string>): Promise<boolean> {
    if (!this.enabled) return false;
    try {
      const res = await fetch(`${this.baseUrl}/v1/sip-users`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ users }),
        signal: AbortSignal.timeout(3000),
      });
      if (!res.ok) {
        this.logger.warn(`PBX sync sip-users -> ${res.status}`);
        return false;
      }
      this.logger.log(`SIP users sync OK (${Object.keys(users).length})`);
      return true;
    } catch (err) {
      this.logger.warn(
        `PBX sync falhou: ${err instanceof Error ? err.message : 'erro'}`,
      );
      return false;
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
