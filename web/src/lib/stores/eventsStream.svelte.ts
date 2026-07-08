import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';
import type { RescaleEvent } from '$lib/types';
import { pendingRescale } from './pendingRescale.svelte';

const MAX_EVENTS = 100;
const INITIAL_BACKOFF_MS = 1000;
const MAX_BACKOFF_MS = 30000;

class EventsStreamStore {
  events = $state<RescaleEvent[]>([]);
  private es: EventSource | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private backoff = INITIAL_BACKOFF_MS;

  connect() {
    if (!browser || this.es) return;
    // Browser EventSource cannot set custom headers, so the SSE handler
    // accepts the shared secret via ?token=… in addition to the
    // X-Internal-Token header. Fall back to import.meta.env so vi.stubEnv
    // still works in tests (mirrors the pattern used in lib/api.ts).
    const token =
      (env.PUBLIC_INTERNAL_TOKEN as string | undefined) ??
      (import.meta.env.PUBLIC_INTERNAL_TOKEN as string | undefined) ??
      '';
    const url = token
      ? `/api/events/stream?token=${encodeURIComponent(token)}`
      : '/api/events/stream';
    this.es = new EventSource(url);
    this.es.addEventListener('message', (e: MessageEvent) => {
      try {
        const data = JSON.parse(e.data) as RescaleEvent | { ok?: boolean };
        if ('id' in data && typeof data.id === 'number') {
          this.events = [data as RescaleEvent, ...this.events].slice(0, MAX_EVENTS);
          // Terminal rescale events clear the in-flight pending state
          // for the affected server — the page's badge disappears
          // without a manual refresh.
          const ev = data as RescaleEvent;
          if (ev.kind === 'rescale_completed' || ev.kind === 'rescale_failed') {
            pendingRescale.clear(ev.server_id);
          }
        }
      } catch {
        /* ignore malformed messages */
      }
    });
    // The `rescale_pending` event name rides phase updates: the
    // Manager broadcasts a fresh rescale_pending row at each phase
    // boundary. We forward to the pendingRescale store, which is
    // keyed by server_id (id is the row id but the map key is the
    // server). The badge component reads from the store via $derived.
    this.es.addEventListener('rescale_pending', (e: MessageEvent) => {
      try {
        const data = JSON.parse(e.data) as RescaleEvent;
        if (typeof data.id === 'number' && data.kind === 'rescale_pending') {
          pendingRescale.upsert(data);
        }
      } catch {
        /* ignore malformed messages */
      }
    });
    this.es.addEventListener('error', () => {
      this.scheduleReconnect();
    });
    this.es.addEventListener('open', () => {
      this.backoff = INITIAL_BACKOFF_MS;
    });
  }

  disconnect() {
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.reconnectTimer = null;
    if (this.es) {
      this.es.close();
      this.es = null;
    }
  }

  replaceAll(items: RescaleEvent[]) {
    this.events = items.slice(0, MAX_EVENTS);
  }

  private scheduleReconnect() {
    if (!browser) return;
    this.disconnect();
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect();
    }, this.backoff);
    this.backoff = Math.min(this.backoff * 2, MAX_BACKOFF_MS);
  }
}

export const eventsStream = new EventsStreamStore();