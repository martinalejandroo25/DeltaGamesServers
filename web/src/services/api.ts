export interface GameServer {
  id: string;
  name: string;
  gameId: string;
  containerId: string;
  status: 'running' | 'stopped' | 'starting' | 'error';
  portMap: Record<string, number>;
  envVars: Record<string, string>;
  volumePath: string;
  autoStart: boolean;
  createdAt: string;
}

export interface GameTemplate {
  id: string;
  name: string;
  steamAppId: string;
  baseImage: string;
  startBinary: string;
  defaultParams: string;
  environment: Record<string, string>;
  ports: Array<{ containerPort: number; protocol: string; description: string }>;
  volumes: Array<{ containerPath: number; description: string }>;
}

const API_BASE = '/api/v1';

export async function fetchHealth() {
  const res = await fetch(`${API_BASE}/health`);
  return res.json();
}

export async function fetchTemplates(): Promise<GameTemplate[]> {
  const res = await fetch(`${API_BASE}/templates`);
  if (!res.ok) throw new Error('Failed to fetch templates');
  return res.json();
}

export async function fetchServers(): Promise<GameServer[]> {
  const res = await fetch(`${API_BASE}/servers`);
  if (!res.ok) throw new Error('Failed to fetch game servers');
  return res.json();
}

export async function createServer(payload: {
  name: string;
  gameId: string;
  customPorts?: Record<string, number>;
  customEnv?: Record<string, string>;
}): Promise<GameServer> {
  const res = await fetch(`${API_BASE}/servers`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });
  if (!res.ok) throw new Error('Failed to create server');
  return res.json();
}

export async function startServer(id: string) {
  const res = await fetch(`${API_BASE}/servers/${id}/start`, { method: 'POST' });
  if (!res.ok) throw new Error('Failed to start server');
  return res.json();
}

export async function stopServer(id: string) {
  const res = await fetch(`${API_BASE}/servers/${id}/stop`, { method: 'POST' });
  if (!res.ok) throw new Error('Failed to stop server');
  return res.json();
}

export async function deleteServer(id: string) {
  const res = await fetch(`${API_BASE}/servers/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error('Failed to delete server');
  return res.json();
}

export async function fetchLogs(id: string, tail: number = 100) {
  const res = await fetch(`${API_BASE}/servers/${id}/logs?tail=${tail}`);
  if (!res.ok) throw new Error('Failed to fetch logs');
  return res.json();
}
