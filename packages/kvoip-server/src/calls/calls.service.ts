import { Injectable, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import type { CallRecord } from '@kvoip/shared';
import { Repository } from 'typeorm';
import { CallRecordEntity } from '../database/entities/call-record.entity';
import { PbxClient } from '../pbx/pbx.client';

export type CdrIngestInput = {
  id: string;
  direction?: string;
  state?: string;
  from: string;
  to: string;
  startedAt: string;
  answeredAt?: string;
  endedAt?: string;
  durationSec?: number;
};

@Injectable()
export class CallsService {
  private readonly logger = new Logger(CallsService.name);

  constructor(
    @InjectRepository(CallRecordEntity)
    private readonly repo: Repository<CallRecordEntity>,
    private readonly pbx: PbxClient,
  ) {}

  async findAll(): Promise<CallRecord[]> {
    const live = await this.pbx.getCalls(false);
    await this.persistEnded(live);

    const history = await this.repo.find({
      order: { startedAt: 'DESC' },
      take: 200,
    });
    const historyDto = history.map((row) => this.toDto(row));
    const ids = new Set(live.map((c) => c.id));
    return [...live, ...historyDto.filter((c) => !ids.has(c.id))];
  }

  async findActive(): Promise<CallRecord[]> {
    const live = await this.pbx.getCalls(true);
    if (live.length > 0) return live;
    return (await this.findAll()).filter((call) => call.state !== 'ended');
  }

  async ingestCdr(input: CdrIngestInput): Promise<CallRecord> {
    const startedAt = new Date(input.startedAt);
    const answeredAt = input.answeredAt ? new Date(input.answeredAt) : null;
    const endedAt = input.endedAt ? new Date(input.endedAt) : new Date();
    let durationSec = input.durationSec ?? 0;
    if (!input.durationSec && !Number.isNaN(startedAt.getTime())) {
      durationSec = Math.max(
        0,
        Math.round((endedAt.getTime() - startedAt.getTime()) / 1000),
      );
    }

    const existing = await this.repo.findOne({ where: { id: input.id } });
    const row =
      existing ??
      this.repo.create({
        id: input.id,
      });

    row.direction = input.direction || 'internal';
    row.state = input.state || 'ended';
    row.from = input.from;
    row.to = input.to;
    row.startedAt = startedAt;
    row.answeredAt = answeredAt;
    row.endedAt = endedAt;
    row.durationSec = durationSec;

    const saved = await this.repo.save(row);
    this.logger.log(`CDR salvo ${saved.id} (${saved.from} → ${saved.to})`);
    return this.toDto(saved);
  }

  async statsToday() {
    const all = await this.findAll();
    const active = all.filter((call) => call.state !== 'ended').length;
    const startOfDay = new Date();
    startOfDay.setHours(0, 0, 0, 0);
    const today = all.filter(
      (call) => new Date(call.startedAt).getTime() >= startOfDay.getTime(),
    ).length;
    const answered = all.filter((call) => call.answeredAt);
    const avgAnswerSec =
      answered.length === 0
        ? 0
        : Math.round(
            answered.reduce((sum, call) => {
              const start = new Date(call.startedAt).getTime();
              const answer = new Date(call.answeredAt!).getTime();
              return sum + (answer - start) / 1000;
            }, 0) / answered.length,
          );
    return { active, today, avgAnswerSec };
  }

  private async persistEnded(live: CallRecord[]) {
    const ended = live.filter((c) => c.state === 'ended');
    for (const call of ended) {
      await this.ingestCdr({
        id: call.id,
        direction: call.direction,
        state: call.state,
        from: call.from,
        to: call.to,
        startedAt: call.startedAt,
        answeredAt: call.answeredAt,
        endedAt: call.endedAt,
        durationSec: call.durationSec,
      });
    }
  }

  private toDto(row: CallRecordEntity): CallRecord {
    return {
      id: row.id,
      direction: row.direction as CallRecord['direction'],
      state: row.state as CallRecord['state'],
      from: row.from,
      to: row.to,
      startedAt: row.startedAt.toISOString(),
      answeredAt: row.answeredAt?.toISOString(),
      endedAt: row.endedAt?.toISOString(),
      durationSec: row.durationSec,
    };
  }
}
