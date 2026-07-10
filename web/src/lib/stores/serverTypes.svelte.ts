import { api } from '$lib/api';
import type { ServerType } from '$lib/types';

// Soft TTL — 5 minutes. The server-types catalog rarely changes (a new
// type appears maybe quarterly, and per-location availability only flips
// on Hetzner stock events), so re-fetching on every page mount is wasteful.
const TTL_MS = 5 * 60 * 1000;

interface LocationCache {
  types: ServerType[];
  loadedAt: number;
}

class ServerTypesStore {
  // The most recently loaded catalog (whatever location it was for).
  // Pages that don't care about location read this; pages that need a
  // specific location call `load(location)` which keeps the cache fresh
  // for that location.
  types = $state<ServerType[]>([]);
  loadedAt = $state<number | null>(null);
  loadError = $state<string | null>(null);

  // Per-location cache. Keyed by location name; the most recent
  // successful load for each location sticks around for one TTL window.
  private byLocation = new Map<string, LocationCache>();

  // Single in-flight promise so two simultaneous `load(loc)` calls
  // share the same network request rather than racing each other.
  // The key includes the location so simultaneous loads for DIFFERENT
  // locations each get their own in-flight promise.
  private inflight = new Map<string, Promise<ServerType[]>>();

  /**
   * Load the catalog for a specific Hetzner location (e.g. "fsn1").
   * Idempotent within the soft TTL: a second call for the same
   * location within 5 minutes returns the cached result. Calls for a
   * different location refetch (the cache is per-location).
   *
   * On error, the prior value stays intact and the new error is
   * recorded in `loadError`. The store's `types` always reflects the
   * last successful load for any location.
   */
  load(location: string, force = false): Promise<ServerType[]> {
    if (!location) {
      return Promise.reject(new Error('location required'));
    }
    const existing = this.inflight.get(location);
    if (existing) return existing;

    const cached = this.byLocation.get(location);
    const fresh = cached !== undefined && Date.now() - cached.loadedAt < TTL_MS;
    if (!force && fresh) {
      return Promise.resolve(cached!.types);
    }

    const promise = api
      .serverTypes(location)
      .then((items) => {
        this.byLocation.set(location, { types: items, loadedAt: Date.now() });
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
        this.inflight.delete(location);
      });

    this.inflight.set(location, promise);
    return promise;
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

  /** Types that Hetzner currently reports as available. */
  available(): ServerType[] {
    return this.types.filter((t) => t.available);
  }

  /**
   * Test helper — clears the cache + in-flight state so the next
   * `load()` re-fetches for every location.
   */
  reset(): void {
    this.types = [];
    this.loadedAt = null;
    this.loadError = null;
    this.byLocation.clear();
    this.inflight.clear();
  }

  /** TEST-ONLY: reset the catalog to an empty state. */
  _reset(): void {
    this.types = [];
    this.loadedAt = null;
    this.loadError = null;
    this.byLocation.clear();
    this.inflight.clear();
  }

  /** TEST-ONLY: replace the catalog with a fixed list (no location set). */
  _setTypesForTest(items: ServerType[]): void {
    this.types = items;
    this.loadedAt = Date.now();
    this.loadError = null;
  }
}

export const serverTypes = new ServerTypesStore();
