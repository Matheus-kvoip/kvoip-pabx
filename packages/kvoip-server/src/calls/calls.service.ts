import { Injectable } from '@nestjs/common';
import type { CallRecord } from '@kvoip/shared';
import { PbxClient } from '../pbx/pbx.client';

@Injectable()
export class CallsService {
  private readonly history: CallRecord[] = [
    {
      id: 'call-history-1',
      direction: 'inbound',
      state: 'ended',
      from: '11987654321',
      to: '1001',
      startedAt: '2026-07-21T18:10:00.000Z',
      answeredAt: '2026-07-21T18:10:12.000Z',
      endedAt: '2026-07-21T18:18:40.000Z',
      durationSec: 508,
    },
  ];

  constructor(private readonly pbx: PbxClient) {}

  async findAll(): Promise<CallRecord[]> {
    const live = await this.pbx.getCalls(false);
    if (live.length === 0) {
      return [...this.history];
    }
    const ids = new Set(live.map((c) => c.id));
    return [...live, ...this.history.filter((c) => !ids.has(c.id))];
  }

  async findActive(): Promise<CallRecord[]> {
    const live = await this.pbx.getCalls(true);
    if (live.length > 0) return live;
    return (await this.findAll()).filter((call) => call.state !== 'ended');
  }

  async statsToday() {
    const all = await this.findAll();
    const active = all.filter((call) => call.state !== 'ended').length;
    const today = all.length;
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
}
