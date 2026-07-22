export function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function formatDateTime(iso: string): string {
  return new Intl.DateTimeFormat('pt-BR', {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(iso));
}

export function statusLabel(
  status: string,
): string {
  const map: Record<string, string> = {
    online: 'Online',
    offline: 'Offline',
    busy: 'Ocupado',
    ringing: 'Tocando',
    up: 'Ativo',
    down: 'Fora',
    degraded: 'Degradado',
    inbound: 'Entrada',
    outbound: 'Saída',
    internal: 'Interna',
    answered: 'Em atendimento',
    held: 'Em espera',
    ended: 'Encerrada',
  };
  return map[status] ?? status;
}
