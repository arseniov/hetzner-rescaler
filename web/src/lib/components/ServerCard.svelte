<script lang="ts">
  import type { Server } from '$lib/types';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  interface Props { server: Server }
  let { server }: Props = $props();

  // Current size = the `to_type` of the most recent event for this
  // server. Falls back to the configured base type when no event
  // has been recorded yet (e.g. a brand-new manual server). This is
  // the same heuristic used by /status/servers and is intentionally
  // permissive — exact reconciliation against Hetzner's live state
  // is a server-detail concern.
  let currentType = $derived.by(() => {
    const latest = eventsStream.events.find((e) => e.server_id === server.id);
    return latest?.to_type ?? server.base_server_type;
  });

  // Health: any failed event in the recent stream → degraded,
  // otherwise ok (events seen and all succeeded), or unknown (no
  // events yet — silent operator or new server).
  let status = $derived.by<'ok' | 'degraded' | 'unknown'>(() => {
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

  const statusClasses: Record<'ok' | 'degraded' | 'unknown', string> = {
    ok: 'bg-success',
    degraded: 'bg-destructive',
    unknown: 'bg-muted-foreground/40'
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