import { Injectable, UnauthorizedException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { JwtService } from '@nestjs/jwt';
import type { AuthUser, LoginInput, LoginResponse } from '@kvoip/shared';
import * as bcrypt from 'bcryptjs';
import { Repository } from 'typeorm';
import { ConfigService } from '@nestjs/config';
import { UserEntity } from '../database/entities/user.entity';

@Injectable()
export class AuthService {
  constructor(
    private readonly config: ConfigService,
    private readonly jwtService: JwtService,
    @InjectRepository(UserEntity)
    private readonly users: Repository<UserEntity>,
  ) {}

  async login(input: LoginInput): Promise<LoginResponse> {
    const email = input.email.trim().toLowerCase();
    const user = await this.users.findOne({ where: { email } });
    if (!user) {
      throw new UnauthorizedException('Credenciais inválidas');
    }

    const ok = await bcrypt.compare(input.password, user.passwordHash);
    if (!ok) {
      throw new UnauthorizedException('Credenciais inválidas');
    }

    const authUser: AuthUser = {
      id: user.id,
      email: user.email,
      name: user.name,
    };

    const expiresIn = this.config.get<string>('JWT_EXPIRES_IN', '8h');
    const accessToken = await this.jwtService.signAsync({
      sub: authUser.id,
      email: authUser.email,
      name: authUser.name,
    });

    return {
      accessToken,
      tokenType: 'Bearer',
      expiresIn,
      user: authUser,
    };
  }
}
