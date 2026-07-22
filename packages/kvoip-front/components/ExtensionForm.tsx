'use client';

import { useRouter } from 'next/navigation';
import { FormEvent, useState } from 'react';
import { api } from '@/lib/api';

export function ExtensionForm() {
  const router = useRouter();
  const [number, setNumber] = useState('');
  const [displayName, setDisplayName] = useState('');
  const [email, setEmail] = useState('');
  const [device, setDevice] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await api.extensions.create({
        number,
        displayName,
        email: email || undefined,
        device: device || undefined,
      });
      setNumber('');
      setDisplayName('');
      setEmail('');
      setDevice('');
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha ao criar ramal');
    } finally {
      setBusy(false);
    }
  }

  return (
    <form className="form-grid" onSubmit={onSubmit}>
      <label>
        Número
        <input
          required
          value={number}
          onChange={(e) => setNumber(e.target.value)}
          placeholder="1005"
        />
      </label>
      <label>
        Nome
        <input
          required
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          placeholder="Nome do ramal"
        />
      </label>
      <label>
        E-mail
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="opcional"
        />
      </label>
      <div className="form-actions" style={{ display: 'grid', gap: '0.5rem' }}>
        <label>
          Aparelho
          <input
            value={device}
            onChange={(e) => setDevice(e.target.value)}
            placeholder="Softphone"
          />
        </label>
        <button className="btn btn-primary" type="submit" disabled={busy}>
          {busy ? 'Salvando…' : 'Adicionar ramal'}
        </button>
      </div>
      {error ? (
        <div className="error-banner" style={{ gridColumn: '1 / -1', margin: 0 }}>
          {error}
        </div>
      ) : null}
    </form>
  );
}
