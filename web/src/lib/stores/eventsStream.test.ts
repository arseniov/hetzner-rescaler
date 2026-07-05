import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';

// In a Vitest environment, SvelteKit's $env/dynamic/public virtual
// module returns an object whose `env` getter is undefined, which makes
// `import { env } from '$env/dynamic/public'` throw at module load.
// Mock the module so the store can be imported in tests; values are
// forwarded to vi.stubEnv's underlying import.meta.env so individual
// tests can drive the token via vi.stubEnv.
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

type Listener = (ev: MessageEvent) => void;
class MockEventSource {
  url: string;
  readyState = 0; // CONNECTING
  static instances: MockEventSource[] = [];
  listeners: Listener[] = [];
  constructor(url: string) {
    this.url = url;
    MockEventSource.instances.push(this);
  }
  addEventListener(kind: string, fn: Listener) {
    if (kind === 'message') this.listeners.push(fn);
  }
  removeEventListener() {}
  close() {}
  open() {
    this.readyState = 1; // OPEN
  }
  emit(data: unknown) {
    this.listeners.forEach((l) => l(new MessageEvent('message', { data: JSON.stringify(data) })));
  }
}

beforeEach(() => {
  MockEventSource.instances = [];
  (globalThis as any).EventSource = MockEventSource;
  vi.stubEnv('PUBLIC_INTERNAL_TOKEN', 'test-token');
});

afterEach(() => {
  vi.resetModules();
  vi.unstubAllEnvs();
});

describe('eventsStream store', () => {
  it('subscribes on connect and emits events into the array', async () => {
    const { eventsStream } = await import('./eventsStream.svelte');
    eventsStream.connect();
    expect(MockEventSource.instances).toHaveLength(1);
    const es = MockEventSource.instances[0];

    // The token should be delivered via the query string because the
    // browser EventSource cannot set custom headers.
    expect(es.url).toBe('/api/events/stream?token=test-token');

    es.open();
    es.emit({ id: 1, kind: 'rescale_up' });
    es.emit({ id: 2, kind: 'rescale_down' });

    expect(eventsStream.events.length).toBe(2);
    // The store prepends new events, so the most recently emitted event
    // sits at index 0.
    expect(eventsStream.events[0].kind).toBe('rescale_down');
    expect(eventsStream.events[1].kind).toBe('rescale_up');
  });

  it('caps the event list to 100', async () => {
    const { eventsStream } = await import('./eventsStream.svelte');
    eventsStream.connect();
    const es = MockEventSource.instances[0];
    es.open();
    for (let i = 0; i < 150; i++) es.emit({ id: i, kind: 'rescale_up' });
    expect(eventsStream.events.length).toBeLessThanOrEqual(100);
  });
});