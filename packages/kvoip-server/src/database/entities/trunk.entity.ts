import {
  Column,
  CreateDateColumn,
  Entity,
  PrimaryGeneratedColumn,
  UpdateDateColumn,
} from 'typeorm';

@Entity('trunks')
export class TrunkEntity {
  @PrimaryGeneratedColumn('uuid')
  id!: string;

  @Column()
  name!: string;

  @Column()
  host!: string;

  @Column({ type: 'int', default: 5060 })
  port!: number;

  @Column({ default: 'udp' })
  protocol!: string;

  @Column({ default: 'down' })
  status!: string;

  @Column({ name: 'concurrent_calls', type: 'int', default: 0 })
  concurrentCalls!: number;

  @Column({ name: 'max_channels', type: 'int', default: 30 })
  maxChannels!: number;

  @Column({ default: true })
  enabled!: boolean;

  @CreateDateColumn({ name: 'created_at' })
  createdAt!: Date;

  @UpdateDateColumn({ name: 'updated_at' })
  updatedAt!: Date;
}
