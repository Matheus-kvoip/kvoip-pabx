import { Module } from '@nestjs/common';
import { APP_GUARD } from '@nestjs/core';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { AuthModule } from './auth/auth.module';
import { AuthGuard } from './auth/auth.guard';
import { HealthModule } from './health/health.module';
import { ExtensionsModule } from './extensions/extensions.module';
import { TrunksModule } from './trunks/trunks.module';
import { CallsModule } from './calls/calls.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { DatabaseModule } from './database/database.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: ['.env', '.env.local'],
    }),
    DatabaseModule,
    AuthModule,
    HealthModule,
    ExtensionsModule,
    TrunksModule,
    CallsModule,
    DashboardModule,
  ],
  providers: [
    {
      provide: APP_GUARD,
      useClass: AuthGuard,
    },
    ConfigService,
  ],
})
export class AppModule {}
