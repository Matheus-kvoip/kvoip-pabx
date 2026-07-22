import {
  Body,
  Controller,
  Get,
  HttpCode,
  Post,
  UnauthorizedException,
} from '@nestjs/common';
import type { AuthUser, LoginInput } from '@kvoip/shared';
import { AuthService } from './auth.service';
import { Public } from './public.decorator';
import { CurrentUser } from './current-user.decorator';

@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Public()
  @Post('login')
  @HttpCode(200)
  login(@Body() body: LoginInput) {
    if (!body?.email || !body?.password) {
      throw new UnauthorizedException('E-mail e senha são obrigatórios');
    }
    return this.authService.login(body);
  }

  @Get('me')
  me(@CurrentUser() user: AuthUser) {
    return user;
  }

  @Post('logout')
  @HttpCode(200)
  logout() {
    return { ok: true };
  }
}
