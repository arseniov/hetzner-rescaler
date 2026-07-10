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

// Mock the api module BEFORE importing the page so the page's
// `import { api }` picks up the stubbed getServer.
vi.mock('$lib/api', () => {
  return {
    api: {
      getServer: vi.fn().mockResolvedValue({
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
      }),
      updateServer: vi.fn(),
      // The page's explicit parent trigger calls serverTypes.load(),
      // which calls api.serverTypes(). Return an empty catalog so the
      // load resolves without making a real network call.
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
// Importing the actual page — verifies the explicit parent trigger
// the user requested ("it should fire each time we go into a server
// edit page"). The page calls `serverTypes.load(server.location)`
// inside its onMount, in addition to the $effect in the child
// component, so the network call is guaranteed to fire.
import EditPage from './+page.svelte';

describe('edit page — explicit parent trigger', () => {
  beforeEach(() => {
    serverTypes._reset();
    cleanup();
  });
  afterEach(() => {
    cleanup();
    serverTypes._reset();
  });

  it('fires serverTypes.load with the server location on mount', async () => {
    const loadSpy = vi.spyOn(serverTypes, 'load');

    // The mocked api.getServer returns a server with location: 'fsn1'.
    expect(api.getServer).toBeDefined();

    render(EditPage);
    // onMount fires after the first render. Wait a few ticks to let
    // the async api.getServer resolve and the parent-side trigger run.
    await tick();
    await tick();
    await tick();
    await tick();

    expect(api.getServer).toHaveBeenCalledWith(42);
    expect(loadSpy).toHaveBeenCalledWith('fsn1');
    loadSpy.mockRestore();
  });
});
