export interface Project {
  id: number;
  name: string;
  has_token: boolean;
  last_error?: string;
  created_at: string;
  updated_at: string;
}

export interface Server {
  id: number;
  project_id: number;
  hcloud_server_id: number;
  name: string;
  label: string;
  base_server_type: string;
  top_server_type: string;
  fallback_chain: string[];
  mode: 'manual' | 'auto_promote' | 'scheduled';
  promote_state?: string | null;
  timezone: string;
  // Live state from Hetzner — populated when the API call succeeds.
  // Both fields are absent (undefined) when Hetzner is unreachable or
  // the server has been deleted out-of-band. The web falls back to
  // event-derived state in that case so the dashboard never blanks.
  status?: string;
  current_type?: string;
  // Live Hetzner-reported location (e.g. `fsn1`, `nbg1`, `hel1`).
  // Populated when the API call to Hetzner succeeds; absent when
  // Hetzner is unreachable or the server has been deleted out-of-band.
  // The frontend uses this to query /api/server-types?location=X for
  // the per-location availability gate; when it's missing the frontend
  // falls back to the unfiltered catalog so dropdowns never blank.
  location?: string;
  created_at?: string;
  updated_at?: string;
  // In-flight rescale_pending event, if any. Mirrors the
  // rescale_pending row in the events table — present only while a
  // rescale is running. The server detail page hydrates from this on
  // mount (survives a hard refresh) and then keeps it fresh via the
  // pendingRescale store fed by SSE.
  pending_event?: RescaleEvent;
}

export interface Window_ {
  id: number;
  server_id: number;
  label: string;
  days_of_week: number;
  start_time: string;
  stop_time: string;
  target_type: string;
  enabled: boolean;
}

export interface RescaleEvent {
  id: number;
  server_id: number;
  kind: string;
  from_type?: string;
  to_type?: string;
  // Current phase of the rescale. Set on rescale_pending rows
  // (`shutting_down`, `changing_type`, `powering_on`); undefined on
  // terminal rows (rescale_completed / rescale_failed). The server
  // detail badge uses this to label the live state.
  phase?: string;
  started_at: string;
  finished_at: string;
  ok: boolean;
  error?: string;
  triggered_by: string;
}

/**
 * ServerType — one row of GET /api/server-types. `name` is the Hetzner
 * type code (e.g. `cpx11`, `cpx31`, `cax11`); the other fields drive
 * the operator-facing dropdowns (description, cores, memory, disk)
 * and the cost chart (price). The full shape is exposed because
 * ServerTypeSelect uses the cores/memory/price to render a richer
 * option label than just the type code.
 */
export interface ServerType {
  name: string;
  available: boolean;
  description?: string;
  cores?: number;
  memory_gb?: number;
  disk_gb?: number;
  price_monthly_eur?: number;
}

export interface CreateProjectRequest { name: string; hcloud_token: string }

/**
 * CreateProjectResult — what POST /api/projects returns. `added` and
 * `skipped` are the auto-sync tally the server performs on create
 * (every Hetzner server it could see, partitioned into ones it linked
 * vs. ones that were already linked to another project). `last_error`
 * is set when the initial sync failed partway; the UI surfaces it as
 * a non-destructive warning so the operator knows the project exists
 * but the servers may be incomplete.
 *
 * `added` and `skipped` are nullable on the wire: when the initial
 * sync fails before the loop completes (most commonly an invalid
 * Hetzner token), the Go handler returns `null` for both arrays and
 * surfaces the failure through `last_error` instead. Consumers should
 * treat null as "no sync data" — see /projects for the pattern.
 */
export interface CreateProjectResult {
  project: Project;
  added: Server[] | null;
  skipped: Server[] | null;
  last_error?: string;
}
export interface CreateServerRequest {
  project_id: number;
  hcloud_server_id: number;
  name: string;
  label: string;
  base_server_type: string;
  top_server_type: string;
  fallback_chain: string[];
  mode: Server['mode'];
  timezone: string;
}
export interface UpdateServerRequest {
  name: string;
  label: string;
  base_server_type: string;
  top_server_type: string;
  fallback_chain: string[];
  mode: Server['mode'];
  timezone: string;
}
export interface CreateWindowRequest {
  label: string;
  days_of_week: number;
  start_time: string;
  stop_time: string;
  target_type: string;
  enabled: boolean;
}
export interface RescaleRequest { direction: 'up' | 'down'; confirm: boolean }
export interface ConfirmRequest { confirm: boolean }
export interface RefreshResponse { added: Server[]; skipped: Server[] }

/**
 * Metrics types — JSON shapes returned by GET /api/metrics. Field names
 * are snake_case to match the Go handler's JSON tags directly (the rest
 * of the Go API is snake_case; mixing camelCase here silently dropped
 * every KPI to "—"). See internal/api/handlers_metrics.go.
 */
export interface RescaleCountsByDayRow {
  date: string;
  ok: number;
  failed: number;
  total: number;
}
export interface HoursAtTypeRow {
  server_id: number;
  server_name: string;
  base: number;
  top: number;
  fallback: number;
  cost_eur: number;
}
export interface SuccessRateRow {
  server_id: number;
  server_name: string;
  ok: number;
  total: number;
  ok_rate: number;
}
export interface MetricsKpis {
  active_server_count: number;
  projects_with_token_count: number;
  rescales_24h_ok: number;
  last_rescale_error: null | RescaleError;
}
export interface RescaleError {
  server_id: number;
  kind: string;
  at: string;
  error: string;
}
export interface MetricsResponse {
  range: '1d' | '7d' | '30d';
  from: string;
  to: string;
  kpis: MetricsKpis;
  rescale_counts_by_day: RescaleCountsByDayRow[];
  hours_at_type: HoursAtTypeRow[];
  success_rate_by_server: SuccessRateRow[];
}
