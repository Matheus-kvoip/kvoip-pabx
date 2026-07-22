import { Column, CreateDateColumn, Entity, PrimaryColumn } from 'typeorm';

@Entity('call_records')
export class CallRecordEntity {
  @PrimaryColumn()
  id!: string;

  @Column()
  direction!: string;

  @Column()
  state!: string;

  @Column({ name: 'from_number' })
  from!: string;

  @Column({ name: 'to_number' })
  to!: string;

  @Column({ name: 'started_at', type: 'datetime' })
  startedAt!: Date;

  @Column({ name: 'answered_at', type: 'datetime', nullable: true })
  answeredAt!: Date | null;

  @Column({ name: 'ended_at', type: 'datetime', nullable: true })
  endedAt!: Date | null;

  @Column({ name: 'duration_sec', type: 'int', default: 0 })
  durationSec!: number;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;
}
