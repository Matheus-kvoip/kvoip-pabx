import { NestFactory } from '@nestjs/core';
import { ConfigService } from '@nestjs/config';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const config = app.get(ConfigService);

  const prefix = config.get<string>('API_PREFIX', 'api');
  app.setGlobalPrefix(prefix);

  const corsOrigin = config.get<string>('CORS_ORIGIN', 'http://localhost:3000');
  app.enableCors({
    origin: corsOrigin.split(',').map((item) => item.trim()),
    credentials: true,
  });

  const port = Number(config.get<string>('PORT', '3001'));
  await app.listen(port);

  const name = config.get<string>('API_NAME', 'kvoip-server');
  console.log(`${name} listening on http://localhost:${port}/${prefix}`);
}

bootstrap();
