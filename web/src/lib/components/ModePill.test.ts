import { describe, it, expect, afterEach } from 'vitest';
import { render, cleanup } from '@testing-library/svelte';
import ModePill from './ModePill.svelte';

describe('ModePill', () => {
  afterEach(() => {
    cleanup();
  });

  it('renders Manual with no extra line for mode=manual', () => {
    const { getByText, queryByText } = render(ModePill, {
      props: { mode: 'manual', promoteState: null, lastTickAt: null, windows: [] },
    });
    expect(getByText(/manual/i)).toBeTruthy();
    expect(queryByText(/request promote/i)).toBeNull();
  });

  it('renders "Ready" for auto_promote with null promote_state', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'auto_promote', promoteState: null, lastTickAt: null, windows: [] },
    });
    expect(getByText(/ready/i)).toBeTruthy();
    expect(getByText(/request promote to start/i)).toBeTruthy();
  });

  it('renders pending state for auto_promote with promote_requested', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'auto_promote', promoteState: 'promote_requested', lastTickAt: null, windows: [] },
    });
    expect(getByText(/promote/i)).toBeTruthy();
    expect(getByText(/requested · waiting/i)).toBeTruthy();
  });

  it('renders warning for scheduled with no windows', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'scheduled', promoteState: null, lastTickAt: null, windows: [] },
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
      props: { mode: 'scheduled', promoteState: null, lastTickAt: null, windows: w, now: new Date('2026-07-08T12:00:00Z') },
    });
    expect(getByText(/cpx31/)).toBeTruthy();
  });

  it('renders last tick age', () => {
    const now = new Date('2026-07-08T12:00:00Z');
    const lastTick = new Date(now.getTime() - 14_000).toISOString();
    const { getByText } = render(ModePill, {
      props: { mode: 'manual', promoteState: null, lastTickAt: lastTick, windows: [], now },
    });
    expect(getByText(/14s ago/i)).toBeTruthy();
  });

  it('renders "never" when no last tick', () => {
    const { getByText } = render(ModePill, {
      props: { mode: 'manual', promoteState: null, lastTickAt: null, windows: [] },
    });
    expect(getByText(/never/i)).toBeTruthy();
  });
});
