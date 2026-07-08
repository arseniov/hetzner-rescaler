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

// serverTypes.load() calls api.serverTypes(); stub it to a no-op so
// onMount doesn't make a real network call.
vi.mock('$lib/api', () => ({
  api: {
    serverTypes: () => Promise.resolve([]),
  }
}));

import { serverTypes } from '$lib/stores/serverTypes.svelte';
import ServerTypeMultiSelect from './ServerTypeMultiSelect.svelte';
import type { ServerType } from '$lib/types';

function typeOverrides(): ServerType[] {
  return [
    { name: 'cpx11', available: true, cores: 2, memory_gb: 2 },
    { name: 'cpx21', available: true, cores: 3, memory_gb: 4 },
    { name: 'cpx31', available: true, cores: 4, memory_gb: 8 },
    { name: 'cx33',  available: false, cores: 8, memory_gb: 16 },
  ];
}

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
    // Scope to <li> chips inside the dnd <ul>, not the <select>
    // options below (those also list unselected types for adding).
    const chipList = container.querySelector('ul[aria-label]');
    expect(chipList).toBeTruthy();
    const chipNames = Array.from(chipList!.querySelectorAll('li')).map((li) => li.textContent?.trim() ?? '');
    expect(chipNames.some((t) => t.includes('cpx11'))).toBe(true);
    expect(chipNames.some((t) => t.includes('cpx31'))).toBe(false);
  });

  it('does not list already-selected types in the addable options', async () => {
    const { container } = render(ServerTypeMultiSelect, {
      props: { value: ['cpx11'] },
    });
    // The add dropdown is a <select>; verify cpx11 is not present as
    // an <option> value (it's already selected).
    const options = container.querySelectorAll('select option');
    const optionValues = Array.from(options).map((o) => o.getAttribute('value'));
    expect(optionValues).not.toContain('cpx11');
    // cpx31 is not in the bound value, so it must be addable.
    expect(optionValues).toContain('cpx31');
  });

  it('excludes types listed in the `excluded` prop from addable options', async () => {
    const { container } = render(ServerTypeMultiSelect, {
      props: { value: [], excluded: ['cpx31'] },
    });
    const options = container.querySelectorAll('select option');
    const optionValues = Array.from(options).map((o) => o.getAttribute('value'));
    expect(optionValues).not.toContain('cpx31');
  });
});