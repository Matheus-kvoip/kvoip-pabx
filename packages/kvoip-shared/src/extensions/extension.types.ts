import type { ExtensionStatus } from './extension.enums';

export type Extension = {
  id: string;
  number: string;
  displayName: string;
  email?: string;
  status: ExtensionStatus | `${ExtensionStatus}`;
  device?: string;
  createdAt: string;
};

export type CreateExtensionInput = {
  number: string;
  displayName: string;
  email?: string;
  device?: string;
  /** Senha SIP Digest (default no seed: kvoip123). */
  sipPassword?: string;
};
