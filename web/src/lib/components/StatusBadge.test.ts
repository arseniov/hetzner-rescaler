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
});