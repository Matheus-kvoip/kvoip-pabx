'use client';

import { FormEvent, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { api } from '@/lib/api';
import { setClientToken } from '@/lib/auth';

export function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [email, setEmail] = useState('admin@kvoip.com.br');
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setBusy(true);
    setError(null);

    try {
      const result = await api.auth.login({ email, password });
      setClientToken(result.accessToken);
      const next = searchParams.get('next') || '/';
      router.replace(next);
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Falha no login');
    } finally {
      setBusy(false);
    }
  }

  return (
    <form className="login-form" onSubmit={onSubmit}>
      <label>
        E-mail
        <input
          type="email"
          autoComplete="username"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
        />
      </label>
      <label>
        Senha
        <input
          type="password"
          autoComplete="current-password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="••••••••"
        />
      </label>

      {error ? <div className="error-banner">{error}</div> : null}

      <button className="btn btn-primary login-submit" type="submit" disabled={busy}>
        {busy ? 'Entrando…' : 'Entrar'}
      </button>

      <p className="login-hint">
        Ambiente de demonstração. Credenciais no <span className="mono">.env</span> do
        servidor.
      </p>
    </form>
  );
}
