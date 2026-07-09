import { describe, it, expect, afterEach } from 'vitest';
import { render, cleanup } from '@testing-library/svelte';
import ModePill from './ModePill.svelte';

describe('ModePill', () => {
  afterEach(() => {
    cleanup();
  });

  it('renders Manual with no extra line for mode=manual', () => {
    const { getByText, queryByText } = render(ModePill, {
      props: { mode: 'manual', promoteState: null, lastTickAt: null, windows: [], timezone: 'UTC' },
    });
    expect(getByText(/manual/i)).toBeTruthy();
    expect(queryByText(/request promote/i)).toBeNull();
  });

  it('renders "Ready" for auto_promote with null promote_state', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'auto_promote', promoteState: null, lastTickAt: null, windows: [], timezone: 'UTC' },
    });
    expect(getByText(/ready/i)).toBeTruthy();
    expect(getByText(/request promote to start/i)).toBeTruthy();
  });

  it('renders pending state for auto_promote with promote_requested', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'auto_promote', promoteState: 'promote_requested', lastTickAt: null, windows: [], timezone: 'UTC' },
    });
    expect(getByText(/promote/i)).toBeTruthy();
    expect(getByText(/requested · waiting/i)).toBeTruthy();
  });

  it('renders warning for scheduled with no windows', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'scheduled', promoteState: null, lastTickAt: null, windows: [], timezone: 'UTC' },
    });
    expect(getByText(/no windows configured/i)).toBeTruthy();
  });

  it('renders next-window label for scheduled with windows', () => {
    const w = [{
      days_of_week: 0b1111111, start_time: '09:00', stop_time: '18:00',
      target_type: 'cpx31', enabled: true,
    }];
    // noon UTC → in_window; component reflects that.
    const { getByText } = render(ModePill, {
      props: { mode: 'scheduled', promoteState: null, lastTickAt: null, windows: w, timezone: 'UTC', now: new Date('2026-07-08T12:00:00Z') },
    });
    expect(getByText(/cpx31/)).toBeTruthy();
  });

  it('evaluates windows in the server timezone, not UTC', () => {
    // Sunday-only window at 09:00 local. At 2026-07-04T12:00:00Z
    // (Saturday noon UTC = Saturday 08:00 NY), the next window opens at
    // 09:00 NY on Sun Jul 5 = 13:00 UTC, NOT 09:00 UTC. The bug would
    // render the UTC-anchored next-start; the fix must render the
    // America/New_York-anchored next-start. Both are Sunday, but the
    // local-clock time differs by the NY↔UTC offset regardless of the
    // test runner's local timezone.
    const w = [{
      days_of_week: 0b0000001, start_time: '09:00', stop_time: '18:00',
      target_type: 'cpx31', enabled: true,
    }];
    const { getByText, queryByText } = render(ModePill, {
      props: {
        mode: 'scheduled',
        promoteState: null,
        lastTickAt: null,
        windows: w,
        timezone: 'America/New_York',
        now: new Date('2026-07-04T12:00:00Z'),
      },
    });
    // Next window starts at 2026-07-05T13:00:00Z (09:00 NY on Sun Jul 5).
    expect(getByText(/Sun 15:00/)).toBeTruthy();
    // The UTC-only bug would have produced 2026-07-05T09:00:00Z, i.e. "Sun 11:00".
    expect(queryByText(/Sun 11:00/)).toBeNull();
  });

  it('renders last tick age', () => {
    const now = new Date('2026-07-08T12:00:00Z');
    const lastTick = new Date(now.getTime() - 14_000).toISOString();
    const { getByText } = render(ModePill, {
      props: { mode: 'manual', promoteState: null, lastTickAt: lastTick, windows: [], timezone: 'UTC', now },
    });
    expect(getByText(/14s ago/i)).toBeTruthy();
  });

  it('renders "never" when no last tick', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'manual', promoteState: null, lastTickAt: null, windows: [], timezone: 'UTC' },
    });
    expect(getByText(/never/i)).toBeTruthy();
  });
});