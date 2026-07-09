import { describe, it, expect } from 'vitest';
import { nextWindow } from './windowSchedule';

const win = (
  daysOfWeek: number,
  start: string,
  stop: string,
  target: string,
  enabled = true
) => ({
  days_of_week: daysOfWeek,
  start_time: start,
  stop_time: stop,
  target_type: target,
  enabled,
});

describe('nextWindow', () => {
  it('returns none when no windows are configured', () => {
    expect(nextWindow([], 'UTC', new Date('2026-07-08T12:00:00Z'))).toEqual({ kind: 'none' });
  });

  it('returns none when all windows are disabled', () => {
    expect(
      nextWindow(
        [win(0b1111111, '00:00', '23:59', 'cpx31', false)],
        'UTC',
        new Date('2026-07-08T12:00:00Z')
      )
    ).toEqual({ kind: 'none' });
  });

  it('detects an in-window state', () => {
    const out = nextWindow(
      [win(0b1111111, '09:00', '18:00', 'cpx31')],
      'UTC',
      new Date('2026-07-08T12:00:00Z')
    );
    expect(out.kind).toBe('in_window');
    if (out.kind !== 'in_window') return;
    expect(out.target).toBe('cpx31');
    expect(out.endsAt.toISOString()).toBe('2026-07-08T18:00:00.000Z');
  });

  it('finds the next window strictly after now', () => {
    // Monday July 6 2026, 23:00 UTC — last today's window ended at 22:00,
    // next occurrence is Tue Jul 7 09:00.
    const out = nextWindow(
      [win(0b1111111, '09:00', '22:00', 'cpx31')],
      'UTC',
      new Date('2026-07-06T23:00:00Z')
    );
    expect(out.kind).toBe('next');
    if (out.kind !== 'next') return;
    expect(out.startsAt.toISOString()).toBe('2026-07-07T09:00:00.000Z');
    expect(out.target).toBe('cpx31');
  });

  it('returns the earliest upcoming occurrence across multiple windows', () => {
    // Wed Jul 8 12:00 UTC. Window A starts same day at 14:00, window B
    // starts next Monday at 09:00. A wins.
    const out = nextWindow(
      [
        win(0b0011111, '14:00', '18:00', 'cpx31'), // Mon-Fri
        win(0b0000001, '09:00', '12:00', 'cpx21'), // Sunday only — bit 0 = Sun
      ],
      'UTC',
      new Date('2026-07-08T12:00:00Z')
    );
    expect(out.kind).toBe('next');
    if (out.kind !== 'next') return;
    expect(out.startsAt.toISOString()).toBe('2026-07-08T14:00:00.000Z');
    expect(out.target).toBe('cpx31');
  });

  it('picks the soonest end when currently inside overlapping windows', () => {
    const out = nextWindow(
      [
        win(0b1111111, '00:00', '23:00', 'cpx31'),
        win(0b1111111, '09:00', '17:00', 'cpx21'),
      ],
      'UTC',
      new Date('2026-07-08T12:00:00Z')
    );
    expect(out.kind).toBe('in_window');
    if (out.kind !== 'in_window') return;
    expect(out.target).toBe('cpx21'); // narrower window wins for in_window report
    expect(out.endsAt.toISOString()).toBe('2026-07-08T17:00:00.000Z');
  });

  it('detects in-window for wrap-around windows in the stop-day tail', () => {
    // Window runs 22:00 -> 06:00 next morning, every day.
    // At 03:00 UTC on a Tuesday morning, the window that started Monday
    // 22:00 is still active. The day bit check must consider yesterday.
    const out = nextWindow(
      [win(0b1111111, '22:00', '06:00', 'cpx31')],
      'UTC',
      new Date('2026-07-07T03:00:00Z') // Tuesday 03:00 UTC
    );
    expect(out.kind).toBe('in_window');
    if (out.kind !== 'in_window') return;
    expect(out.target).toBe('cpx31');
    // The window ends at 06:00 UTC on the same day (Tuesday 06:00).
    expect(out.endsAt.toISOString()).toBe('2026-07-07T06:00:00.000Z');
    // And it started at 22:00 UTC on the previous day (Monday 22:00).
    expect(out.startedAt.toISOString()).toBe('2026-07-06T22:00:00.000Z');
  });
});