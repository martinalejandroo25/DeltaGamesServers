import React, { useState, useEffect } from 'react';
import { GameTemplate, fetchTemplates, createServer } from '../services/api';
import { X, Gamepad2, Layers, Cpu } from 'lucide-react';

interface CreateServerModalProps {
  onClose: () => void;
  onCreated: () => void;
}

export const CreateServerModal: React.FC<CreateServerModalProps> = ({ onClose, onCreated }) => {
  const [templates, setTemplates] = useState<GameTemplate[]>([]);
  const [name, setName] = useState('');
  const [selectedGameId, setSelectedGameId] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchTemplates()
      .then((data) => {
        setTemplates(data);
        if (data.length > 0) setSelectedGameId(data[0].id);
      })
      .catch(console.error);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !selectedGameId) return;

    setLoading(true);
    try {
      await createServer({
        name,
        gameId: selectedGameId
      });
      onCreated();
      onClose();
    } catch (err) {
      alert('Error al crear el servidor: ' + err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      position: 'fixed',
      inset: 0,
      background: 'rgba(5, 8, 15, 0.85)',
      backdropFilter: 'blur(16px)',
      zIndex: 200,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '24px'
    }}>
      <div className="glass-panel" style={{ width: '100%', maxWidth: '540px', padding: '28px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <Gamepad2 size={24} color="var(--accent-primary)" />
            <h2 style={{ fontSize: '1.2rem', color: '#fff' }}>Desplegar Nuevo Servidor</h2>
          </div>
          <button className="btn btn-secondary" style={{ padding: '6px' }} onClick={onClose}>
            <X size={18} />
          </button>
        </div>

        <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
          <div>
            <label style={{ display: 'block', fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '6px' }}>
              Nombre del Servidor
            </label>
            <input
              type="text"
              required
              placeholder="Ej: Mi Servidor Garry's Mod Sandbox"
              value={name}
              onChange={(e) => setName(e.target.value)}
              style={{
                width: '100%',
                padding: '10px 14px',
                borderRadius: '8px',
                border: '1px solid var(--border-glass)',
                background: 'rgba(0, 0, 0, 0.3)',
                color: '#fff',
                fontFamily: 'var(--font-sans)',
                outline: 'none'
              }}
            />
          </div>

          <div>
            <label style={{ display: 'block', fontSize: '0.85rem', color: 'var(--text-muted)', marginBottom: '6px' }}>
              Plantilla de Juego (OCI Engine)
            </label>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              {templates.map((tmpl) => (
                <div
                  key={tmpl.id}
                  onClick={() => setSelectedGameId(tmpl.id)}
                  style={{
                    padding: '12px 16px',
                    borderRadius: '8px',
                    border: selectedGameId === tmpl.id ? '1px solid var(--accent-primary)' : '1px solid var(--border-glass)',
                    background: selectedGameId === tmpl.id ? 'rgba(59, 130, 246, 0.12)' : 'rgba(255, 255, 255, 0.02)',
                    cursor: 'pointer',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    transition: 'all 0.2s ease'
                  }}
                >
                  <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                    <Layers size={18} color="var(--accent-cyan)" />
                    <div>
                      <strong style={{ color: '#fff', fontSize: '0.9rem' }}>{tmpl.name}</strong>
                      <div style={{ fontSize: '0.75rem', color: 'var(--text-muted)' }}>Steam AppID: {tmpl.steamAppId}</div>
                    </div>
                  </div>
                  <Cpu size={16} color="var(--text-dim)" />
                </div>
              ))}
            </div>
          </div>

          <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '10px', marginTop: '12px' }}>
            <button type="button" className="btn btn-secondary" onClick={onClose}>
              Cancelar
            </button>
            <button type="submit" className="btn btn-primary" disabled={loading}>
              {loading ? 'Creando Contenedor...' : 'Crear y Desplegar'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
