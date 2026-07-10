import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// SvelteKit's $env/dynamic/public is a virtual module not provided by Vitest.
// Mock it so api.ts can import `env` without throwing.
vi.mock('$env/dynamic/public', () => ({ env: {} }));

import { apiFetch } from './api';
import { api } from './api';

describe('apiFetch', () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllEnvs();
  });
  afterEach(() => {
    vi.unstubAllEnvs();
  });

  it('attaches X-Internal-Token header', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ ok: true }), { status: 200, headers: { 'Content-Type': 'application/json' } })
    );
    vi.stubGlobal('fetch', fetchMock);
    vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 'test-token');

    await apiFetch('/api/healthz', { method: 'GET' });

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe('/api/healthz');
    // Headers is a Headers instance, not a plain object; use .get().
    const headers = init?.headers as Headers;
    expect(headers.get('X-Internal-Token')).toBe('test-token');
  });

  it('throws on non-2xx with parsed JSON error body', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ error: 'unauthorized' }), { status: 401 })
    ));
    vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 't');

    await expect(apiFetch('/api/projects')).rejects.toThrow(/unauthorized/);
  });

  it('returns parsed JSON on 2xx', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify([{ id: 1, name: 'p' }]), { status: 200 })
    ));
    vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 't');

    const result = await apiFetch('/api/projects');
    expect(result).toEqual([{ id: 1, name: 'p' }]);
  });

  it('logs request and response details at debug level', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ ok: true }), { status: 200 })
    ));
    vi.stubEnv('PUBLIC_LOG_LEVEL', 'debug');
    const debugSpy = vi.spyOn(console, 'debug').mockImplementation(() => undefined);
    const infoSpy = vi.spyOn(console, 'info').mockImplementation(() => undefined);

    await apiFetch('/api/projects', { method: 'POST', body: JSON.stringify({ name: 'p' }) });

    expect(debugSpy).toHaveBeenCalledWith('[api] request', {
      method: 'POST',
      path: '/api/projects'
    });
    expect(infoSpy).toHaveBeenCalledWith('[api] response', expect.objectContaining({
      method: 'POST',
      path: '/api/projects',
      status: 200,
      ok: true,
      duration_ms: expect.any(Number)
    }));
  });

  it('logs response details at info level without debug request logs', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ ok: true }), { status: 201 })
    ));
    vi.stubEnv('PUBLIC_LOG_LEVEL', 'info');
    const debugSpy = vi.spyOn(console, 'debug').mockImplementation(() => undefined);
    const infoSpy = vi.spyOn(console, 'info').mockImplementation(() => undefined);

    await apiFetch('/api/projects', { method: 'POST', body: JSON.stringify({ name: 'p' }) });

    expect(debugSpy).not.toHaveBeenCalled();
    expect(infoSpy).toHaveBeenCalledWith('[api] response', expect.objectContaining({
      method: 'POST',
      path: '/api/projects',
      status: 201,
      ok: true,
      duration_ms: expect.any(Number)
    }));
  });
});

describe('api.serverTypes', () => {
  beforeEach(() => vi.restoreAllMocks());
  afterEach(() => vi.restoreAllMocks());

  it('hits /api/server-types?location=fsn1 with the auth header', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify([{ name: 'cpx11', available: true }]), {
        status: 200, headers: { 'Content-Type': 'application/json' }
      })
    );
    vi.stubGlobal('fetch', fetchMock);
    vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 't');

    const result = await api.serverTypes('fsn1');

    expect(fetchMock).toHaveBeenCalledTimes(1);
    const [url, init] = fetchMock.mock.calls[0];
    expect(url).toBe('/api/server-types?location=fsn1');
    expect((init?.headers as Headers).get('X-Internal-Token')).toBe('t');
    expect(result).toEqual([{ name: 'cpx11', available: true }]);
  });

  it('encodes special characters in the location', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify([]), { status: 200 })
    );
    vi.stubGlobal('fetch', fetchMock);
    vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 't');

    await api.serverTypes('fsn1 dc');

    const [url] = fetchMock.mock.calls[0];
    expect(url).toBe('/api/server-types?location=fsn1%20dc');
  });
});
