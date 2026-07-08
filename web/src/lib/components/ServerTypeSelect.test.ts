import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, cleanup } from '@testing-library/svelte';

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
// component's onMount doesn't make a real network call.
vi.mock('$lib/api', () => ({
  api: {
    serverTypes: () => Promise.resolve([]),
  }
}));

import { serverTypes } from '$lib/stores/serverTypes.svelte';
import ServerTypeSelect from './ServerTypeSelect.svelte';
import type { ServerType, Server } from '$lib/types';
import { roleFor } from '$lib/utils/serverTypeRoles';

function typeOverrides(): ServerType[] {
  return [
    { name: 'cpx11', available: true, cores: 2, memory_gb: 2, price_monthly_eur: 4.85 },
    { name: 'cpx21', available: true, cores: 3, memory_gb: 4, price_monthly_eur: 9.85 },
    { name: 'cpx31', available: true, cores: 4, memory_gb: 8, price_monthly_eur: 19.85 },
    { name: 'cx33', available: false, cores: 8, memory_gb: 16, price_monthly_eur: 39.0 },
  ];
}

const baseServer: Server = {
  id: 1, project_id: 1, hcloud_server_id: 1,
  name: 'w', label: 'w',
  base_server_type: 'cpx11', top_server_type: 'cpx31',
  fallback_chain: ['cpx21'],
  mode: 'manual', timezone: 'UTC',
  status: 'running', current_type: 'cpx21',
};

describe('roleFor (pure)', () => {
  const types = typeOverrides();

  it('returns "current" when the type matches current_type', () => {
    expect(roleFor(types.find((t) => t.name === 'cpx21')!, baseServer)).toBe('current');
  });

  it('returns "base" when the type matches base_server_type', () => {
    expect(roleFor(types.find((t) => t.name === 'cpx11')!, baseServer)).toBe('base');
  });

  it('returns "top" when the type matches top_server_type', () => {
    expect(roleFor(types.find((t) => t.name === 'cpx31')!, baseServer)).toBe('top');
  });

  it('returns "fallback" when the type is in fallback_chain', () => {
    // cpx21 is current *and* fallback for this server; current wins.
    // Use a fresh server where cpx21 is only fallback.
    const svr: Server = { ...baseServer, current_type: 'cx33' };
    expect(roleFor(types.find((t) => t.name === 'cpx21')!, svr)).toBe('fallback');
  });

  it('returns null when the type has no role on this server', () => {
    // cpx31 is "top" by default; switch to a server where it's neither.
    const svr: Server = { ...baseServer, top_server_type: 'cx33' };
    expect(roleFor(types.find((t) => t.name === 'cpx31')!, svr)).toBeNull();
  });

  it('prioritises current > base > top > fallback when a type appears in multiple buckets', () => {
    const svr: Server = { ...baseServer, current_type: 'cpx11' }; // cpx11 is both current AND base
    expect(roleFor(types.find((t) => t.name === 'cpx11')!, svr)).toBe('current');
  });
});

describe('ServerTypeSelect', () => {
  beforeEach(() => {
    serverTypes._reset();
    serverTypes._setTypesForTest(typeOverrides());
    cleanup();
  });
  afterEach(() => {
    cleanup();
    serverTypes._reset();
  });

  it('renders the trigger with the selected type name', () => {
    const { getByRole } = render(ServerTypeSelect, {
      props: { value: 'cpx31', server: baseServer, id: 'sel' },
    });
    // The bits-ui Select.Trigger renders as a button whose accessible
    // name includes the currently-selected value.
    const trigger = getByRole('button', { name: /cpx31/ });
    expect(trigger).toBeTruthy();
    expect(trigger.getAttribute('id')).toBe('sel');
  });

  it('renders the placeholder when no value is selected and not required', () => {
    const { getByRole } = render(ServerTypeSelect, {
      props: { value: '', server: baseServer, id: 'sel' },
    });
    // Em-dash placeholder is shown in the trigger label.
    expect(getByRole('button', { name: /—/ })).toBeTruthy();
  });
});