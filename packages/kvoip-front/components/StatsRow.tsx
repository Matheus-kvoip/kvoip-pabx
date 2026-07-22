import type { DashboardStats } from '@kvoip/shared';

export function StatsRow({ stats }: { stats: DashboardStats }) {
  const items = [
    {
      label: 'Chamadas ativas',
      value: String(stats.activeCalls),
      meta: `${stats.callsToday} no período`,
    },
    {
      label: 'Ramais online',
      value: `${stats.extensionsOnline}`,
      meta: `de ${stats.extensionsTotal} cadastrados`,
    },
    {
      label: 'Troncos SIP',
      value: `${stats.trunksUp}`,
      meta: `ativos de ${stats.trunksTotal}`,
    },
    {
      label: 'Núcleo PABX',
      value: stats.pbxOnline ? 'Online' : 'Offline',
      meta: stats.pbxOnline
        ? 'dados ao vivo do kvoip-pbx'
        : 'usando fallback / mock',
    },
  ];

  return (
    <div className="stats">
      {items.map((item) => (
        <div key={item.label} className="stat">
          <div className="stat-label">{item.label}</div>
          <div className="stat-value">{item.value}</div>
          <div className="stat-meta">{item.meta}</div>
        </div>
      ))}
    </div>
  );
}
