import {
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Patch,
  Post,
} from '@nestjs/common';
import type { CreateExtensionInput } from '@kvoip/shared';
import { ExtensionsService } from './extensions.service';

@Controller('extensions')
export class ExtensionsController {
  constructor(private readonly extensionsService: ExtensionsService) {}

  @Get()
  findAll() {
    return this.extensionsService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.extensionsService.findOne(id);
  }

  @Post()
  create(@Body() body: CreateExtensionInput) {
    return this.extensionsService.create(body);
  }

  @Patch(':id')
  update(
    @Param('id') id: string,
    @Body() body: Partial<CreateExtensionInput>,
  ) {
    return this.extensionsService.update(id, body);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.extensionsService.remove(id);
  }
}
