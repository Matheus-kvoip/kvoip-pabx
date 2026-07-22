'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';
import { BrandLogo } from '@/components/BrandLogo';
import { api } from '@/lib/api';
import { clearClientToken } from '@/lib/auth';
import { appConfig } from '@/lib/config';

const links = [
  { href: '/', label: 'Visão geral' },
  { href: '/extensions', label: 'Ramais' },
  { href: '/trunks', label: 'Troncos SIP' },
  { href: '/calls', label: 'Chamadas' },
];

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [apiOk, setApiOk] = useState<boolean | null>(null);
  const [userName, setUserName] = useState<string | null>(null);

  useEffect(() => {
    let alive = true;
    api
      .health()
      .then(() => alive && setApiOk(true))
      .catch(() => alive && setApiOk(false));
    api
      .auth.me()
      .then((user) => alive && setUserName(user.name))
      .catch(() => alive && setUserName(null));
    return () => {
      alive = false;
    };
  }, []);

  async function onLogout() {
    try {
      await api.auth.logout();
    } catch {
      // ignore network errors on logout
    }
    clearClientToken();
    router.replace('/login');
    router.refresh();
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <BrandLogo size="sm" priority showName />
          <div className="brand-sub">{appConfig.appTagline}</div>
        </div>

        <nav className="nav">
          {links.map((link) => {
            const active =
              link.href === '/'
                ? pathname === '/'
                : pathname.startsWith(link.href);
            return (
              <Link
                key={link.href}
                href={link.href}
                className={active ? 'active' : undefined}
              >
                {link.label}
              </Link>
            );
          })}
        </nav>

        <div className="sidebar-foot">
          {userName ? <div style={{ fontWeight: 700, color: 'var(--ink)' }}>{userName}</div> : null}
          <div>
            <span className={`health-dot ${apiOk ? 'ok' : ''}`} />
            {apiOk === null
              ? 'Checando API…'
              : apiOk
                ? 'Plataforma online'
                : 'API indisponível'}
          </div>
          <div className="mono">{appConfig.apiHostLabel}</div>
          <button type="button" className="btn logout-btn" onClick={onLogout}>
            Sair
          </button>
        </div>
      </aside>

      <main className="main">{children}</main>
    </div>
  );
}
