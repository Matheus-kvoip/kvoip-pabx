import { MigrationInterface, QueryRunner } from 'typeorm';

export class InitialSchema1742600000000 implements MigrationInterface {
  name = 'InitialSchema1742600000000';

  public async up(queryRunner: QueryRunner): Promise<void> {
    const db = queryRunner.connection.options.type;

    if (db === 'postgres') {
      await queryRunner.query(`
        CREATE TABLE IF NOT EXISTS "users" (
          "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
          "email" varchar NOT NULL UNIQUE,
          "password_hash" varchar NOT NULL,
          "name" varchar NOT NULL,
          "role" varchar NOT NULL DEFAULT 'admin',
          "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
          "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
        );
      `);
      await queryRunner.query(`
        CREATE TABLE IF NOT EXISTS "extensions" (
          "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
          "number" varchar NOT NULL UNIQUE,
          "display_name" varchar NOT NULL,
          "email" varchar,
          "device" varchar,
          "sip_password" varchar NOT NULL,
          "enabled" boolean NOT NULL DEFAULT true,
          "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
          "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
        );
      `);
      await queryRunner.query(`
        CREATE TABLE IF NOT EXISTS "trunks" (
          "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
          "name" varchar NOT NULL,
          "host" varchar NOT NULL,
          "port" int NOT NULL DEFAULT 5060,
          "protocol" varchar NOT NULL DEFAULT 'udp',
          "status" varchar NOT NULL DEFAULT 'down',
          "concurrent_calls" int NOT NULL DEFAULT 0,
          "max_channels" int NOT NULL DEFAULT 30,
          "enabled" boolean NOT NULL DEFAULT true,
          "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
          "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
        );
      `);
      await queryRunner.query(`
        CREATE TABLE IF NOT EXISTS "call_records" (
          "id" varchar PRIMARY KEY,
          "direction" varchar NOT NULL,
          "state" varchar NOT NULL,
          "from_number" varchar NOT NULL,
          "to_number" varchar NOT NULL,
          "started_at" TIMESTAMPTZ NOT NULL,
          "answered_at" TIMESTAMPTZ,
          "ended_at" TIMESTAMPTZ,
          "duration_sec" int NOT NULL DEFAULT 0,
          "created_at" TIMESTAMPTZ NOT NULL DEFAULT now()
        );
      `);
      await queryRunner.query(`
        CREATE INDEX IF NOT EXISTS "idx_call_records_started_at"
        ON "call_records" ("started_at" DESC);
      `);
      return;
    }

    // sql.js / sqlite fallback
    await queryRunner.query(`
      CREATE TABLE IF NOT EXISTS "users" (
        "id" varchar PRIMARY KEY,
        "email" varchar NOT NULL UNIQUE,
        "password_hash" varchar NOT NULL,
        "name" varchar NOT NULL,
        "role" varchar NOT NULL DEFAULT 'admin',
        "created_at" datetime NOT NULL DEFAULT (datetime('now')),
        "updated_at" datetime NOT NULL DEFAULT (datetime('now'))
      );
    `);
    await queryRunner.query(`
      CREATE TABLE IF NOT EXISTS "extensions" (
        "id" varchar PRIMARY KEY,
        "number" varchar NOT NULL UNIQUE,
        "display_name" varchar NOT NULL,
        "email" varchar,
        "device" varchar,
        "sip_password" varchar NOT NULL,
        "enabled" boolean NOT NULL DEFAULT 1,
        "created_at" datetime NOT NULL DEFAULT (datetime('now')),
        "updated_at" datetime NOT NULL DEFAULT (datetime('now'))
      );
    `);
    await queryRunner.query(`
      CREATE TABLE IF NOT EXISTS "trunks" (
        "id" varchar PRIMARY KEY,
        "name" varchar NOT NULL,
        "host" varchar NOT NULL,
        "port" integer NOT NULL DEFAULT 5060,
        "protocol" varchar NOT NULL DEFAULT 'udp',
        "status" varchar NOT NULL DEFAULT 'down',
        "concurrent_calls" integer NOT NULL DEFAULT 0,
        "max_channels" integer NOT NULL DEFAULT 30,
        "enabled" boolean NOT NULL DEFAULT 1,
        "created_at" datetime NOT NULL DEFAULT (datetime('now')),
        "updated_at" datetime NOT NULL DEFAULT (datetime('now'))
      );
    `);
    await queryRunner.query(`
      CREATE TABLE IF NOT EXISTS "call_records" (
        "id" varchar PRIMARY KEY,
        "direction" varchar NOT NULL,
        "state" varchar NOT NULL,
        "from_number" varchar NOT NULL,
        "to_number" varchar NOT NULL,
        "started_at" datetime NOT NULL,
        "answered_at" datetime,
        "ended_at" datetime,
        "duration_sec" integer NOT NULL DEFAULT 0,
        "created_at" datetime NOT NULL DEFAULT (datetime('now'))
      );
    `);
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    await queryRunner.query(`DROP TABLE IF EXISTS "call_records"`);
    await queryRunner.query(`DROP TABLE IF EXISTS "trunks"`);
    await queryRunner.query(`DROP TABLE IF EXISTS "extensions"`);
    await queryRunner.query(`DROP TABLE IF EXISTS "users"`);
  }
}
