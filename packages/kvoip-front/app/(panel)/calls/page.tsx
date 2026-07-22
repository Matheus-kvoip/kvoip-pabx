import { api } from '@/lib/api-server';
import { formatDateTime, formatDuration, statusLabel } from '@/lib/format';
import { StatusBadge } from '@/components/StatusBadge';

export const dynamic = 'force-dynamic';

export default async function CallsPage() {
  let error: string | null = null;
  let calls: Awaited<ReturnType<typeof api.calls.list>> = [];

  try {
    calls = await api.calls.list();
  } catch (err) {
    error =
      err instanceof Error ? err.message : 'Falha ao carregar chamadas';
  }

  return (
    <>
      <div className="page-head">
        <div>
          <h1>Chamadas</h1>
          <p>Histórico e sessões ativas — base para gravação e CDR</p>
        </div>
      </div>

      {error ? <div className="error-banner">{error}</div> : null}

      <section className="panel">
        <div className="panel-head">
          <h2>Registro de chamadas</h2>
        </div>
        <div className="table-wrap">
          {calls.length === 0 ? (
            <div className="empty">Sem registros de chamada.</div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Início</th>
                  <th>De</th>
                  <th>Para</th>
                  <th>Direção</th>
                  <th>Estado</th>
                  <th>Duração</th>
                </tr>
              </thead>
              <tbody>
                {calls.map((call) => (
                  <tr key={call.id}>
                    <td className="mono">{formatDateTime(call.startedAt)}</td>
                    <td className="mono">{call.from}</td>
                    <td className="mono">{call.to}</td>
                    <td>{statusLabel(call.direction)}</td>
                    <td>
                      <StatusBadge status={call.state} />
                    </td>
                    <td className="mono">
                      {formatDuration(call.durationSec)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </section>
    </>
  );
}
