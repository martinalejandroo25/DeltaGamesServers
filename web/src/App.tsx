import React, { useState, useEffect } from 'react';
import { Header } from './components/Header';
import { ServerCard } from './components/ServerCard';
import { TerminalConsole } from './components/TerminalConsole';
import { CreateServerModal } from './components/CreateServerModal';
import { GameServer, fetchServers, startServer, stopServer, deleteServer } from './services/api';
import { Server, RefreshCw } from 'lucide-react';

export const App: React.FC = () => {
  const [servers, setServers] = useState<GameServer[]>([]);
  const [selectedServerConsole, setSelectedServerConsole] = useState<GameServer | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);

  const loadServers = async () => {
    try {
      const data = await fetchServers();
      setServers(data || []);
    } catch (err) {
      console.error('Failed to load servers:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadServers();
    const interval = setInterval(loadServers, 5000);
    return () => clearInterval(interval);
  }, []);

  const handleStart = async (id: string) => {
    await startServer(id);
    loadServers();
  };

  const handleStop = async (id: string) => {
    await stopServer(id);
    loadServers();
  };

  const handleDelete = async (id: string) => {
    if (confirm('¿Seguro que deseas eliminar este servidor y su contenedor Podman?')) {
      await deleteServer(id);
      loadServers();
    }
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Header
        serverCount={servers.filter((s) => s.status === 'running').length}
        onOpenCreateModal={() => setIsCreateModalOpen(true)}
      />

      <main style={{ flex: 1, padding: '32px', maxWidth: '1400px', margin: '0 auto', width: '100%' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '24px' }}>
          <div>
            <h2 style={{ fontSize: '1.4rem', fontWeight: 800, color: '#fff' }}>Instancias de Servidores</h2>
            <p style={{ fontSize: '0.85rem', color: 'var(--text-muted)' }}>
              Gestión automatizada de contenedores Podman en tu Home Server
            </p>
          </div>

          <button className="btn btn-secondary" onClick={loadServers}>
            <RefreshCw size={16} /> Actualizar
          </button>
        </div>

        {loading ? (
          <div style={{ textAlign: 'center', padding: '60px', color: 'var(--text-muted)' }}>
            Cargando instancias de servidores...
          </div>
        ) : servers.length === 0 ? (
          <div className="glass-panel" style={{ textAlign: 'center', padding: '60px 24px', display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '16px' }}>
            <div style={{ width: '60px', height: '60px', borderRadius: '50%', background: 'rgba(255, 255, 255, 0.05)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              <Server size={30} color="var(--text-dim)" />
            </div>
            <h3 style={{ color: '#fff', fontSize: '1.2rem' }}>No hay servidores creados</h3>
            <p style={{ color: 'var(--text-muted)', maxWidth: '400px', fontSize: '0.9rem' }}>
              Despliega tu primer servidor de juegos dedicado en un contenedor Podman en pocos segundos.
            </p>
            <button className="btn btn-primary" onClick={() => setIsCreateModalOpen(true)}>
              Desplegar Servidor
            </button>
          </div>
        ) : (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(340px, 1fr))', gap: '20px' }}>
            {servers.map((server) => (
              <ServerCard
                key={server.id}
                server={server}
                onStart={handleStart}
                onStop={handleStop}
                onDelete={handleDelete}
                onSelectConsole={(srv) => setSelectedServerConsole(srv)}
              />
            ))}
          </div>
        )}
      </main>

      {/* Modals */}
      {isCreateModalOpen && (
        <CreateServerModal
          onClose={() => setIsCreateModalOpen(false)}
          onCreated={loadServers}
        />
      )}

      {selectedServerConsole && (
        <TerminalConsole
          server={selectedServerConsole}
          onClose={() => setSelectedServerConsole(null)}
        />
      )}
    </div>
  );
};
