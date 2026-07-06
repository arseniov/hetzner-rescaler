<script lang="ts">
  import type { Server } from '$lib/types';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { Status } from './StatusBadge.svelte';
  interface Props { server: Server }
  let { server }: Props = $props();

  // Current size — preferred source is the live `current_type` field
  // the API populates from Hetzner (`server.current_type`). When the
  // API didn't return it (Hetzner was unreachable, the row hasn't been
  // refreshed since the field was added), we fall back to the most
  // recent event's `to_type` and finally to the configured base type.
  // This is the same precedence used by /status/servers so a server
  // looks the same everywhere on the dashboard.
  let currentType = $derived.by(() => {
    if (server.current_type) return server.current_type;
    const latest = eventsStream.events.find((e) => e.server_id === server.id);
    return latest?.to_type ?? server.base_server_type;
  });

  // Health — same precedence: live `server.status` first (Hetzner's
  // authoritative view, e.g. `running` / `initializing` / `off`),
  // then event-derived ok/degraded from the live SSE stream, then
  // unknown when nothing has ever recorded a state.
  let status = $derived.by<Status>(() => {
    if (server.status) return server.status as Status;
    const recent = eventsStream.events.filter((e) => e.server_id === server.id);
    if (recent.length === 0) return 'unknown';
    if (recent.some((e) => !e.ok)) return 'degraded';
    return 'ok';
  });

  // Last activity timestamp — used as a relative-time hint. ISO is
  // stored verbatim; the user-facing label is computed in the markup.
  let lastActivity = $derived.by(() => {
    return eventsStream.events.find((e) => e.server_id === server.id)?.started_at;
  });

  function relativeTime(iso: string | undefined): string {
    if (!iso) return '—';
    const ms = Date.now() - new Date(iso).getTime();
    if (ms < 0) return 'just now';
    const sec = Math.floor(ms / 1000);
    if (sec < 60) return `${sec}s ago`;
    const min = Math.floor(sec / 60);
    if (min < 60) return `${min}m ago`;
    const hr = Math.floor(min / 60);
    if (hr < 24) return `${hr}h ago`;
    const day = Math.floor(hr / 24);
    return `${day}d ago`;
  }

  const modeLabel = {
    manual: 'Manual',
    auto_promote: 'Auto-promote',
    scheduled: 'Scheduled'
  } as const;

  // Mode-as-badge: filled primary for the only mode that drives
  // action on its own (auto-promote), bordered for the rest. Avoids
  // the visual noise of three identically-coloured chips.
  const modeClasses: Record<Server['mode'], string> = {
    manual: 'border-border bg-transparent text-muted-foreground',
    auto_promote: 'border-primary/40 bg-primary/10 text-primary',
    scheduled: 'border-border bg-muted text-foreground'
  };

  // Dot-colour map for every Status value StatusBadge renders. The
  // operator-facing vocabulary groups Hetzner live states into three
  // buckets: healthy (green), transitional (amber), stopped/destructive
  // (red/grey). When StatusBadge gains a new state this map must grow
  // to match — a fallback would silently render a missing colour.
  const statusClasses: Record<Status, string> = {
    ok: 'bg-success',
    degraded: 'bg-destructive',
    unknown: 'bg-muted-foreground/40',
    running: 'bg-success',
    initializing: 'bg-warning',
    starting: 'bg-warning',
    stopping: 'bg-warning',
    off: 'bg-muted-foreground/40',
    deleting: 'bg-destructive'
  };
</script>

<!--
  Horizontal server card. The page that hosts this card places the
  status column to the right; this component intentionally avoids
  nesting another column inside so the parent grid stays flat.
-->
<a
  href="/servers/{server.id}"
  class="flex items-center gap-3 rounded-md border border-border bg-card px-3 py-2.5 text-card-foreground transition-colors hover:bg-muted"
>
  <!-- Left: name + hcloud id (mono caption). The truncate on the
       second line keeps long names from breaking the grid. -->
  <div class="min-w-0 flex-1">
    <div class="truncate font-medium text-foreground">{server.name}</div>
    <div class="mt-0.5 truncate font-mono text-xs text-muted-foreground">
      hcloud #{server.hcloud_server_id}
      <span class="mx-1 text-foreground/30">·</span>
      {relativeTime(lastActivity)}
    </div>
  </div>

  <!-- Center: current server type — the field the operator most often
       wants at a glance. Mono so the size code lines up across rows. -->
  <div class="shrink-0 text-right">
    <div class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
      Current
    </div>
    <div class="font-mono text-sm font-semibold text-foreground tabular">
      {currentType}
    </div>
  </div>

  <!-- Right: mode as a small badge + status dot. The badge colour
       hints at automation level without screaming it. -->
  <div class="flex shrink-0 items-center gap-2">
    <span
      class="rounded-sm border px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-wider {modeClasses[server.mode]}"
    >
      {modeLabel[server.mode]}
    </span>
    <span
      class="inline-block size-1.5 rounded-full {statusClasses[status]}"
      aria-label="Status: {status}"
    ></span>
  </div>
</a>