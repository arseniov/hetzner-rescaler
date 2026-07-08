import { describe, it, expect } from 'vitest';
import type { ServerType, Server } from '$lib/types';
import { roleFor } from './serverTypeRoles';

const t = (name: string): ServerType => ({
  name, available: true, cores: 2, memory_gb: 4, price_monthly_eur: 5,
});

const baseServer: Server = {
  id: 1, project_id: 1, hcloud_server_id: 1,
  name: 'w', label: 'w',
  base_server_type: 'cpx11', top_server_type: 'cpx31',
  fallback_chain: ['cpx21'],
  mode: 'manual', timezone: 'UTC',
  status: 'running', current_type: 'cpx21',
};

describe('roleFor', () => {
  it('returns "current" when the type matches current_type', () => {
    expect(roleFor(t('cpx21'), baseServer)).toBe('current');
  });

  it('returns "base" when the type matches base_server_type', () => {
    expect(roleFor(t('cpx11'), baseServer)).toBe('base');
  });

  it('returns "top" when the type matches top_server_type', () => {
    expect(roleFor(t('cpx31'), baseServer)).toBe('top');
  });

  it('returns "fallback" when the type is in fallback_chain', () => {
    const svr: Server = { ...baseServer, current_type: 'cx33' };
    expect(roleFor(t('cpx21'), svr)).toBe('fallback');
  });

  it('returns null when the type plays no role', () => {
    const svr: Server = { ...baseServer, top_server_type: 'cx33' };
    expect(roleFor(t('cpx31'), svr)).toBeNull();
  });

  it('prioritises current over base when a type appears in both', () => {
    const svr: Server = { ...baseServer, current_type: 'cpx11' };
    expect(roleFor(t('cpx11'), svr)).toBe('current');
  });

  it('handles missing current_type (e.g. server not yet provisioned)', () => {
    const svr: Server = { ...baseServer, current_type: '' };
    // current_type is empty so we fall through to base.
    expect(roleFor(t('cpx11'), svr)).toBe('base');
  });

  it('handles empty fallback_chain', () => {
    const svr: Server = { ...baseServer, fallback_chain: [] };
    expect(roleFor(t('cpx21'), svr)).not.toBe('fallback');
  });
});