import type { Metadata } from 'next';
import { JetBrains_Mono, Montserrat, Source_Sans_3 } from 'next/font/google';
import { appConfig } from '@/lib/config';
import './globals.css';

const montserrat = Montserrat({
  subsets: ['latin'],
  variable: '--font-montserrat',
});

const sourceSans = Source_Sans_3({
  subsets: ['latin'],
  variable: '--font-source',
});

const jetbrains = JetBrains_Mono({
  subsets: ['latin'],
  variable: '--font-jetbrains',
});

export const metadata: Metadata = {
  title: `${appConfig.appName} | PABX Virtual`,
  description:
    'Painel do PABX Virtual KVOIP — ramais, troncos SIP, filas e chamadas.',
  icons: {
    icon: '/logo-kvoip.png',
    apple: '/logo-kvoip.png',
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="pt-BR"
      className={`${montserrat.variable} ${sourceSans.variable} ${jetbrains.variable}`}
    >
      <body>{children}</body>
    </html>
  );
}
