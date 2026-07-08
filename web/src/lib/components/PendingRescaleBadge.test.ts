import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { render, cleanup } from '@testing-library/svelte';
import PendingRescaleBadge from './PendingRescaleBadge.svelte';
import type { RescaleEvent } from '$lib/types';

function makeEvent(overrides: Partial<RescaleEvent> = {}): RescaleEvent {
  return {
    id: 1, server_id: 42, kind: 'rescale_pending',
    started_at: new Date().toISOString(),
    finished_at: '0001-01-01T00:00:00Z',
    ok: false, triggered_by: 'api',
    ...overrides,
  };
}

describe('PendingRescaleBadge', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date('2026-07-07T00:00:30Z'));
  });
  afterEach(() => {
    cleanup();
    vi.useRealTimers();
  });

  it('renders the shutting_down phase label', () => {
    const { getByText } = render(PendingRescaleBadge, {
      props: { event: makeEvent({ started_at: '2026-07-07T00:00:00Z', phase: 'shutting_down' }) },
    });
    expect(getByText(/Shutting down…/)).toBeTruthy();
  });

  it('renders the powering_on phase label', () => {
    const { getByText } = render(PendingRescaleBadge, {
      props: { event: makeEvent({ started_at: '2026-07-07T00:00:00Z', phase: 'powering_on' }) },
    });
    expect(getByText(/Powering on…/)).toBeTruthy();
  });

  it('shows the elapsed time in seconds', () => {
    const { getByText } = render(PendingRescaleBadge, {
      props: { event: makeEvent({ started_at: '2026-07-07T00:00:00Z' }) },
    });
    expect(getByText(/30s/)).toBeTruthy();
  });

  it('shows the elapsed time in minutes + seconds past 60s', () => {
    vi.setSystemTime(new Date('2026-07-07T00:02:05Z'));
    const { getByText } = render(PendingRescaleBadge, {
      props: { event: makeEvent({ started_at: '2026-07-07T00:00:00Z' }) },
    });
    expect(getByText(/2m 5s/)).toBeTruthy();
  });
});
