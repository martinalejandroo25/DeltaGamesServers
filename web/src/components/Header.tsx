import React from 'react';
import { Server, ShieldCheck, Plus, Activity } from 'lucide-react';

interface HeaderProps {
  onOpenCreateModal: () => void;
  serverCount: number;
}

export const Header: React.FC<HeaderProps> = ({ onOpenCreateModal, serverCount }) => {
  return (
    <header style={{
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'space-between',
      padding: '16px 32px',
      background: 'rgba(18, 25, 41, 0.85)',
      backdropFilter: 'blur(12px)',
      borderBottom: '1px solid var(--border-glass)',
      position: 'sticky',
      top: 0,
      zIndex: 100
    }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
        <div style={{
          width: '42px',
          height: '42px',
          borderRadius: '12px',
          background: 'linear-gradient(135deg, var(--accent-primary), var(--accent-cyan))',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: 'var(--shadow-glow)'
        }}>
          <Server size={24} color="#fff" />
        </div>
        <div>
          <h1 style={{ fontSize: '1.25rem', fontWeight: 800, letterSpacing: '-0.02em', background: 'linear-gradient(90deg, #fff, #9ca3af)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
            DeltaGamesServers
          </h1>
          <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '0.75rem', color: 'var(--text-muted)' }}>
            <ShieldCheck size={14} color="var(--accent-green)" />
            <span>Podman Rootless Engine Active</span>
          </div>
        </div>
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px', padding: '6px 12px', background: 'rgba(255, 255, 255, 0.04)', borderRadius: '8px', border: '1px solid var(--border-glass)', fontSize: '0.85rem' }}>
          <Activity size={16} color="var(--accent-cyan)" />
          <span>Servidores Activos: <strong>{serverCount}</strong></span>
        </div>

        <button className="btn btn-primary" onClick={onOpenCreateModal}>
          <Plus size={18} />
          Nuevo Servidor
        </button>
      </div>
    </header>
  );
};
