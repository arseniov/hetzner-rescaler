<script lang="ts">
  import { onMount } from 'svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { Server } from '$lib/types';
  import StatusBadge, { type Status } from '$lib/components/StatusBadge.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);

  // Two precomputed indices over the live event stream, recomputed
  // when the stream changes: (a) servers that have at least one
  // failed event → degraded, (b) the most-recent event per server
  // → drives both "current type" and "last activity" columns.
  const eventIndex = $derived.by(() => {
    const failed = new Set<number>();
    const latest = new Map<number, { started_at: string; to_type?: string; ok: boolean }>();
    for (const ev of eventsStream.events) {
      if (!ev.ok) failed.add(ev.server_id);
      const prev = latest.get(ev.server_id);
      if (!prev || prev.started_at < ev.started_at) {
        latest.set(ev.server_id, {
          started_at: ev.started_at,
          to_type: ev.to_type,
          ok: ev.ok
        });
      }
    }
    return { failed, latest };
  });

  function statusFor(s: Server): Status {
    if (eventIndex.failed.has(s.id)) return 'degraded';
    if (eventIndex.latest.has(s.id)) return 'ok';
    return 'unknown';
  }

  // "Current type" derivation: prefer the most recent event's
  // to_type (the operator's most recent reality), fall back to the
  // configured top type when the server has never rescaled (typical
  // for `manual` servers that have only ever sat at one shape).
  function currentTypeFor(s: Server): string {
    return eventIndex.latest.get(s.id)?.to_type ?? s.top_server_type;
  }

  function lastActivityFor(s: Server): string | undefined {
    return eventIndex.latest.get(s.id)?.started_at;
  }

  function relativeTime(iso: string | undefined): string {
    if (!iso) return m.servers_status_no_activity();
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

  // Mode-as-badge vocabulary. Kept in sync with ServerCard.
  const modeLabel: Record<Server['mode'], string> = {
    manual: 'Manual',
    auto_promote: 'Auto-promote',
    scheduled: 'Scheduled'
  };
  const modeClasses: Record<Server['mode'], string> = {
    manual: 'border-border bg-transparent text-muted-foreground',
    auto_promote: 'border-primary/40 bg-primary/10 text-primary',
    scheduled: 'border-border bg-muted text-foreground'
  };

  onMount(async () => {
    try {
      servers = await api.listServers();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }
  });
</script>

<svelte:head>
  <title>{m.servers_status_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Server status — flat list of every server, derived from the live
  event stream. Columns are sized for triage, not for inventory: name,
  current size (latest event's to_type, falls back to top), mode as a
  badge, last activity, and overall status. The operator comes here
  when something looks wrong on the dashboard.
-->
<div class="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.servers_status_title()}
    </h1>
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <section aria-label="Server status" class="rounded-md border border-border bg-card">
    {#if servers.length === 0}
      <p class="px-4 py-6 text-sm text-muted-foreground">{m.servers_status_empty()}</p>
    {:else}
      <!--
        Header row aligned with the data grid. Columns are tuned for
        triage: name (wide left), current size (mono), mode badge,
        last activity (mono), status (right). Below sm the row layout
        stacks; the header itself collapses so it doesn't leave a
        stranded label on mobile.
      -->
      <div
        class="hidden border-b border-border px-4 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground sm:grid sm:grid-cols-[1.4fr_6rem_8rem_7rem_8rem] sm:gap-3"
      >
        <span>{m.servers_status_col_name()}</span>
        <span>{m.servers_status_col_current()}</span>
        <span>{m.servers_status_col_mode()}</span>
        <span>{m.servers_status_col_last_activity()}</span>
        <span>{m.servers_status_col_status()}</span>
      </div>
      <ul>
        {#each servers as s, i (s.id)}
          <li
            class="px-4 py-3 text-sm sm:grid sm:grid-cols-[1.4fr_6rem_8rem_7rem_8rem] sm:items-center sm:gap-3 {i > 0 ? 'border-t border-border' : ''}"
          >
            <a
              href="/servers/{s.id}"
              class="block truncate font-medium text-foreground hover:underline"
            >
              {s.name}
            </a>
            <span class="mt-1 block font-mono text-xs font-semibold tabular text-foreground sm:mt-0">
              {currentTypeFor(s)}
            </span>
            <span class="mt-1 block sm:mt-0">
              <span
                class="inline-flex rounded-sm border px-1.5 py-0.5 font-mono text-[10px] uppercase tracking-wider {modeClasses[s.mode]}"
              >
                {modeLabel[s.mode]}
              </span>
            </span>
            <span class="mt-1 block font-mono text-xs text-muted-foreground sm:mt-0">
              {relativeTime(lastActivityFor(s))}
            </span>
            <span class="mt-1 block sm:mt-0">
              <StatusBadge status={statusFor(s)} />
            </span>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>