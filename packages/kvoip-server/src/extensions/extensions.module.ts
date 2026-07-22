import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { DatabaseModule } from '../database/database.module';
import { ExtensionEntity } from '../database/entities/extension.entity';
import { PbxModule } from '../pbx/pbx.module';
import { ExtensionsController } from './extensions.controller';
import { ExtensionsService } from './extensions.service';

@Module({
  imports: [
    PbxModule,
    DatabaseModule,
    TypeOrmModule.forFeature([ExtensionEntity]),
  ],
  controllers: [ExtensionsController],
  providers: [ExtensionsService],
  exports: [ExtensionsService],
})
export class ExtensionsModule {}
