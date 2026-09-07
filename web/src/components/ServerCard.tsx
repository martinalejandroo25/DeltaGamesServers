import React from 'react';
import { GameServer } from '../services/api';
import { Play, Square, Terminal, Trash2, Cpu, HardDrive, Network } from 'lucide-react';

interface ServerCardProps {
  server: GameServer;
  onStart: (id: string) => void;
  onStop: (id: string) => void;
  onDelete: (id: string) => void;
  onSelectConsole: (server: GameServer) => void;
}

export const ServerCard: React.FC<ServerCardProps> = ({
  server,
  onStart,
  onStop,
  onDelete,
  onSelectConsole
}) => {
  const isRunning = server.status === 'running';

  return (
    <div className="glass-panel" style={{ padding: '24px', display: 'flex', flexDirection: 'column', gap: '16px', position: 'relative' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h3 style={{ fontSize: '1.1rem', fontWeight: 700, color: '#fff', marginBottom: '4px' }}>
            {server.name}
          </h3>
          <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)', fontFamily: 'var(--font-mono)' }}>
            ID: {server.id}
          </span>
        </div>

        <div className={`badge ${isRunning ? 'badge-running' : 'badge-stopped'}`}>
          <span className="status-dot"></span>
          {server.status}
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '8px', padding: '12px', background: 'rgba(0, 0, 0, 0.2)', borderRadius: '8px', border: '1px solid var(--border-glass)', fontSize: '0.8rem' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px', color: 'var(--text-muted)' }}>
          <Cpu size={14} color="var(--accent-cyan)" />
          <span>CPU: {isRunning ? '1.2%' : '0%'}</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px', color: 'var(--text-muted)' }}>
          <HardDrive size={14} color="var(--accent-primary)" />
          <span>RAM: {isRunning ? '480 MB' : '0 MB'}</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px', color: 'var(--text-muted)' }}>
          <Network size={14} color="var(--accent-green)" />
          <span>Port: {Object.values(server.portMap)[0] || 27015}</span>
        </div>
      </div>

      <div style={{ display: 'flex', gap: '8px', marginTop: 'auto', paddingTop: '8px' }}>
        {isRunning ? (
          <button className="btn btn-secondary" style={{ color: 'var(--accent-red)' }} onClick={() => onStop(server.id)}>
            <Square size={16} /> Detener
          </button>
        ) : (
          <button className="btn btn-primary" onClick={() => onStart(server.id)}>
            <Play size={16} /> Iniciar
          </button>
        )}

        <button className="btn btn-secondary" onClick={() => onSelectConsole(server)}>
          <Terminal size={16} color="var(--accent-cyan)" /> Consola
        </button>

        <button className="btn btn-secondary" style={{ marginLeft: 'auto', padding: '10px' }} onClick={() => onDelete(server.id)} title="Eliminar Servidor">
          <Trash2 size={16} color="var(--text-dim)" />
        </button>
      </div>
    </div>
  );
};
