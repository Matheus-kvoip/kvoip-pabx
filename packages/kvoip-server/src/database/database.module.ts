import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { TypeOrmModule, type TypeOrmModuleOptions } from '@nestjs/typeorm';
import { mkdirSync } from 'fs';
import { dirname } from 'path';
import { PbxModule } from '../pbx/pbx.module';
import { DatabaseSeedService } from './database.seed';
import { CallRecordEntity } from './entities/call-record.entity';
import { ExtensionEntity } from './entities/extension.entity';
import { TrunkEntity } from './entities/trunk.entity';
import { UserEntity } from './entities/user.entity';
import { InitialSchema1742600000000 } from './migrations/1742600000000-InitialSchema';

const entities = [UserEntity, ExtensionEntity, TrunkEntity, CallRecordEntity];

@Module({
  imports: [
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (config: ConfigService): TypeOrmModuleOptions => {
        const dbType = (
          config.get<string>('DB_TYPE', 'sqlite') || 'sqlite'
        ).toLowerCase();
        const isPostgres = dbType === 'postgres' || dbType === 'postgresql';
        const syncExplicit = config.get<string>('DB_SYNC');
        const synchronize =
          syncExplicit !== undefined
            ? syncExplicit === 'true'
            : !isPostgres;
        const migrationsRun =
          (config.get<string>('DB_MIGRATIONS_RUN') ??
            (isPostgres ? 'true' : 'false')) === 'true';
        const logging = config.get<string>('DB_LOGGING', 'false') === 'true';

        if (isPostgres) {
          const url = config.get<string>('DATABASE_URL');
          const base: TypeOrmModuleOptions = {
            type: 'postgres',
            entities,
            migrations: [InitialSchema1742600000000],
            synchronize,
            migrationsRun,
            logging,
          };
          if (url) {
            return { ...base, url };
          }
          return {
            ...base,
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
          type: 'sqljs',
          location: database,
          autoSave: true,
          entities,
          migrations: [InitialSchema1742600000000],
          synchronize,
          migrationsRun,
          logging,
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
