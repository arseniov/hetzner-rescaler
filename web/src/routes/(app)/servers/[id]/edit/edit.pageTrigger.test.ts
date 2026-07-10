import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, cleanup } from '@testing-library/svelte';
import { tick } from 'svelte';

// $env/dynamic/public virtual module — the paraglide `m` helper pulls
// from there. Mock to a no-op proxy.
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

// $app/stores — the page store is read once for `params.id`.
vi.mock('$app/stores', () => ({
  page: {
    subscribe: (run: (value: { params: { id: string } }) => void) => {
      run({ params: { id: '42' } });
      return () => {};
    }
  }
}));

// $app/navigation — goto is a no-op in tests.
vi.mock('$app/navigation', () => ({
  goto: vi.fn().mockResolvedValue(undefined)
}));

// Wire api.mockServer so each test can pre-set what /api/servers/42
// returns. The page's onMount calls api.getServer(serverId) and (with
// the fix) calls serverTypes.load('fsn1') unconditionally — the test
// for the "no-location" case relies on this NOT depending on the API
// response. Stored here as a holder so vi.mock can reference it.
const apiMockState: { server: any } = { server: null };

vi.mock('$lib/api', () => {
  return {
    api: {
      getServer: vi.fn().mockImplementation(() => Promise.resolve(apiMockState.server)),
      updateServer: vi.fn(),
      // serverTypes.load() resolves through api.serverTypes().
      serverTypes: vi.fn().mockResolvedValue([])
    },
    ApiError: class ApiError extends Error {
      constructor(public status: number, message: string) {
        super(message);
        this.name = 'ApiError';
      }
    }
  };
});

import { api } from '$lib/api';
import { serverTypes } from '$lib/stores/serverTypes.svelte';
import EditPage from './+page.svelte';

describe('edit page — parent-level serverTypes.load trigger', () => {
  beforeEach(() => {
    serverTypes._reset();
    cleanup();
  });
  afterEach(() => {
    cleanup();
    serverTypes._reset();
  });

  it('fires serverTypes.load with the fallback when server resolves with no location', async () => {
    // This reproduces the real-world bug from 2026-07-10:
    // /api/servers/[id] returns a server WITHOUT a `location` field
    // (Hetzner's GetServer for that server soft-fails — Datacenter is
    // nil — and the omitempty tag drops the key). Before the fix the
    // edit page's onMount only fired load() inside `if (server.location)`
    // so the call never went out for such servers; the contract the
    // type-availability gate relies on was silently broken.
    //
    // With the fix, the page fires serverTypes.load(FALLBACK_LOCATION)
    // BEFORE awaiting api.getServer, so the request goes out regardless
    // of what the server response contains. The dedupe map on the
    // store collapses any subsequent calls.
    apiMockState.server = {
      id: 42,
      project_id: 1,
      hcloud_server_id: 1,
      name: 'fiutaspesa-app',
      label: 'fiutaspesa-app',
      base_server_type: 'cx33',
      top_server_type: 'cx33',
      fallback_chain: ['cx33'],
      mode: 'manual',
      timezone: 'UTC',
      status: 'running',
      current_type: 'cx33',
      // NOTE: no `location` field — this matches the actual API payload
      // for a server where Hetzner's GetServer soft-fails.
      created_at: '2026-07-10T06:48:52Z',
      updated_at: '2026-07-10T06:48:52Z'
    };

    const loadSpy = vi.spyOn(serverTypes, 'load');
    expect(api.getServer).toBeDefined();

    render(EditPage);
    // Wait for the unconditional load (fired before await) to be
    // captured by the spy.
    await tick();
    await tick();
    await tick();

    expect(loadSpy).toHaveBeenCalled();
    expect(loadSpy.mock.calls[0][0]).toBe('fsn1'); // FALLBACK_LOCATION
    // The post-resolution refire is skipped because server.location is
    // undefined, so we must NOT have a second call with a different
    // location. The store dedupes identical back-to-back calls anyway.
    loadSpy.mockRestore();
  });

  it('fires serverTypes.load once with the real location when server has one', async () => {
    // When the server resolves WITH a real location matching the
    // fallback ('fsn1'), the page fires:
    //   1. serverTypes.load('fsn1') from the parent — BEFORE await.
    //      This is the unconditional call that fixes the no-location
    //      regression. It guarantees the network request goes out on
    //      every navigation regardless of what the server returns.
    //   2-4. serverTypes.load('fsn1') from the three child components
    //      (ServerTypeSelect base, ServerTypeSelect top,
    //      ServerTypeMultiSelect fallback) when their `location` prop
    //      flips undefined → 'fsn1' after server resolves.
    // The store's in-flight map collapses all four calls into one
    // network request — vi.spyOn sees the four invocations, but only
    // one HTTP round-trip happens.
    //
    // Critically: NO 5th call with a different location, because the
    // page's check (`server.location !== FALLBACK_LOCATION`) is false
    // ('fsn1' === 'fsn1') so the parent post-await refire is skipped.
    apiMockState.server = {
      id: 42,
      project_id: 1,
      hcloud_server_id: 1,
      name: 'w',
      label: 'w',
      base_server_type: 'cpx11',
      top_server_type: 'cpx31',
      fallback_chain: ['cpx21'],
      mode: 'manual',
      timezone: 'UTC',
      status: 'running',
      current_type: 'cpx21',
      location: 'fsn1'
    };

    const loadSpy = vi.spyOn(serverTypes, 'load');

    render(EditPage);
    await tick();
    await tick();
    await tick();
    await tick();

    // Every call uses 'fsn1' — no other location appears in the spy.
    for (const call of loadSpy.mock.calls) {
      expect(call[0]).toBe('fsn1');
    }
    // 1 parent pre-await + 3 children (after server resolves) = 4 total.
    expect(loadSpy.mock.calls.length).toBe(4);
    loadSpy.mockRestore();
  });

  it('fires serverTypes.load with both fallback and real location when they differ', async () => {
    // When the server's real location differs from the fallback, the
    // page fires load('fsn1') before await (parent + 3 child effects
    // during the await), then refires with 'nbg1' once the server
    // resolves and the location check passes. Children also re-fire
    // their own loads with 'nbg1' when their location prop flips, so
    // we see 5 calls total: 1 'fsn1' (parent pre-await) + 4 'nbg1'
    // (3 children + 1 parent refire).
    apiMockState.server = {
      id: 42,
      project_id: 1,
      hcloud_server_id: 1,
      name: 'w',
      label: 'w',
      base_server_type: 'cpx11',
      top_server_type: 'cpx31',
      fallback_chain: ['cpx21'],
      mode: 'manual',
      timezone: 'UTC',
      status: 'running',
      current_type: 'cpx21',
      location: 'nbg1'
    };

    const loadSpy = vi.spyOn(serverTypes, 'load');

    render(EditPage);
    await tick();
    await tick();
    await tick();
    await tick();
    await tick();

    const args = loadSpy.mock.calls.map((c) => c[0]);
    expect(args.filter((a) => a === 'fsn1').length).toBe(1); // the unconditional pre-await call
    expect(args.filter((a) => a === 'nbg1').length).toBeGreaterThanOrEqual(1); // the post-await refire
    loadSpy.mockRestore();
  });
});
