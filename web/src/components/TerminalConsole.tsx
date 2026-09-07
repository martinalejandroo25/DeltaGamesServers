import React, { useEffect, useRef } from 'react';
import { Terminal as XTerm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';
import { X, RefreshCw } from 'lucide-react';
import { GameServer } from '../services/api';

interface TerminalConsoleProps {
  server: GameServer;
  onClose: () => void;
}

export const TerminalConsole: React.FC<TerminalConsoleProps> = ({ server, onClose }) => {
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermInstance = useRef<XTerm | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);

  useEffect(() => {
    if (!terminalRef.current) return;

    // Initialize xterm
    const term = new XTerm({
      cursorBlink: true,
      fontFamily: 'JetBrains Mono, monospace',
      fontSize: 13,
      theme: {
        background: '#0a0d14',
        foreground: '#e5e7eb',
        cursor: '#3b82f6',
        selectionBackground: 'rgba(59, 130, 246, 0.3)'
      }
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    fitAddon.fit();

    xtermInstance.current = term;
    fitAddonRef.current = fitAddon;

    term.writeln(`\x1b[1;34m=== Conectando a Consola Live de DeltaGamesServers [${server.name}] ===\x1b[0m`);

    // Connect WebSocket
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/servers/${server.id}/ws`;
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      term.writeln('\x1b[1;32m[WS Connected] WebSocket streaming iniciado.\x1b[0m\r\n');
    };

    ws.onmessage = (event) => {
      term.writeln(event.data);
    };

    ws.onerror = () => {
      term.writeln('\x1b[1;31m[WS Error] Fallo en conexión con el socket del servidor.\x1b[0m');
    };

    ws.onclose = () => {
      term.writeln('\r\n\x1b[1;33m[WS Closed] Conexión cerrada.\x1b[0m');
    };

    const handleResize = () => fitAddon.fit();
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      ws.close();
      term.dispose();
    };
  }, [server]);

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
      <div className="glass-panel" style={{
        width: '100%',
        maxWidth: '900px',
        height: '600px',
        display: 'flex',
        flexDirection: 'column',
        overflow: 'hidden',
        boxShadow: '0 20px 50px rgba(0,0,0,0.8)'
      }}>
        {/* Modal Header */}
        <div style={{
          padding: '14px 20px',
          background: 'rgba(255,255,255,0.03)',
          borderBottom: '1px solid var(--border-glass)',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <span className="status-dot" style={{ color: server.status === 'running' ? 'var(--accent-green)' : 'var(--accent-red)' }}></span>
            <strong style={{ color: '#fff', fontSize: '0.95rem' }}>Consola STDOUT/STDIN: {server.name}</strong>
          </div>

          <div style={{ display: 'flex', gap: '8px' }}>
            <button className="btn btn-secondary" style={{ padding: '6px' }} onClick={() => fitAddonRef.current?.fit()}>
              <RefreshCw size={16} />
            </button>
            <button className="btn btn-secondary" style={{ padding: '6px' }} onClick={onClose}>
              <X size={18} />
            </button>
          </div>
        </div>

        {/* Xterm Box */}
        <div ref={terminalRef} style={{ flex: 1, padding: '12px', background: '#0a0d14' }} />
      </div>
    </div>
  );
};
