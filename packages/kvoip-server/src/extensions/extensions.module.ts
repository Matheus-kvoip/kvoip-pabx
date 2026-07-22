import { Module } from '@nestjs/common';
import { PbxModule } from '../pbx/pbx.module';
import { ExtensionsController } from './extensions.controller';
import { ExtensionsService } from './extensions.service';

@Module({
  imports: [PbxModule],
  controllers: [ExtensionsController],
  providers: [ExtensionsService],
  exports: [ExtensionsService],
})
export class ExtensionsModule {}
