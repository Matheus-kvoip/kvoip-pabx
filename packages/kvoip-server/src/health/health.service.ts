import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { HealthStatus } from '@kvoip/shared';

@Injectable()
export class HealthService {
  private readonly startedAt = Date.now();

  constructor(private readonly config: ConfigService) {}

  getStatus(): HealthStatus {
    return {
      status: 'ok',
      service: this.config.get<string>('API_NAME', 'kvoip-server'),
      version: this.config.get<string>('API_VERSION', '0.1.0'),
      uptimeSec: Math.floor((Date.now() - this.startedAt) / 1000),
      timestamp: new Date().toISOString(),
    };
  }
}
