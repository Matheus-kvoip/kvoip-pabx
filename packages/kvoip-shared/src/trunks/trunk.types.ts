import type { TrunkProtocol, TrunkStatus } from './trunk.enums';

export type Trunk = {
  id: string;
  name: string;
  host: string;
  port: number;
  protocol: TrunkProtocol | `${TrunkProtocol}`;
  status: TrunkStatus | `${TrunkStatus}`;
  concurrentCalls: number;
  maxChannels: number;
};
