import React from 'react';

interface StatItemProps {
  value: string;
  label: string;
}

const StatItem: React.FC<StatItemProps> = ({ value, label }) => (
  <div className="stat-item">
    <div className="stat-value">{value}</div>
    <div className="stat-label">{label}</div>
  </div>
);

interface StatsSectionProps {
  stats: {
    free: string;
    models: string;
    unlimited: string;
    fast: string;
  };
}

export const StatsSection: React.FC<StatsSectionProps> = ({ stats }) => {
  return (
    <section className="py-20 px-6 relative overflow-hidden">
      <div className="stats-bg-decoration"></div>
      <div className="container mx-auto max-w-6xl relative z-10">
        <div className="stats-grid">
          <StatItem value="100%" label={stats.free} />
          <StatItem value="GLM-4" label={stats.models} />
          <StatItem value="∞" label={stats.unlimited} />
          <StatItem value="Fast" label={stats.fast} />
        </div>
      </div>
    </section>
  );
};
