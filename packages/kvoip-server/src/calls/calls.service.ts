import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import type { CallRecord } from '@kvoip/shared';
import { Repository } from 'typeorm';
import { CallRecordEntity } from '../database/entities/call-record.entity';
import { PbxClient } from '../pbx/pbx.client';

@Injectable()
export class CallsService {
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
      const exists = await this.repo.findOne({ where: { id: call.id } });
      if (exists) continue;
      await this.repo.save(
        this.repo.create({
          id: call.id,
          direction: call.direction,
          state: call.state,
          from: call.from,
          to: call.to,
          startedAt: new Date(call.startedAt),
          answeredAt: call.answeredAt ? new Date(call.answeredAt) : null,
          endedAt: call.endedAt ? new Date(call.endedAt) : null,
          durationSec: call.durationSec,
        }),
      );
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
