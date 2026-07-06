import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import StatusBadge from './StatusBadge.svelte';

describe('StatusBadge', () => {
  it('shows OK for ok status', () => {
    const { getByText } = render(StatusBadge, { status: 'ok' });
    expect(getByText('OK')).toBeTruthy();
  });

  it('shows DEGRADED for degraded status', () => {
    const { getByText } = render(StatusBadge, { status: 'degraded' });
    expect(getByText('DEGRADED')).toBeTruthy();
  });

  it('shows UNKNOWN for unknown status', () => {
    const { getByText } = render(StatusBadge, { status: 'unknown' });
    expect(getByText('UNKNOWN')).toBeTruthy();
  });

  // Hetzner Cloud live states — populated from server.status when
  // the API call succeeds. The badge renders the raw state word
  // uppercased so the operator can grep logs verbatim.
  for (const status of [
    'running',
    'initializing',
    'starting',
    'stopping',
    'off',
    'deleting'
  ] as const) {
    it(`renders ${status.toUpperCase()} for Hetzner live status`, () => {
      const { getByText } = render(StatusBadge, { status });
      expect(getByText(status.toUpperCase())).toBeTruthy();
    });
  }
});