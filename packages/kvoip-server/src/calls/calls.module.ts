import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { CallRecordEntity } from '../database/entities/call-record.entity';
import { PbxModule } from '../pbx/pbx.module';
import { CallsController } from './calls.controller';
import { CallsService } from './calls.service';

@Module({
  imports: [PbxModule, TypeOrmModule.forFeature([CallRecordEntity])],
  controllers: [CallsController],
  providers: [CallsService],
  exports: [CallsService],
})
export class CallsModule {}
