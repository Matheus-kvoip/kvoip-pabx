import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import request from 'supertest';
import { App } from 'supertest/types';
import { AppModule } from './../src/app.module';

describe('Kvoip API (e2e)', () => {
  let app: INestApplication<App>;
  let token: string;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    app.setGlobalPrefix('api');
    await app.init();

    const login = await request(app.getHttpServer())
      .post('/api/auth/login')
      .send({
        email: process.env.AUTH_EMAIL ?? 'admin@kvoip.com.br',
        password: process.env.AUTH_PASSWORD ?? 'kvoip123',
      })
      .expect(200);

    token = login.body.accessToken;
  });

  afterEach(async () => {
    await app.close();
  });

  it('/api/health (GET) is public', () => {
    return request(app.getHttpServer())
      .get('/api/health')
      .expect(200)
      .expect((res) => {
        expect(res.body.service).toBe('kvoip-server');
        expect(res.body.status).toBe('ok');
      });
  });

  it('/api/extensions (GET) requires auth', () => {
    return request(app.getHttpServer()).get('/api/extensions').expect(401);
  });

  it('/api/extensions (GET) with token', () => {
    return request(app.getHttpServer())
      .get('/api/extensions')
      .set('Authorization', `Bearer ${token}`)
      .expect(200)
      .expect((res) => {
        expect(Array.isArray(res.body)).toBe(true);
        expect(res.body.length).toBeGreaterThan(0);
      });
  });
});
