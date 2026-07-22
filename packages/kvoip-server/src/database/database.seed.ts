import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { InjectRepository } from '@nestjs/typeorm';
import * as bcrypt from 'bcryptjs';
import { Repository } from 'typeorm';
import { PbxClient } from '../pbx/pbx.client';
import { CallRecordEntity } from './entities/call-record.entity';
import { ExtensionEntity } from './entities/extension.entity';
import { TrunkEntity } from './entities/trunk.entity';
import { UserEntity } from './entities/user.entity';

@Injectable()
export class DatabaseSeedService implements OnModuleInit {
  private readonly logger = new Logger(DatabaseSeedService.name);

  constructor(
    private readonly config: ConfigService,
    private readonly pbx: PbxClient,
    @InjectRepository(UserEntity)
    private readonly users: Repository<UserEntity>,
    @InjectRepository(ExtensionEntity)
    private readonly extensions: Repository<ExtensionEntity>,
    @InjectRepository(TrunkEntity)
    private readonly trunks: Repository<TrunkEntity>,
    @InjectRepository(CallRecordEntity)
    private readonly calls: Repository<CallRecordEntity>,
  ) {}

  async onModuleInit() {
    await this.seedUsers();
    await this.seedExtensions();
    await this.seedTrunks();
    await this.seedCalls();
    await this.syncSipUsers();
  }

  private async seedUsers() {
    const count = await this.users.count();
    if (count > 0) return;

    const email = this.config.get<string>('AUTH_EMAIL', 'admin@kvoip.com.br');
    const password = this.config.get<string>('AUTH_PASSWORD', 'kvoip123');
    const name = this.config.get<string>('AUTH_NAME', 'Administrador');
    const passwordHash = await bcrypt.hash(password, 10);

    await this.users.save(
      this.users.create({
        email: email.toLowerCase(),
        passwordHash,
        name,
        role: 'admin',
      }),
    );
    this.logger.log(`Usuário admin seed: ${email}`);
  }

  private async seedExtensions() {
    const count = await this.extensions.count();
    if (count > 0) return;

    const rows: Partial<ExtensionEntity>[] = [
      {
        number: '1001',
        displayName: 'Ana Souza',
        email: 'ana@kvoip.com.br',
        device: 'Yealink T46U',
        sipPassword: 'kvoip123',
      },
      {
        number: '1002',
        displayName: 'Bruno Lima',
        email: 'bruno@kvoip.com.br',
        device: 'Softphone',
        sipPassword: 'kvoip123',
      },
      {
        number: '1003',
        displayName: 'Carla Mendes',
        email: 'carla@kvoip.com.br',
        device: 'Grandstream GXP2170',
        sipPassword: 'kvoip123',
      },
      {
        number: '1004',
        displayName: 'Diego Rocha',
        email: 'diego@kvoip.com.br',
        device: 'Softphone',
        sipPassword: 'kvoip123',
      },
      {
        number: '2000',
        displayName: 'Recepção',
        device: 'Fila / IVR',
        sipPassword: 'kvoip123',
      },
    ];

    await this.extensions.save(rows.map((row) => this.extensions.create(row)));
    this.logger.log(`Ramais seed: ${rows.length}`);
  }

  private async seedTrunks() {
    const count = await this.trunks.count();
    if (count > 0) return;

    await this.trunks.save([
      this.trunks.create({
        name: 'Tronco Principal',
        host: 'sip.operadora.com.br',
        port: 5060,
        protocol: 'udp',
        status: 'up',
        concurrentCalls: 0,
        maxChannels: 30,
      }),
      this.trunks.create({
        name: 'Tronco Backup',
        host: 'backup.sip.kvoip.local',
        port: 5061,
        protocol: 'tls',
        status: 'up',
        concurrentCalls: 0,
        maxChannels: 10,
      }),
      this.trunks.create({
        name: 'Tronco Móvel',
        host: 'gw-mobile.kvoip.local',
        port: 5060,
        protocol: 'tcp',
        status: 'degraded',
        concurrentCalls: 0,
        maxChannels: 8,
      }),
    ]);
    this.logger.log('Troncos seed OK');
  }

  private async seedCalls() {
    const count = await this.calls.count();
    if (count > 0) return;

    await this.calls.save(
      this.calls.create({
        id: 'call-history-1',
        direction: 'inbound',
        state: 'ended',
        from: '11987654321',
        to: '1001',
        startedAt: new Date('2026-07-21T18:10:00.000Z'),
        answeredAt: new Date('2026-07-21T18:10:12.000Z'),
        endedAt: new Date('2026-07-21T18:18:40.000Z'),
        durationSec: 508,
      }),
    );
  }

  async syncSipUsers() {
    const rows = await this.extensions.find({ where: { enabled: true } });
    const users: Record<string, string> = {};
    for (const row of rows) {
      users[row.number] = row.sipPassword;
    }
    await this.pbx.syncSipUsers(users);
  }
}
