import { Injectable, UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { JwtService } from '@nestjs/jwt';
import type { AuthUser, LoginInput, LoginResponse } from '@kvoip/shared';

@Injectable()
export class AuthService {
  constructor(
    private readonly config: ConfigService,
    private readonly jwtService: JwtService,
  ) {}

  async login(input: LoginInput): Promise<LoginResponse> {
    const email = this.config.getOrThrow<string>('AUTH_EMAIL');
    const password = this.config.getOrThrow<string>('AUTH_PASSWORD');
    const name = this.config.get<string>('AUTH_NAME', 'Administrador');

    if (
      input.email.trim().toLowerCase() !== email.toLowerCase() ||
      input.password !== password
    ) {
      throw new UnauthorizedException('Credenciais inválidas');
    }

    const user: AuthUser = {
      id: 'admin',
      email,
      name,
    };

    const expiresIn = this.config.get<string>('JWT_EXPIRES_IN', '8h');
    const accessToken = await this.jwtService.signAsync({
      sub: user.id,
      email: user.email,
      name: user.name,
    });

    return {
      accessToken,
      tokenType: 'Bearer',
      expiresIn,
      user,
    };
  }
}
