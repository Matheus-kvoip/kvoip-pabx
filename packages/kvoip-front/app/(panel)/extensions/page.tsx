import { api } from '@/lib/api-server';
import { formatDateTime } from '@/lib/format';
import { ExtensionForm } from '@/components/ExtensionForm';
import { DeleteExtensionButton } from '@/components/DeleteExtensionButton';
import { StatusBadge } from '@/components/StatusBadge';

export const dynamic = 'force-dynamic';

export default async function ExtensionsPage() {
  let error: string | null = null;
  let extensions: Awaited<ReturnType<typeof api.extensions.list>> = [];

  try {
    extensions = await api.extensions.list();
  } catch (err) {
    error =
      err instanceof Error ? err.message : 'Falha ao carregar ramais';
  }

  return (
    <>
      <div className="page-head">
        <div>
          <h1>Ramais</h1>
          <p>Endpoints SIP do PABX Virtual — softphone, IP phone ou fila</p>
        </div>
      </div>

      {error ? <div className="error-banner">{error}</div> : null}

      <section className="panel">
        <div className="panel-head">
          <h2>Cadastrar ramal</h2>
        </div>
        <ExtensionForm />
        <div className="table-wrap">
          {extensions.length === 0 ? (
            <div className="empty">Nenhum ramal cadastrado.</div>
          ) : (
            <table>
              <thead>
                <tr>
                  <th>Número</th>
                  <th>Nome</th>
                  <th>Aparelho</th>
                  <th>Status</th>
                  <th>Criado</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {extensions.map((ext) => (
                  <tr key={ext.id}>
                    <td className="mono">{ext.number}</td>
                    <td>
                      <div>{ext.displayName}</div>
                      {ext.email ? (
                        <div style={{ color: 'var(--muted)', fontSize: '0.8rem' }}>
                          {ext.email}
                        </div>
                      ) : null}
                    </td>
                    <td>{ext.device ?? '—'}</td>
                    <td>
                      <StatusBadge status={ext.status} />
                    </td>
                    <td className="mono">{formatDateTime(ext.createdAt)}</td>
                    <td>
                      <DeleteExtensionButton id={ext.id} />
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
