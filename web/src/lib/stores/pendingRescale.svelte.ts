import type { RescaleEvent } from '$lib/types';

/**
 * pendingRescale — keyed map of active rescale_pending events by
 * server_id. Fed by:
 *   - the `rescale_pending` SSE listener (phase updates)
 *   - the `rescale` SSE listener (terminal → clear)
 *   - `setFromServer` when the page loads (seeds from ServerResponse.pending_event)
 *
 * `upsert` is keyed by event id; a phase update for the same pending row
 * replaces the existing entry. `clear` removes the entry for a server
 * (or all servers) when the terminal event arrives.
 *
 * The store is reactive — components read `get(serverID)` inside a
 * `$derived` to track changes without manual subscriptions.
 */
class PendingRescaleStore {
  #byServer = $state<Map<number, RescaleEvent>>(new Map());

  get(serverID: number): RescaleEvent | undefined {
    return this.#byServer.get(serverID);
  }

  upsert(event: RescaleEvent): void {
    if (event.kind !== 'rescale_pending') return;
    this.#byServer.set(event.server_id, event);
  }

  clear(serverID?: number): void {
    if (serverID === undefined) {
      this.#byServer = new Map();
      return;
    }
    this.#byServer.delete(serverID);
  }

  setFromServer(event: RescaleEvent | undefined | null): void {
    if (event == null) return;
    this.upsert(event);
  }
}

export const pendingRescale = new PendingRescaleStore();
