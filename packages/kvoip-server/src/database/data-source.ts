import { DataSource, DataSourceOptions } from 'typeorm';
import { CallRecordEntity } from './entities/call-record.entity';
import { ExtensionEntity } from './entities/extension.entity';
import { TrunkEntity } from './entities/trunk.entity';
import { UserEntity } from './entities/user.entity';
import { InitialSchema1742600000000 } from './migrations/1742600000000-InitialSchema';

const entities = [UserEntity, ExtensionEntity, TrunkEntity, CallRecordEntity];

export function buildDataSourceOptions(): DataSourceOptions {
  const dbType = (process.env.DB_TYPE || 'sqlite').toLowerCase();
  const sync = (process.env.DB_SYNC || 'false').toLowerCase() === 'true';
  const logging = (process.env.DB_LOGGING || 'false').toLowerCase() === 'true';

  if (dbType === 'postgres' || dbType === 'postgresql') {
    const url = process.env.DATABASE_URL;
    const base = {
      type: 'postgres' as const,
      entities,
      migrations: [InitialSchema1742600000000],
      synchronize: sync,
      logging,
    };
    if (url) {
      return { ...base, url };
    }
    return {
      ...base,
      host: process.env.DB_HOST || '127.0.0.1',
      port: Number(process.env.DB_PORT || '5432'),
      username: process.env.DB_USER || 'kvoip',
      password: process.env.DB_PASSWORD || 'kvoip',
      database: process.env.DB_NAME || 'kvoip',
    };
  }

  return {
    type: 'sqljs',
    location: process.env.SQLITE_PATH || 'data/kvoip.sqlite',
    autoSave: true,
    entities,
    migrations: [InitialSchema1742600000000],
    synchronize: sync || dbType === 'sqlite',
    logging,
  };
}

export default new DataSource(buildDataSourceOptions());
