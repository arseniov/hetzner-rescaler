import { api } from '$lib/api';
import type { ServerType } from '$lib/types';

// Soft TTL — 5 minutes. The server-types catalog rarely changes (a new
// type appears maybe quarterly), so re-fetching on every page mount is
// wasteful. Five minutes is the same window used by Hetzner Cloud UI
// for its server-type dropdowns.
const TTL_MS = 5 * 60 * 1000;

class ServerTypesStore {
  types = $state<ServerType[]>([]);
  loadedAt = $state<number | null>(null);
  loadError = $state<string | null>(null);

  // Single in-flight promise so two simultaneous `load()` calls
  // share the same network request rather than racing each other.
  private inflight: Promise<ServerType[]> | null = null;

  /**
   * Load the catalog, idempotent within the soft TTL. The first call
   * always fetches; subsequent calls within TTL return the cached
   * result without touching the network. On error, the prior value
   * stays intact and the new error is recorded in `loadError`.
   */
  load(force = false): Promise<ServerType[]> {
    if (this.inflight) return this.inflight;

    const fresh = this.loadedAt !== null && Date.now() - this.loadedAt < TTL_MS;
    if (!force && fresh) {
      return Promise.resolve(this.types);
    }

    this.inflight = api
      .serverTypes()
      .then((items) => {
        this.types = items;
        this.loadedAt = Date.now();
        this.loadError = null;
        return items;
      })
      .catch((err: unknown) => {
        this.loadError = err instanceof Error ? err.message : String(err);
        throw err;
      })
      .finally(() => {
        this.inflight = null;
      });

    return this.inflight;
  }

  /** Look up a type by name. Returns undefined when the catalog hasn't loaded. */
  byName(name: string): ServerType | undefined {
    return this.types.find((t) => t.name === name);
  }

  /** Look up several types at once, preserving the input order. */
  byNames(names: string[]): (ServerType | undefined)[] {
    const map = new Map(this.types.map((t) => [t.name, t]));
    return names.map((n) => map.get(n));
  }

  /** Types that Hetzner currently reports as available (not sold out). */
  available(): ServerType[] {
    return this.types.filter((t) => t.available);
  }

  /**
   * Test helper — clears the cache + in-flight state so the next
   * `load()` re-fetches. Not used by application code.
   */
  reset(): void {
    this.types = [];
    this.loadedAt = null;
    this.loadError = null;
    this.inflight = null;
  }

  /** TEST-ONLY: reset the catalog to an empty state. */
  _reset(): void {
    this.types = [];
    this.loadedAt = null;
    this.loadError = null;
  }

  /** TEST-ONLY: replace the catalog with a fixed list. */
  _setTypesForTest(items: ServerType[]): void {
    this.types = items;
    this.loadedAt = Date.now();
    this.loadError = null;
  }
}

export const serverTypes = new ServerTypesStore();