import { Injectable, NotFoundException } from '@nestjs/common';
import type { CreateExtensionInput, Extension } from '@kvoip/shared';
import { randomUUID } from 'crypto';
import { PbxClient } from '../pbx/pbx.client';

@Injectable()
export class ExtensionsService {
  private extensions: Extension[] = [
    {
      id: 'ext-1001',
      number: '1001',
      displayName: 'Ana Souza',
      email: 'ana@kvoip.com.br',
      status: 'offline',
      device: 'Yealink T46U',
      createdAt: '2026-01-10T12:00:00.000Z',
    },
    {
      id: 'ext-1002',
      number: '1002',
      displayName: 'Bruno Lima',
      email: 'bruno@kvoip.com.br',
      status: 'offline',
      device: 'Softphone',
      createdAt: '2026-01-12T09:30:00.000Z',
    },
    {
      id: 'ext-1003',
      number: '1003',
      displayName: 'Carla Mendes',
      email: 'carla@kvoip.com.br',
      status: 'offline',
      device: 'Grandstream GXP2170',
      createdAt: '2026-02-01T15:10:00.000Z',
    },
    {
      id: 'ext-1004',
      number: '1004',
      displayName: 'Diego Rocha',
      email: 'diego@kvoip.com.br',
      status: 'offline',
      device: 'Softphone',
      createdAt: '2026-02-18T11:00:00.000Z',
    },
    {
      id: 'ext-2000',
      number: '2000',
      displayName: 'Recepção',
      status: 'offline',
      device: 'Fila / IVR',
      createdAt: '2026-01-05T08:00:00.000Z',
    },
  ];

  constructor(private readonly pbx: PbxClient) {}

  async findAll(): Promise<Extension[]> {
    const directory = [...this.extensions];
    const regs = await this.pbx.getRegistrations();
    const byNumber = new Map(regs.map((r) => [r.number, r]));
    const activeCalls = await this.pbx.getCalls(true);
    const busy = new Set<string>();
    const ringing = new Set<string>();
    for (const call of activeCalls) {
      const parties = [call.from, call.to];
      if (call.state === 'ringing') {
        parties.forEach((p) => ringing.add(p));
      } else if (call.state === 'answered' || call.state === 'held') {
        parties.forEach((p) => busy.add(p));
      }
    }

    const merged = directory.map((ext) => {
      const reg = byNumber.get(ext.number);
      if (!reg) {
        return { ...ext, status: 'offline' as const };
      }
      byNumber.delete(ext.number);
      let status: Extension['status'] = 'online';
      if (busy.has(ext.number)) status = 'busy';
      else if (ringing.has(ext.number)) status = 'ringing';
      return {
        ...ext,
        status,
        device: reg.contact || ext.device,
      };
    });

    for (const reg of byNumber.values()) {
      let status: Extension['status'] = 'online';
      if (busy.has(reg.number)) status = 'busy';
      else if (ringing.has(reg.number)) status = 'ringing';
      merged.push({
        id: `sip-${reg.number}`,
        number: reg.number,
        displayName: reg.number,
        status,
        device: reg.contact,
        createdAt: reg.updatedAt,
      });
    }

    return merged.sort((a, b) => a.number.localeCompare(b.number));
  }

  async findOne(id: string): Promise<Extension> {
    const all = await this.findAll();
    const extension = all.find((item) => item.id === id);
    if (!extension) {
      throw new NotFoundException(`Ramal ${id} não encontrado`);
    }
    return extension;
  }

  create(input: CreateExtensionInput): Extension {
    const extension: Extension = {
      id: `ext-${randomUUID().slice(0, 8)}`,
      number: input.number,
      displayName: input.displayName,
      email: input.email,
      device: input.device,
      status: 'offline',
      createdAt: new Date().toISOString(),
    };
    this.extensions.push(extension);
    return extension;
  }

  update(id: string, input: Partial<CreateExtensionInput>): Extension {
    const extension = this.extensions.find((item) => item.id === id);
    if (!extension) {
      throw new NotFoundException(`Ramal ${id} não encontrado`);
    }
    Object.assign(extension, {
      number: input.number ?? extension.number,
      displayName: input.displayName ?? extension.displayName,
      email: input.email ?? extension.email,
      device: input.device ?? extension.device,
    });
    return extension;
  }

  remove(id: string): { ok: true } {
    const exists = this.extensions.some((item) => item.id === id);
    if (!exists) {
      throw new NotFoundException(`Ramal ${id} não encontrado`);
    }
    this.extensions = this.extensions.filter((item) => item.id !== id);
    return { ok: true };
  }

  async countByStatus() {
    const list = await this.findAll();
    const online = list.filter(
      (item) =>
        item.status === 'online' ||
        item.status === 'busy' ||
        item.status === 'ringing',
    ).length;
    return { online, total: list.length };
  }
}
