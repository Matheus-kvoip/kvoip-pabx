import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import type { Trunk } from '@kvoip/shared';
import { Repository } from 'typeorm';
import { TrunkEntity } from '../database/entities/trunk.entity';

@Injectable()
export class TrunksService {
  constructor(
    @InjectRepository(TrunkEntity)
    private readonly repo: Repository<TrunkEntity>,
  ) {}

  async findAll(): Promise<Trunk[]> {
    const rows = await this.repo.find({
      where: { enabled: true },
      order: { name: 'ASC' },
    });
    return rows.map((row) => this.toDto(row));
  }

  async findOne(id: string): Promise<Trunk> {
    const row = await this.repo.findOne({ where: { id } });
    if (!row) {
      throw new NotFoundException(`Tronco ${id} não encontrado`);
    }
    return this.toDto(row);
  }

  async countByStatus() {
    const rows = await this.findAll();
    const up = rows.filter((item) => item.status === 'up').length;
    return { up, total: rows.length };
  }

  private toDto(row: TrunkEntity): Trunk {
    return {
      id: row.id,
      name: row.name,
      host: row.host,
      port: row.port,
      protocol: row.protocol as Trunk['protocol'],
      status: row.status as Trunk['status'],
      concurrentCalls: row.concurrentCalls,
      maxChannels: row.maxChannels,
    };
  }
}
