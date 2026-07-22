/** Estado do tronco SIP. */
export enum TrunkStatus {
  Up = 'up',
  Down = 'down',
  Degraded = 'degraded',
}

/** Transporte SIP do tronco. */
export enum TrunkProtocol {
  Udp = 'udp',
  Tcp = 'tcp',
  Tls = 'tls',
}
