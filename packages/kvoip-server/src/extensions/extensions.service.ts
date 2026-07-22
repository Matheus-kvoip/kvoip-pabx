import {
  Injectable,
  NotFoundException,
  ConflictException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import type { CreateExtensionInput, Extension } from '@kvoip/shared';
import { Repository } from 'typeorm';
import { ExtensionEntity } from '../database/entities/extension.entity';
import { DatabaseSeedService } from '../database/database.seed';
import { PbxClient } from '../pbx/pbx.client';

@Injectable()
export class ExtensionsService {
  constructor(
    @InjectRepository(ExtensionEntity)
    private readonly repo: Repository<ExtensionEntity>,
    private readonly pbx: PbxClient,
    private readonly seed: DatabaseSeedService,
  ) {}

  async findAll(): Promise<Extension[]> {
    const directory = await this.repo.find({ order: { number: 'ASC' } });
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
      const base = this.toDto(ext);
      const reg = byNumber.get(ext.number);
      if (!reg) {
        return { ...base, status: 'offline' as const };
      }
      byNumber.delete(ext.number);
      let status: Extension['status'] = 'online';
      if (busy.has(ext.number)) status = 'busy';
      else if (ringing.has(ext.number)) status = 'ringing';
      return {
        ...base,
        status,
        device: reg.contact || base.device,
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

  async create(input: CreateExtensionInput): Promise<Extension> {
    const exists = await this.repo.findOne({ where: { number: input.number } });
    if (exists) {
      throw new ConflictException(`Ramal ${input.number} já existe`);
    }
    const entity = await this.repo.save(
      this.repo.create({
        number: input.number,
        displayName: input.displayName,
        email: input.email ?? null,
        device: input.device ?? null,
        sipPassword: input.sipPassword?.trim() || 'kvoip123',
        enabled: true,
      }),
    );
    await this.seed.syncSipUsers();
    return this.toDto(entity);
  }

  async update(
    id: string,
    input: Partial<CreateExtensionInput>,
  ): Promise<Extension> {
    const entity = await this.repo.findOne({ where: { id } });
    if (!entity) {
      throw new NotFoundException(`Ramal ${id} não encontrado`);
    }
    if (input.number && input.number !== entity.number) {
      const clash = await this.repo.findOne({ where: { number: input.number } });
      if (clash) {
        throw new ConflictException(`Ramal ${input.number} já existe`);
      }
      entity.number = input.number;
    }
    entity.displayName = input.displayName ?? entity.displayName;
    entity.email = input.email ?? entity.email;
    entity.device = input.device ?? entity.device;
    if (input.sipPassword?.trim()) {
      entity.sipPassword = input.sipPassword.trim();
    }
    await this.repo.save(entity);
    await this.seed.syncSipUsers();
    return this.toDto(entity);
  }

  async remove(id: string): Promise<{ ok: true }> {
    const entity = await this.repo.findOne({ where: { id } });
    if (!entity) {
      throw new NotFoundException(`Ramal ${id} não encontrado`);
    }
    await this.repo.remove(entity);
    await this.seed.syncSipUsers();
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

  private toDto(entity: ExtensionEntity): Extension {
    return {
      id: entity.id,
      number: entity.number,
      displayName: entity.displayName,
      email: entity.email ?? undefined,
      device: entity.device ?? undefined,
      status: 'offline',
      createdAt: entity.createdAt.toISOString(),
    };
  }
}
