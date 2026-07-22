import { Controller, Get, Query } from '@nestjs/common';
import { CallsService } from './calls.service';

@Controller('calls')
export class CallsController {
  constructor(private readonly callsService: CallsService) {}

  @Get()
  findAll(@Query('active') active?: string) {
    if (active === 'true') {
      return this.callsService.findActive();
    }
    return this.callsService.findAll();
  }
}
