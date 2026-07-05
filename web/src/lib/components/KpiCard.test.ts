import { describe, it, expect } from 'vitest';
import { render } from '@testing-library/svelte';
import KpiCard from './KpiCard.svelte';

describe('KpiCard', () => {
  it('renders label and value', () => {
    const { getByText } = render(KpiCard, { label: 'Active servers', value: 12 });
    expect(getByText('Active servers')).toBeTruthy();
    expect(getByText('12')).toBeTruthy();
  });

  it('renders hint when provided', () => {
    const { getByText } = render(KpiCard, {
      label: 'Active servers',
      value: 12,
      hint: 'Servers in auto_promote or scheduled mode'
    });
    expect(getByText('Servers in auto_promote or scheduled mode')).toBeTruthy();
  });

  it('omits hint element when not provided', () => {
    const { queryByText, container } = render(KpiCard, { label: 'Active servers', value: 12 });
    expect(queryByText('Servers in auto_promote or scheduled mode')).toBeNull();
    // The hint is rendered as a <p> with text-xs class; no element with text-xs should exist
    expect(container.querySelectorAll('p.text-xs').length).toBe(0);
  });

  it('shows ellipsis placeholder while loading', () => {
    const { container } = render(KpiCard, {
      label: 'Active servers',
      value: 12,
      loading: true
    });
    // The numeric value should NOT be rendered while loading; instead an
    // animated "…" placeholder is shown.
    expect(container.textContent).not.toContain('12');
    expect(container.querySelector('.animate-pulse')).toBeTruthy();
  });

  it('shows em-dash when value is null after fetch completes', () => {
    const { container } = render(KpiCard, {
      label: 'Last error',
      value: null,
      loading: false
    });
    // When fetch completed but the API returned null, render an em-dash
    // rather than "Loading…".
    expect(container.textContent).toContain('—');
    expect(container.querySelector('.animate-pulse')).toBeNull();
  });

  it('shows em-dash when value is undefined after fetch completes', () => {
    const { container } = render(KpiCard, {
      label: 'Last error',
      value: undefined,
      loading: false
    });
    expect(container.textContent).toContain('—');
  });
});