import { describe, it, expect, beforeEach, vi } from 'vitest';
import { serverTypes } from './serverTypes.svelte';

vi.mock('$lib/api', () => ({
  api: {
    serverTypes: vi.fn()
  }
}));

import { api } from '$lib/api';

const mockedApi = api as unknown as { serverTypes: ReturnType<typeof vi.fn> };

beforeEach(() => {
  serverTypes._reset();
  mockedApi.serverTypes.mockReset();
});

describe('ServerTypesStore.load(location)', () => {
  it('passes the location to api.serverTypes', async () => {
    mockedApi.serverTypes.mockResolvedValue([{ name: 'cpx11', available: true }]);
    await serverTypes.load('fsn1');
    expect(mockedApi.serverTypes).toHaveBeenCalledWith('fsn1');
  });

  it('caches per location within TTL', async () => {
    mockedApi.serverTypes.mockResolvedValue([{ name: 'cpx11', available: true }]);
    await serverTypes.load('fsn1');
    await serverTypes.load('fsn1');
    expect(mockedApi.serverTypes).toHaveBeenCalledTimes(1);
  });

  it('refetches when the location changes', async () => {
    mockedApi.serverTypes.mockResolvedValue([{ name: 'cpx11', available: true }]);
    await serverTypes.load('fsn1');
    await serverTypes.load('nbg1');
    expect(mockedApi.serverTypes).toHaveBeenCalledTimes(2);
    expect(mockedApi.serverTypes).toHaveBeenNthCalledWith(1, 'fsn1');
    expect(mockedApi.serverTypes).toHaveBeenNthCalledWith(2, 'nbg1');
  });

  it('updates types to the new location when refetching', async () => {
    mockedApi.serverTypes.mockImplementation(async (loc: string) =>
      loc === 'fsn1'
        ? [{ name: 'cpx11', available: true }]
        : [{ name: 'cpx11', available: false }]
    );
    await serverTypes.load('fsn1');
    expect(serverTypes.types[0].available).toBe(true);
    await serverTypes.load('nbg1');
    expect(serverTypes.types[0].available).toBe(false);
  });

  it('force=true bypasses the per-location cache', async () => {
    mockedApi.serverTypes.mockResolvedValue([{ name: 'cpx11', available: true }]);
    await serverTypes.load('fsn1');
    await serverTypes.load('fsn1', true);
    expect(mockedApi.serverTypes).toHaveBeenCalledTimes(2);
  });
});
