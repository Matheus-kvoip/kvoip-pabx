import { Module } from '@nestjs/common';
import { ExtensionsModule } from '../extensions/extensions.module';
import { TrunksModule } from '../trunks/trunks.module';
import { CallsModule } from '../calls/calls.module';
import { PbxModule } from '../pbx/pbx.module';
import { DashboardController } from './dashboard.controller';
import { DashboardService } from './dashboard.service';

@Module({
  imports: [ExtensionsModule, TrunksModule, CallsModule, PbxModule],
  controllers: [DashboardController],
  providers: [DashboardService],
})
export class DashboardModule {}
