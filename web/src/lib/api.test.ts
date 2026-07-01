import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { apiFetch } from './api';

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
});
