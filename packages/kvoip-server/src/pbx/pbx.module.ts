import { Module } from '@nestjs/common';
import { PbxClient } from './pbx.client';

@Module({
  providers: [PbxClient],
  exports: [PbxClient],
})
export class PbxModule {}
