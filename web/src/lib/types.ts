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
  started_at: string;
  finished_at: string;
  ok: boolean;
  error?: string;
  triggered_by: string;
}

export interface ServerType { name: string; available: boolean }

export interface CreateProjectRequest { name: string; hcloud_token: string }
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

export interface RescaleCountsByDayRow { date: string; ok: number; failed: number; total: number }
export interface HoursAtTypeRow { serverId: number; serverName: string; base: number; top: number; fallback: number; costEur: number }
export interface SuccessRateRow { serverId: number; serverName: string; ok: number; total: number; okRate: number }
export interface MetricsKpis {
  activeServerCount: number;
  projectsWithTokenCount: number;
  rescales24hOk: number;
  lastRescaleError: null | RescaleError;
}
export interface RescaleError {
  serverId: number;
  kind: string;
  at: string;
  error: string;
}
export interface MetricsResponse {
  range: '1d' | '7d' | '30d';
  from: string;
  to: string;
  kpis: MetricsKpis;
  rescaleCountsByDay: RescaleCountsByDayRow[];
  hoursAtType: HoursAtTypeRow[];
  successRateByServer: SuccessRateRow[];
}
