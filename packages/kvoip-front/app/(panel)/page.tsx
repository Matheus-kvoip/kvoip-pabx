import { api } from '@/lib/api-server';
import { formatDateTime, formatDuration, statusLabel } from '@/lib/format';
import { StatsRow } from '@/components/StatsRow';
import { StatusBadge } from '@/components/StatusBadge';
import Link from 'next/link';

export const dynamic = 'force-dynamic';

export default async function HomePage() {
  let error: string | null = null;
  let stats = null;
  let activeCalls: Awaited<ReturnType<typeof api.calls.list>> = [];
  let extensions: Awaited<ReturnType<typeof api.extensions.list>> = [];

  try {
    [stats, activeCalls, extensions] = await Promise.all([
      api.dashboard(),
      api.calls.list(true),
      api.extensions.list(),
    ]);
  } catch (err) {
    error =
      err instanceof Error
        ? err.message
        : 'Não foi possível carregar o painel. Suba a API em :3001.';
  }

  return (
    <>
      <div className="page-head">
        <div>
          <h1>PABX Virtual</h1>
          <p>Central telefônica na nuvem — operação em tempo real</p>
        </div>
        <Link className="btn btn-primary" href="/extensions">
          Gerenciar ramais
        </Link>
      </div>

      <div className="product-strip">
        <span className="product-chip">Ramais</span>
        <span className="product-chip">Filas / URA</span>
        <span className="product-chip">Troncos SIP</span>
        <span className="product-chip">Gravação</span>
        <span className="product-chip">0800 &amp; DID</span>
      </div>

      {error ? <div className="error-banner">{error}</div> : null}

      {stats ? <StatsRow stats={stats} /> : null}

      <div className="grid-2">
        <section className="panel">
          <div className="panel-head">
            <h2>Chamadas em andamento</h2>
            <Link href="/calls">Ver histórico</Link>
          </div>
          <div className="table-wrap">
            {activeCalls.length === 0 ? (
              <div className="empty">Nenhuma chamada ativa no momento.</div>
            ) : (
              <table>
                <thead>
                  <tr>
                    <th>De</th>
                    <th>Para</th>
                    <th>Direção</th>
                    <th>Estado</th>
                    <th>Duração</th>
                  </tr>
                </thead>
                <tbody>
                  {activeCalls.map((call) => (
                    <tr key={call.id}>
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

        <section className="panel">
          <div className="panel-head">
            <h2>Ramais do PABX</h2>
            <Link href="/extensions">Abrir</Link>
          </div>
          <div className="table-wrap">
            {extensions.length === 0 ? (
              <div className="empty">Sem ramais cadastrados.</div>
            ) : (
              <table>
                <thead>
                  <tr>
                    <th>Nº</th>
                    <th>Nome</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {extensions.slice(0, 6).map((ext) => (
                    <tr key={ext.id}>
                      <td className="mono">{ext.number}</td>
                      <td>{ext.displayName}</td>
                      <td>
                        <StatusBadge status={ext.status} />
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </section>
      </div>

      {stats ? (
        <p style={{ marginTop: '1rem', color: 'var(--muted)', fontSize: '0.85rem' }}>
          Atualizado em {formatDateTime(new Date().toISOString())} · KVOIP Brasil
          Telecom
        </p>
      ) : null}
    </>
  );
}
