// Read the internal token at call time so tests using vi.stubEnv can
// override it. SvelteKit 2 only inlines `import.meta.env.PUBLIC_*` at
// build time when explicit env declarations are configured; otherwise
// the access resolves to undefined and Rollup tree-shakes any
// conditional on it. The runtime container has PUBLIC_INTERNAL_TOKEN
// set in process.env, which SvelteKit exposes via $env/dynamic/public.
// We prefer that (static-safe) and fall back to import.meta.env so
// vi.stubEnv still works in tests.
import { env } from '$env/dynamic/public';
function getInternalToken(): string | undefined {
  return (env.PUBLIC_INTERNAL_TOKEN ?? import.meta.env.PUBLIC_INTERNAL_TOKEN) as string | undefined;
}

type ApiLogLevel = 'debug' | 'info' | 'silent';

function getApiLogLevel(): ApiLogLevel {
  const raw = (env.PUBLIC_LOG_LEVEL ?? import.meta.env.PUBLIC_LOG_LEVEL) as string | undefined;
  if (raw === 'debug' || raw === 'info') return raw;
  return 'silent';
}

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

export async function apiFetch<T = unknown>(
  path: string,
  init: RequestInit = {}
): Promise<T> {
  const headers = new Headers(init.headers);
  const token = getInternalToken();
  if (token) {
    headers.set('X-Internal-Token', token);
  }
  if (init.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const logLevel = getApiLogLevel();
  const method = init.method ?? 'GET';
  if (logLevel === 'debug') {
    console.debug('[api] request', { method, path });
  }
  const startedAt = Date.now();
  const resp = await fetch(path, { ...init, headers, credentials: 'omit' });
  if (logLevel === 'debug' || logLevel === 'info') {
    console.info('[api] response', {
      method,
      path,
      status: resp.status,
      ok: resp.ok,
      duration_ms: Date.now() - startedAt
    });
  }

  if (!resp.ok) {
    let msg = `${resp.status} ${resp.statusText}`;
    try {
      const body = (await resp.json()) as { error?: string };
      if (body?.error) msg = body.error;
    } catch {
      // body wasn't JSON; keep status text
    }
    throw new ApiError(resp.status, msg);
  }

  if (resp.status === 204) return undefined as T;
  return (await resp.json()) as T;
}

import type {
  Project, Server, Window_ as Window, RescaleEvent,
  CreateProjectRequest, CreateProjectResult,
  CreateServerRequest, UpdateServerRequest,
  CreateWindowRequest, RescaleRequest, ConfirmRequest,
  RefreshResponse, ServerType
} from './types';

export const api = {
  healthz: () => apiFetch<{ status: string }>('/api/healthz'),
  listProjects: () => apiFetch<Project[]>('/api/projects'),
  createProject: (body: CreateProjectRequest) =>
    apiFetch<CreateProjectResult>('/api/projects', { method: 'POST', body: JSON.stringify(body) }),
  deleteProject: (id: number) =>
    apiFetch<void>(`/api/projects/${id}`, { method: 'DELETE' }),
  refreshProject: (id: number) =>
    apiFetch<RefreshResponse>(`/api/projects/${id}/refresh`, { method: 'POST' }),

  listServers: () => apiFetch<Server[]>('/api/servers'),
  getServer: (id: number) => apiFetch<Server>(`/api/servers/${id}`),
  createServer: (body: CreateServerRequest) =>
    apiFetch<Server>('/api/servers', { method: 'POST', body: JSON.stringify(body) }),
  updateServer: (id: number, body: UpdateServerRequest) =>
    apiFetch<Server>(`/api/servers/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  deleteServer: (id: number) =>
    apiFetch<void>(`/api/servers/${id}`, { method: 'DELETE' }),

  listWindows: (serverId: number) =>
    apiFetch<Window[]>(`/api/servers/${serverId}/windows`),
  createWindow: (serverId: number, body: CreateWindowRequest) =>
    apiFetch<Window>(`/api/servers/${serverId}/windows`, { method: 'POST', body: JSON.stringify(body) }),
  updateWindow: (id: number, body: CreateWindowRequest) =>
    apiFetch<Window>(`/api/windows/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  deleteWindow: (id: number) =>
    apiFetch<void>(`/api/windows/${id}`, { method: 'DELETE' }),

  rescale: (id: number, body: RescaleRequest) =>
    apiFetch<{ status: string }>(`/api/servers/${id}/rescale`, { method: 'POST', body: JSON.stringify(body) }),
  promote: (id: number, body: ConfirmRequest) =>
    apiFetch<{ status: string }>(`/api/servers/${id}/promote`, { method: 'POST', body: JSON.stringify(body) }),
  demote: (id: number, body: ConfirmRequest) =>
    apiFetch<{ status: string }>(`/api/servers/${id}/demote`, { method: 'POST', body: JSON.stringify(body) }),

  serverEvents: (id: number, limit = 50) =>
    apiFetch<RescaleEvent[]>(`/api/servers/${id}/events?limit=${limit}`),
  globalEvents: (opts: { serverId?: number; limit?: number } = {}) => {
    const q = new URLSearchParams();
    if (opts.serverId) q.set('server_id', String(opts.serverId));
    if (opts.limit) q.set('limit', String(opts.limit));
    return apiFetch<RescaleEvent[]>(`/api/events?${q}`);
  },

  serverTypes: (location: string) =>
    apiFetch<ServerType[]>(`/api/server-types?location=${encodeURIComponent(location)}`),

  metrics: (range: '1d' | '7d' | '30d') =>
    apiFetch<import('./types').MetricsResponse>(`/api/metrics?range=${range}`)
};
