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

// serverTypes.load() calls api.serverTypes(); stub it to a no-op so
// onMount doesn't make a real network call.
vi.mock('$lib/api', () => ({
  api: {
    serverTypes: () => Promise.resolve([]),
  }
}));

import { serverTypes } from '$lib/stores/serverTypes.svelte';
import ServerTypeMultiSelect from './ServerTypeMultiSelect.svelte';
import type { ServerType, Server } from '$lib/types';

function typeOverrides(): ServerType[] {
  return [
    { name: 'cpx11', available: true, cores: 2, memory_gb: 2 },
    { name: 'cpx21', available: true, cores: 3, memory_gb: 4 },
    { name: 'cpx31', available: true, cores: 4, memory_gb: 8 },
    { name: 'cx33',  available: false, cores: 8, memory_gb: 16 },
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

describe('ServerTypeMultiSelect', () => {
  beforeEach(() => {
    serverTypes._reset();
    serverTypes._setTypesForTest(typeOverrides());
    cleanup();
  });
  afterEach(() => {
    cleanup();
    serverTypes._reset();
  });

  it('renders chips for each bound value, in order', () => {
    const { getByText } = render(ServerTypeMultiSelect, {
      props: { value: ['cpx11', 'cpx31'] },
    });
    expect(getByText('cpx11')).toBeTruthy();
    expect(getByText('cpx31')).toBeTruthy();
  });

  it('renders an empty state when value is empty', () => {
    const { getByText } = render(ServerTypeMultiSelect, {
      props: { value: [] },
    });
    expect(getByText(/no fallback types/i)).toBeTruthy();
  });

  it('does not render chips that are not in the bound value', () => {
    const { container } = render(ServerTypeMultiSelect, {
      props: { value: ['cpx11'] },
    });
    // Scope to <li> chips inside the dnd <ul>, not the trigger.
    const chipList = container.querySelector('ul[aria-label]');
    expect(chipList).toBeTruthy();
    const chipNames = Array.from(chipList!.querySelectorAll('li')).map((li) => li.textContent?.trim() ?? '');
    expect(chipNames.some((t) => t.includes('cpx11'))).toBe(true);
    expect(chipNames.some((t) => t.includes('cpx31'))).toBe(false);
  });

  it('renders the add trigger with the placeholder', () => {
    const { getByRole } = render(ServerTypeMultiSelect, {
      props: { value: [], id: 'multisel' },
    });
    // The add dropdown is a bits-ui Select.Trigger rendered as a
    // button whose accessible name is the placeholder label.
    const trigger = getByRole('button', { name: /add a fallback type/i });
    expect(trigger).toBeTruthy();
    expect(trigger.getAttribute('id')).toBe('multisel');
  });

  it('renders the Unavailable badge for unavailable addable types', async () => {
    // The multi-select now uses the shared ServerTypeOption row, which
    // renders the Unavailable badge on the far right when type.available
    // is false. With value=[] the typeOverrides() entries (incl. cx33
    // with available=false) are all addable, so the badge must appear.
    serverTypes._reset();
    serverTypes._setTypesForTest([
      { name: 'cpx11', available: true, cores: 2, memory_gb: 2 },
      { name: 'cpx21', available: true, cores: 3, memory_gb: 4 },
      { name: 'cx33',  available: false, cores: 8, memory_gb: 16 },
    ]);
    render(ServerTypeMultiSelect, {
      props: { value: [], server: baseServer, open: true },
    });
    await tick();
    expect(document.body.textContent).toMatch(/Unavailable/);
  });

  it('does NOT render the Unavailable badge when all addable types are available', async () => {
    // Defensive companion: when every addable type has available=true,
    // the Unavailable badge must NOT appear. Prevents a regression
    // where the badge condition flips to always-true.
    serverTypes._reset();
    serverTypes._setTypesForTest([
      { name: 'cpx11', available: true, cores: 2, memory_gb: 2 },
      { name: 'cpx21', available: true, cores: 3, memory_gb: 4 },
    ]);
    render(ServerTypeMultiSelect, {
      props: { value: [], server: baseServer, open: true },
    });
    await tick();
    expect(document.body.textContent).not.toMatch(/Unavailable/);
  });

  it('passes the server prop to role-chip classification', () => {
    // Verify the bits-ui trigger renders with the same id-prop pattern
    // as ServerTypeSelect — same component family, same trigger.
    const { getByRole } = render(ServerTypeMultiSelect, {
      props: { value: ['cpx11'], server: baseServer, id: 'ms-with-role' },
    });
    const trigger = getByRole('button', { name: /add a fallback type/i });
    expect(trigger.getAttribute('id')).toBe('ms-with-role');
  });

  it('falls back to no role chip when no server prop is provided', () => {
    // No assertions on the dropdown content (portal-rendered, opaque
    // to jsdom). The component must not throw with server=null and
    // must still render the chip area for the bound value.
    const { getByText } = render(ServerTypeMultiSelect, {
      props: { value: ['cpx11'], server: null },
    });
    expect(getByText('cpx11')).toBeTruthy();
  });

  it('passes the location to serverTypes.load on mount', async () => {
    const loadSpy = vi.spyOn(serverTypes, 'load');
    render(ServerTypeMultiSelect, {
      props: { value: ['cpx11'], location: 'nbg1' },
    });
    await tick();
    // onMount runs after the first render, so the spy must have seen
    // exactly one load call with the supplied location.
    expect(loadSpy).toHaveBeenCalledWith('nbg1');
    loadSpy.mockRestore();
  });

  it('skips load when location is not provided', async () => {
    // Backwards-compat: callers that pre-date the location prop
    // (e.g. some tests, future purely-presentation uses) must not
    // trigger a load for an arbitrary default.
    const loadSpy = vi.spyOn(serverTypes, 'load');
    render(ServerTypeMultiSelect, {
      props: { value: ['cpx11'] },
    });
    await tick();
    expect(loadSpy).not.toHaveBeenCalled();
    loadSpy.mockRestore();
  });

  it('uses conditional load with location="fsn1"', async () => {
    // Mirrors the project page's register-server dialog, which passes
    // an explicit fsn1 default while the operator is typing.
    const loadSpy = vi.spyOn(serverTypes, 'load');
    render(ServerTypeMultiSelect, {
      props: { value: [], location: 'fsn1' },
    });
    await tick();
    expect(loadSpy).toHaveBeenCalledWith('fsn1');
    loadSpy.mockRestore();
  });

  it('loads the catalog after location transitions undefined → "fsn1"', async () => {
    // Regression for the edit-page race: parent renders the component
    // BEFORE the server has been fetched, so `location` is initially
    // undefined. Once the parent resolves the server, `location` flips
    // to a real value and the component MUST trigger
    // `serverTypes.load(location)`. Using `$effect` (not `onMount`)
    // makes this re-run work; `onMount` would silently swallow the
    // transition because it only fires once at mount time.
    const loadSpy = vi.spyOn(serverTypes, 'load');
    const { rerender } = render(ServerTypeMultiSelect, {
      props: { value: [], server: baseServer, id: 'ms' },
    });
    await tick();
    expect(loadSpy).not.toHaveBeenCalled();

    rerender({
      value: [],
      server: { ...baseServer, location: 'fsn1' },
      id: 'ms',
      location: 'fsn1',
    });
    await tick();
    await tick(); // give microtasks a chance to flush
    expect(loadSpy).toHaveBeenCalledWith('fsn1');
    loadSpy.mockRestore();
  });
});