import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, cleanup, within } from '@testing-library/svelte';
import { tick } from 'svelte';

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

vi.mock('$app/stores', () => ({
  page: {
    subscribe: (run: (value: { params: { id: string } }) => void) => {
      run({ params: { id: '42' } });
      return () => {};
    }
  }
}));

const apiMockState: { server: any; windows: any[] } = {
  server: null,
  windows: []
};

vi.mock('$lib/api', () => ({
  api: {
    getServer: vi.fn().mockImplementation(() => Promise.resolve(apiMockState.server)),
    listWindows: vi.fn().mockImplementation(() => Promise.resolve(apiMockState.windows)),
    createWindow: vi.fn(),
    updateWindow: vi.fn(),
    deleteWindow: vi.fn(),
    serverTypes: vi.fn().mockResolvedValue([])
  },
  ApiError: class ApiError extends Error {
    constructor(public status: number, message: string) {
      super(message);
      this.name = 'ApiError';
    }
  }
}));

import { serverTypes } from '$lib/stores/serverTypes.svelte';
import WindowsPage from './+page.svelte';

describe('windows page — weekly agenda preview', () => {
  beforeEach(() => {
    serverTypes._reset();
    cleanup();
    apiMockState.server = {
      id: 42,
      project_id: 1,
      hcloud_server_id: 100,
      name: 'worker-1',
      label: 'worker-1',
      base_server_type: 'cpx11',
      top_server_type: 'cpx31',
      fallback_chain: ['cpx21'],
      mode: 'scheduled',
      timezone: 'Europe/Rome',
      status: 'running',
      current_type: 'cpx11',
      location: 'fsn1'
    };
    apiMockState.windows = [
      {
        id: 1,
        server_id: 42,
        label: 'weekday-peak',
        days_of_week: 0b00101010,
        start_time: '09:00',
        stop_time: '18:00',
        target_type: 'cpx31',
        enabled: true
      },
      {
        id: 2,
        server_id: 42,
        label: 'disabled-test',
        days_of_week: 0b00000001,
        start_time: '01:00',
        stop_time: '02:00',
        target_type: 'cpx21',
        enabled: false
      }
    ];
  });

  afterEach(() => {
    cleanup();
    serverTypes._reset();
  });

  it('renders enabled windows in the days they affect', async () => {
    const { getByLabelText } = render(WindowsPage);
    await tick();
    await tick();
    await tick();

    const agenda = getByLabelText('Weekly agenda');
    const monday = within(agenda).getByLabelText('Mon agenda');
    expect(within(monday).getByText('weekday-peak')).toBeTruthy();
    expect(within(monday).getByText('09:00–18:00')).toBeTruthy();
    expect(within(monday).getByText('cpx31')).toBeTruthy();

    const wednesday = within(agenda).getByLabelText('Wed agenda');
    expect(within(wednesday).getByText('weekday-peak')).toBeTruthy();

    const sunday = within(agenda).getByLabelText('Sun agenda');
    expect(within(sunday).queryByText('disabled-test')).toBeNull();
    expect(within(sunday).getByText('No active windows')).toBeTruthy();
  });
});
