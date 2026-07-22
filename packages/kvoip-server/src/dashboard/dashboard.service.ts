import { Injectable } from '@nestjs/common';
import type { DashboardStats } from '@kvoip/shared';
import { ExtensionsService } from '../extensions/extensions.service';
import { TrunksService } from '../trunks/trunks.service';
import { CallsService } from '../calls/calls.service';
import { PbxClient } from '../pbx/pbx.client';

@Injectable()
export class DashboardService {
  constructor(
    private readonly extensionsService: ExtensionsService,
    private readonly trunksService: TrunksService,
    private readonly callsService: CallsService,
    private readonly pbx: PbxClient,
  ) {}

  async getStats(): Promise<DashboardStats & { pbxOnline: boolean }> {
    const [extensions, trunks, calls, pbxOnline] = await Promise.all([
      this.extensionsService.countByStatus(),
      this.trunksService.countByStatus(),
      this.callsService.statsToday(),
      this.pbx.health(),
    ]);

    return {
      activeCalls: calls.active,
      extensionsOnline: extensions.online,
      extensionsTotal: extensions.total,
      trunksUp: trunks.up,
      trunksTotal: trunks.total,
      callsToday: calls.today,
      avgAnswerSec: calls.avgAnswerSec,
      pbxOnline,
    };
  }
}
