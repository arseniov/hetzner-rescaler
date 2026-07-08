import type { Server, ServerType } from '$lib/types';

/**
 * Role a server type plays for a given server. Exactly one role per
 * type (priority order if a type happens to be in multiple buckets):
 *
 *   current > base > top > fallback
 *
 * The current chip wins because it answers "what is this server
 * *actually* running on right now", which is the operator's first
 * question when picking a new size.
 */
export type ServerTypeRole = 'current' | 'base' | 'top' | 'fallback';

export function roleFor(t: ServerType, server: Server | undefined | null): ServerTypeRole | null {
  if (!server) return null;
  if (server.current_type && t.name === server.current_type) return 'current';
  if (t.name === server.base_server_type) return 'base';
  if (t.name === server.top_server_type) return 'top';
  if (server.fallback_chain?.includes(t.name)) return 'fallback';
  return null;
}