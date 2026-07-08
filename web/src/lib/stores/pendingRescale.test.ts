import { describe, it, expect, beforeEach } from 'vitest';
import { pendingRescale } from './pendingRescale.svelte';
import type { RescaleEvent } from '$lib/types';

function makeEvent(overrides: Partial<RescaleEvent> = {}): RescaleEvent {
  return {
    id: 1, server_id: 42, kind: 'rescale_pending',
    started_at: '2026-07-07T00:00:00Z',
    finished_at: '0001-01-01T00:00:00Z',
    ok: false, triggered_by: 'api',
    ...overrides,
  };
}

describe('pendingRescale store', () => {
  beforeEach(() => {
    pendingRescale.clear();
  });

  it('starts empty', () => {
    expect(pendingRescale.get(42)).toBeUndefined();
  });

  it('upsert adds a new pending entry', () => {
    pendingRescale.upsert(makeEvent({ id: 1, server_id: 42, phase: 'shutting_down' }));
    const got = pendingRescale.get(42);
    expect(got?.id).toBe(1);
    expect(got?.phase).toBe('shutting_down');
  });

  it('upsert replaces an existing entry by id', () => {
    pendingRescale.upsert(makeEvent({ id: 1, server_id: 42, phase: 'shutting_down' }));
    pendingRescale.upsert(makeEvent({ id: 1, server_id: 42, phase: 'changing_type' }));
    expect(pendingRescale.get(42)?.phase).toBe('changing_type');
  });

  it('upsert does not collide across servers', () => {
    pendingRescale.upsert(makeEvent({ id: 1, server_id: 42, phase: 'shutting_down' }));
    pendingRescale.upsert(makeEvent({ id: 2, server_id: 43, phase: 'powering_on' }));
    expect(pendingRescale.get(42)?.phase).toBe('shutting_down');
    expect(pendingRescale.get(43)?.phase).toBe('powering_on');
  });

  it('clear removes a single server entry', () => {
    pendingRescale.upsert(makeEvent({ id: 1, server_id: 42, phase: 'shutting_down' }));
    pendingRescale.upsert(makeEvent({ id: 2, server_id: 43, phase: 'shutting_down' }));
    pendingRescale.clear(42);
    expect(pendingRescale.get(42)).toBeUndefined();
    expect(pendingRescale.get(43)?.id).toBe(2);
  });

  it('clear() with no args removes all', () => {
    pendingRescale.upsert(makeEvent({ id: 1, server_id: 42 }));
    pendingRescale.upsert(makeEvent({ id: 2, server_id: 43 }));
    pendingRescale.clear();
    expect(pendingRescale.get(42)).toBeUndefined();
    expect(pendingRescale.get(43)).toBeUndefined();
  });
});
