import { api } from '@/lib/api-server';
import { StatusBadge } from '@/components/StatusBadge';

export const dynamic = 'force-dynamic';

export default async function TrunksPage() {
  let error: string | null = null;
  let trunks: Awaited<ReturnType<typeof api.trunks.list>> = [];

  try {
    trunks = await api.trunks.list();
  } catch (err) {
    error =
      err instanceof Error ? err.message : 'Falha ao carregar troncos';
  }

  return (
    <>
      <div className="page-head">
        <div>
          <h1>Troncos SIP</h1>
          <p>Conectividade com operadoras, gateways e rotas 0800 / DID</p>
        </div>
      </div>

      {error ? <div className="error-banner">{error}</div> : null}

      <section className="panel">
        <div className="panel-head">
          <h2>Troncos da conta</h2>
        </div>
        <div className="table-wrap">
          {trunks.length === 0 ? (
            <div className="empty">Nenhum tronco configurado.</div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Nome</th>
                  <th>Host</th>
                  <th>Protocolo</th>
                  <th>Canais</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {trunks.map((trunk) => (
                  <tr key={trunk.id}>
                    <td>{trunk.name}</td>
                    <td className="mono">
                      {trunk.host}:{trunk.port}
                    </td>
                    <td className="mono">{trunk.protocol.toUpperCase()}</td>
                    <td className="mono">
                      {trunk.concurrentCalls}/{trunk.maxChannels}
                    </td>
                    <td>
                      <StatusBadge status={trunk.status} />
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
