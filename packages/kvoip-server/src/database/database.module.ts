import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule } from '@nestjs/typeorm';
import { mkdirSync } from 'fs';
import { dirname } from 'path';
import { PbxModule } from '../pbx/pbx.module';
import { DatabaseSeedService } from './database.seed';
import { CallRecordEntity } from './entities/call-record.entity';
import { ExtensionEntity } from './entities/extension.entity';
import { TrunkEntity } from './entities/trunk.entity';
import { UserEntity } from './entities/user.entity';

const entities = [UserEntity, ExtensionEntity, TrunkEntity, CallRecordEntity];

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (config: ConfigService) => {
        const dbType = (
          config.get<string>('DB_TYPE', 'sqlite') || 'sqlite'
        ).toLowerCase();
        const common = {
          entities,
          synchronize: config.get<string>('DB_SYNC', 'true') !== 'false',
          logging: config.get<string>('DB_LOGGING', 'false') === 'true',
        };

        if (dbType === 'postgres' || dbType === 'postgresql') {
          const url = config.get<string>('DATABASE_URL');
          if (url) {
            return { ...common, type: 'postgres' as const, url };
          }
          return {
            ...common,
            type: 'postgres' as const,
            host: config.get<string>('DB_HOST', '127.0.0.1'),
            port: Number(config.get<string>('DB_PORT', '5432')),
            username: config.get<string>('DB_USER', 'kvoip'),
            password: config.get<string>('DB_PASSWORD', 'kvoip'),
            database: config.get<string>('DB_NAME', 'kvoip'),
          };
        }

        const database = config.get<string>(
          'SQLITE_PATH',
          'data/kvoip.sqlite',
        );
        mkdirSync(dirname(database), { recursive: true });
        return {
          ...common,
          type: 'sqljs' as const,
          location: database,
          autoSave: true,
          autoSaveInterval: 1000,
        };
      },
    }),
    TypeOrmModule.forFeature(entities),
    PbxModule,
  ],
  providers: [DatabaseSeedService],
  exports: [TypeOrmModule, DatabaseSeedService],
})
export class DatabaseModule {}
