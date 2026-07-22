import { Injectable, NotFoundException } from '@nestjs/common';
import type { Trunk } from '@kvoip/shared';

@Injectable()
export class TrunksService {
  private readonly trunks: Trunk[] = [
    {
      id: 'trk-primary',
      name: 'Tronco Principal',
      host: 'sip.operadora.com.br',
      port: 5060,
      protocol: 'udp',
      status: 'up',
      concurrentCalls: 4,
      maxChannels: 30,
    },
    {
      id: 'trk-backup',
      name: 'Tronco Backup',
      host: 'backup.sip.kvoip.local',
      port: 5061,
      protocol: 'tls',
      status: 'up',
      concurrentCalls: 0,
      maxChannels: 10,
    },
    {
      id: 'trk-mobile',
      name: 'Tronco Móvel',
      host: 'gw-mobile.kvoip.local',
      port: 5060,
      protocol: 'tcp',
      status: 'degraded',
      concurrentCalls: 2,
      maxChannels: 8,
    },
  ];

  findAll(): Trunk[] {
    return this.trunks;
  }

  findOne(id: string): Trunk {
    const trunk = this.trunks.find((item) => item.id === id);
    if (!trunk) {
      throw new NotFoundException(`Tronco ${id} não encontrado`);
    }
    return trunk;
  }

  countByStatus() {
    const up = this.trunks.filter((item) => item.status === 'up').length;
    return { up, total: this.trunks.length };
  }
}
