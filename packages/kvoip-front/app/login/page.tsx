import { Suspense } from 'react';
import { BrandLogo } from '@/components/BrandLogo';
import { appConfig } from '@/lib/config';
import { LoginForm } from './login-form';

export default function LoginPage() {
  return (
    <div className="login-screen">
      <div className="login-stage">
        <div className="login-brand">
          <BrandLogo size="lg" priority showName />
          <p className="brand-sub">{appConfig.appTagline}</p>
          <p className="login-lead">
            Acesse o painel do PABX Virtual: ramais, troncos SIP, filas e
            chamadas em um só lugar.
          </p>
          <div className="login-trust product-strip">
            <span className="product-chip">PABX Virtual</span>
            <span className="product-chip">0800 / DID</span>
            <span className="product-chip">Licenciada Anatel</span>
          </div>
        </div>

        <div className="login-card">
          <h2 className="login-card-title">Área do cliente</h2>
          <Suspense fallback={<div className="empty">Carregando…</div>}>
            <LoginForm />
          </Suspense>
          <div className="login-support">
            Suporte: <a href="tel:1140404838">4040-4838</a>
            {' · '}
            <a href="tel:08004444838">0800-444-4838</a>
            <br />
            <a href="mailto:equipe@kvoip.com.br">equipe@kvoip.com.br</a>
          </div>
        </div>
      </div>
    </div>
  );
}
