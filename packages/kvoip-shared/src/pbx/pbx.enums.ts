/** Métodos SIP suportados pelo núcleo. */
export enum SipMethod {
  Invite = 'INVITE',
  Ack = 'ACK',
  Bye = 'BYE',
  Cancel = 'CANCEL',
  Options = 'OPTIONS',
  Register = 'REGISTER',
  Notify = 'NOTIFY',
}

/** Realm Digest padrão do ambiente local. */
export enum SipAuthRealm {
  Local = 'kvoip.local',
}

/** Chaves de ambiente SIP compartilhadas (docs / validação). */
export enum SipEnvKey {
  AuthEnabled = 'SIP_AUTH_ENABLED',
  AuthRealm = 'SIP_AUTH_REALM',
  Users = 'SIP_USERS',
  BindHost = 'SIP_BIND_HOST',
  AdvertisedHost = 'SIP_ADVERTISED_HOST',
  BufferSize = 'SIP_BUFFER_SIZE',
  Port = 'PORT_SERVER_SIP',
}
