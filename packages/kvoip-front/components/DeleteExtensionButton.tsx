'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { api } from '@/lib/api';

export function DeleteExtensionButton({ id }: { id: string }) {
  const router = useRouter();
  const [busy, setBusy] = useState(false);

  async function onDelete() {
    if (!confirm('Remover este ramal?')) return;
    setBusy(true);
    try {
      await api.extensions.remove(id);
      router.refresh();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Falha ao remover');
    } finally {
      setBusy(false);
    }
  }

  return (
    <button
      type="button"
      className="btn btn-danger"
      onClick={onDelete}
      disabled={busy}
    >
      {busy ? '…' : 'Remover'}
    </button>
  );
}
