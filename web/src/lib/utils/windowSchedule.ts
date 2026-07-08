/**
 * Pure helpers for evaluating scheduled rescale windows against the
 * server's timezone on the frontend.
 *
 * The source of truth for "are we in a window right now?" lives in Go
 * (`internal/scheduler.EvaluateWindows`). This module is a *display-only*
 * copy used to populate the server detail page's "next window" label.
 * Differences between this and the Go implementation surface as the user
 * seeing a stale label for at most one tick; the next dispatched rescale
 * still uses the authoritative Go path.
 */

export interface WindowSpec {
  days_of_week: number; // bitmask: bit 0 = Sunday ... bit 6 = Saturday
  start_time: string;   // 'HH:MM' (24h)
  stop_time: string;    // 'HH:MM' (24h)
  target_type: string;
  enabled: boolean;
}

export type WindowState =
  | { kind: 'none' }
  | { kind: 'in_window'; target: string; startedAt: Date; endsAt: Date }
  | { kind: 'next'; target: string; startsAt: Date; endsAt: Date };

const MINUTES_PER_DAY = 24 * 60;

/** Parse 'HH:MM' into { hour, minute }. */
function parseHHMM(hhmm: string): { hour: number; minute: number } | null {
  const m = /^(\d{2}):(\d{2})$/.exec(hhmm);
  if (!m) return null;
  const hour = parseInt(m[1], 10);
  const minute = parseInt(m[2], 10);
  if (!Number.isFinite(hour) || !Number.isFinite(minute)) return null;
  if (hour < 0 || hour > 23 || minute < 0 || minute > 59) return null;
  return { hour, minute };
}

/** Bit 0 = Sunday ... bit 6 = Saturday. */
function dayBit(date: Date, timezone: string): number {
  const weekday = new Intl.DateTimeFormat('en-US', {
    timeZone: timezone,
    weekday: 'short',
  }).format(date);
  const order = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
  return order.indexOf(weekday); // -1 if unknown
}

/** Minutes-since-00:00 in the timezone of `date`. */
function minutesInTz(date: Date, timezone: string): number {
  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone: timezone,
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  }).format(date);
  // '24' can appear at midnight in some locales; normalise.
  const [hh, mm] = parts.split(':').map((s) => parseInt(s, 10));
  return (hh === 24 ? 0 : hh) * 60 + mm;
}

/** Get the year/month/day components of `date` in `timezone`. */
function datePartsInTz(date: Date, timezone: string): { year: number; month: number; day: number } {
  const parts = new Intl.DateTimeFormat('en-CA', {
    timeZone: timezone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).formatToParts(date);
  const get = (type: string) => parts.find((p) => p.type === type)?.value ?? '';
  return {
    year: parseInt(get('year'), 10),
    month: parseInt(get('month'), 10),
    day: parseInt(get('day'), 10),
  };
}

/** Convert a wall-clock (year/month/day/hour/minute) in `timezone` to a UTC Date.
 * Iterates candidate offsets (a timezone is at most ±14h from UTC) and picks the
 * one whose tz-formatted wall time matches the desired components. This is
 * DST-correct: the formatted value of the candidate UTC instant must match
 * the input wall components exactly. */
function wallToUtc(
  year: number, month: number, day: number, hour: number, minute: number,
  timezone: string,
): Date {
  // Build a UTC instant that *would* be the wall time if the timezone were UTC.
  const naiveUtc = Date.UTC(year, month - 1, day, hour, minute);
  const fmt = new Intl.DateTimeFormat('en-CA', {
    timeZone: timezone,
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  });
  // Try offsets in 15-minute increments from -14h to +14h.
  // Walk in order of ascending absolute offset so the smallest-offset match wins.
  const offsets: number[] = [];
  for (let mag = 0; mag <= 14 * 4; mag++) {
    if (mag === 0) offsets.push(0);
    else {
      offsets.push(mag);
      offsets.push(-mag);
    }
  }
  for (const offsetQuarters of offsets) {
    const offsetH = Math.trunc(offsetQuarters / 4);
    const offsetM = (offsetQuarters % 4) * 15;
    const candidate = new Date(naiveUtc - (offsetH * 60 + offsetM) * 60_000);
    const parts = fmt.formatToParts(candidate);
    const get = (type: string) => parts.find((p) => p.type === type)?.value ?? '';
    if (
      get('year') === String(year) &&
      get('month') === String(month).padStart(2, '0') &&
      get('day') === String(day).padStart(2, '0') &&
      (get('hour') === String(hour).padStart(2, '0') || (hour === 0 && get('hour') === '24')) &&
      get('minute') === String(minute).padStart(2, '0')
    ) {
      return candidate;
    }
  }
  // Fallback (shouldn't happen for any valid IANA timezone).
  return new Date(naiveUtc);
}

/** Compute the in-window range for a window on the wall-day that contains `now` in `timezone`.
 * Returns null if the window doesn't cover `now`. */
function inWindowRange(
  now: Date,
  timezone: string,
  w: WindowSpec
): { startedAt: Date; endsAt: Date; target: string } | null {
  const start = parseHHMM(w.start_time);
  const stop = parseHHMM(w.stop_time);
  if (!start || !stop) return null;

  const cur = minutesInTz(now, timezone);
  const startMin = start.hour * 60 + start.minute;
  const stopMin = stop.hour * 60 + stop.minute;

  const todayBit = dayBit(now, timezone);
  if (todayBit < 0 || (w.days_of_week & (1 << todayBit)) === 0) {
    return null;
  }

  const today = datePartsInTz(now, timezone);
  const startUtc = wallToUtc(today.year, today.month, today.day, start.hour, start.minute, timezone);
  const stopUtc = stopMin > startMin
    ? wallToUtc(today.year, today.month, today.day, stop.hour, stop.minute, timezone)
    : wallToUtc(today.year, today.month, today.day + 1, stop.hour, stop.minute, timezone);

  // Same-day: start <= cur < stop
  if (stopMin > startMin && cur >= startMin && cur < stopMin) {
    return { startedAt: startUtc, endsAt: stopUtc, target: w.target_type };
  }
  // Wrap-around window (stop <= start), we're on the start day: cur >= startMin.
  if (stopMin <= startMin && cur >= startMin && cur < MINUTES_PER_DAY) {
    return { startedAt: startUtc, endsAt: stopUtc, target: w.target_type };
  }
  return null;
}

export function nextWindow(
  windows: WindowSpec[],
  timezone: string,
  now: Date
): WindowState {
  const enabled = windows.filter((w) => w.enabled);
  if (enabled.length === 0) return { kind: 'none' };

  // In-window: pick the soonest-ending matching window.
  let best: { startedAt: Date; endsAt: Date; target: string } | null = null;
  for (const w of enabled) {
    const r = inWindowRange(now, timezone, w);
    if (!r) continue;
    if (best === null || r.endsAt.getTime() < best.endsAt.getTime()) best = r;
  }
  if (best !== null) {
    return { kind: 'in_window', target: best.target, startedAt: best.startedAt, endsAt: best.endsAt };
  }

  // Next-upcoming: walk forward in 1-minute steps up to 7 days, find the first
  // window-start strictly after now whose day bit matches. 1-minute resolution
  // is required because window start_time is minute-precision (e.g. '09:00'
  // would be skipped by any stride that doesn't divide 60).
  for (let offset = 1; offset <= 7 * 24 * 60; offset++) {
    const cand = new Date(now.getTime() + offset * 60_000);
    const bit = dayBit(cand, timezone);
    if (bit < 0) continue;
    const cur = minutesInTz(cand, timezone);
    const day = datePartsInTz(cand, timezone);

    for (const w of enabled) {
      if ((w.days_of_week & (1 << bit)) === 0) continue;
      const start = parseHHMM(w.start_time);
      const stop = parseHHMM(w.stop_time);
      if (!start || !stop) continue;
      const startMin = start.hour * 60 + start.minute;
      const stopMin = stop.hour * 60 + stop.minute;
      if (cur !== startMin) continue;

      // The window starts at this instant.
      const startsAt = wallToUtc(day.year, day.month, day.day, start.hour, start.minute, timezone);
      const endsAt = stopMin > startMin
        ? wallToUtc(day.year, day.month, day.day, stop.hour, stop.minute, timezone)
        : wallToUtc(day.year, day.month, day.day + 1, stop.hour, stop.minute, timezone);

      return { kind: 'next', target: w.target_type, startsAt, endsAt };
    }
  }

  return { kind: 'none' };
}