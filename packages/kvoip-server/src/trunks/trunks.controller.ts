import { Controller, Get, Param } from '@nestjs/common';
import { TrunksService } from './trunks.service';

@Controller('trunks')
export class TrunksController {
  constructor(private readonly trunksService: TrunksService) {}

  @Get()
  findAll() {
    return this.trunksService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.trunksService.findOne(id);
  }
}
