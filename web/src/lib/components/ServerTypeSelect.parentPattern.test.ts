import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, cleanup } from '@testing-library/svelte';
import { tick } from 'svelte';

// $env/dynamic/public's virtual module returns undefined for `env` in
// vitest; the paraglide `m` helper pulls from there. Mock to no-op.
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

// serverTypes.load() calls api.serverTypes(); stub it to a no-op so the
// component's effect doesn't make a real network call.
vi.mock('$lib/api', () => ({
  api: {
    serverTypes: () => Promise.resolve([]),
  }
}));

import { serverTypes } from '$lib/stores/serverTypes.svelte';
import ParentPattern from './ServerTypeSelect.ParentPattern.test.svelte';
import type { Server } from '$lib/types';

const baseServer: Server = {
  id: 1, project_id: 1, hcloud_server_id: 1,
  name: 'w', label: 'w',
  base_server_type: 'cpx11', top_server_type: 'cpx31',
  fallback_chain: ['cpx21'],
  mode: 'manual', timezone: 'UTC',
  status: 'running', current_type: 'cpx21',
};

describe('ServerTypeSelect — parent-derived location pattern', () => {
  beforeEach(() => {
    serverTypes._reset();
    cleanup();
  });
  afterEach(() => {
    cleanup();
    serverTypes._reset();
  });

  it('fires serverTypes.load when the parent flips server from null → defined', async () => {
    // This is the EXACT pattern used by the edit page: the parent
    // keeps `let server = $state<Server | null>(null);`, renders the
    // child with `location={server?.location}`, and resolves `server`
    // inside onMount. The component must fire `serverTypes.load`
    // once the derived location flips from undefined to "fsn1".
    const loadSpy = vi.spyOn(serverTypes, 'load');

    const { rerender } = render(ParentPattern, {
      props: { server: null },
    });
    await tick();
    expect(loadSpy).not.toHaveBeenCalled();

    // Simulate the parent's onMount resolving the server.
    rerender({ server: { ...baseServer, location: 'fsn1' } });
    await tick();
    await tick();
    await tick();

    expect(loadSpy).toHaveBeenCalledWith('fsn1');
    loadSpy.mockRestore();
  });
});
