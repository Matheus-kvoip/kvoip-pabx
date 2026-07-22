import {
  BadRequestException,
  Body,
  Controller,
  ForbiddenException,
  Get,
  Headers,
  Post,
  Query,
} from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Public } from '../auth/public.decorator';
import { CallsService, type CdrIngestInput } from './calls.service';

@Controller('calls')
export class CallsController {
  constructor(
    private readonly callsService: CallsService,
    private readonly config: ConfigService,
  ) {}

  @Get()
  findAll(@Query('active') active?: string) {
    if (active === 'true') {
      return this.callsService.findActive();
    }
    return this.callsService.findAll();
  }

  /** Webhook do PBX no hangup — autenticado por segredo compartilhado. */
  @Public()
  @Post('cdr')
  ingestCdr(
    @Body() body: CdrIngestInput,
    @Headers('x-cdr-secret') secret?: string,
  ) {
    const expected = this.config.get<string>('CDR_WEBHOOK_SECRET', '');
    if (expected && secret !== expected) {
      throw new ForbiddenException('CDR secret inválido');
    }
    if (!body?.id || !body?.from || !body?.to || !body?.startedAt) {
      throw new BadRequestException('CDR inválido');
    }
    return this.callsService.ingestCdr(body);
  }
}
