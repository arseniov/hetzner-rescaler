import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// Mock $env/dynamic/public so api.ts can import without throwing in
// Vitest. Same Proxy pattern used by eventsStream.test.ts.
vi.mock('$env/dynamic/public', () => ({
  env: new Proxy(
    {},
    {
      get: (_target, prop) => {
        const viteEnv = (import.meta as any).env ?? {};
        return viteEnv[prop as string];
      }
    }
  )
}));

const SAMPLE_TYPES = [
  {
    name: 'cpx11',
    description: 'Intel/AMD shared',
    cores: 2,
    memory_gb: 2,
    disk_gb: 40,
    available: true,
    price_monthly_eur: 3.29
  },
  {
    name: 'cpx31',
    description: 'Intel/AMD dedicated',
    cores: 4,
    memory_gb: 8,
    disk_gb: 80,
    available: false,
    price_monthly_eur: 0
  }
];

function mockFetchOnce(body: unknown, status = 200) {
  (globalThis as any).fetch = vi.fn().mockResolvedValue({
    ok: status >= 200 && status < 300,
    status,
    statusText: status === 200 ? 'OK' : 'ERR',
    json: async () => body
  });
}

beforeEach(() => {
  vi.resetModules();
  vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 'test-token');
});

afterEach(() => {
  vi.unstubAllEnvs();
  vi.restoreAllMocks();
});

describe('serverTypes store', () => {
  it('loads from /api/server-types and caches', async () => {
    mockFetchOnce(SAMPLE_TYPES);
    const { serverTypes } = await import('./serverTypes.svelte');
    serverTypes.reset();
    await serverTypes.load();
    expect(serverTypes.types).toHaveLength(2);
    expect(serverTypes.loadedAt).not.toBeNull();
  });

  it('is idempotent within the soft TTL', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      statusText: 'OK',
      json: async () => SAMPLE_TYPES
    });
    (globalThis as any).fetch = fetchMock;
    const { serverTypes } = await import('./serverTypes.svelte');
    serverTypes.reset();
    await serverTypes.load();
    await serverTypes.load();
    await serverTypes.load();
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it('byName returns the matching type', async () => {
    mockFetchOnce(SAMPLE_TYPES);
    const { serverTypes } = await import('./serverTypes.svelte');
    serverTypes.reset();
    await serverTypes.load();
    expect(serverTypes.byName('cpx11')?.cores).toBe(2);
    expect(serverTypes.byName('missing')).toBeUndefined();
  });

  it('available filters out sold-out types', async () => {
    mockFetchOnce(SAMPLE_TYPES);
    const { serverTypes } = await import('./serverTypes.svelte');
    serverTypes.reset();
    await serverTypes.load();
    const avail = serverTypes.available();
    expect(avail.map((t) => t.name)).toEqual(['cpx11']);
  });

  it('records loadError on fetch failure', async () => {
    mockFetchOnce(SAMPLE_TYPES);
    const { serverTypes } = await import('./serverTypes.svelte');
    serverTypes.reset();
    await serverTypes.load();
    expect(serverTypes.types).toHaveLength(2);
    // Force a re-fetch; the mocked fetch now fails.
    serverTypes.reset();
    (globalThis as any).fetch = vi.fn().mockRejectedValue(new Error('boom'));
    await expect(serverTypes.load(true)).rejects.toThrow('boom');
    // After reset+failed load, types is empty and loadError is set.
    expect(serverTypes.types).toHaveLength(0);
    expect(serverTypes.loadError).toBe('boom');
  });
});