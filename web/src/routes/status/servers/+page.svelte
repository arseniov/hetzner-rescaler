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

  // Walk the live event stream and tag each server that has at least
  // one failed event. Recomputes whenever the stream changes.
  const recentFailureByServer = $derived.by(() => {
    const map = new Map<number, boolean>();
    for (const ev of eventsStream.events) {
      if (!ev.ok) map.set(ev.server_id, true);
    }
    return map;
  });

  function statusFor(s: Server): Status {
    // Server.last_error is not exposed by the API today; classify
    // status purely from the live event stream.
    if (recentFailureByServer.get(s.id)) return 'degraded';
    if (eventsStream.events.some((e) => e.server_id === s.id)) return 'ok';
    return 'unknown';
  }

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
  Server status — flat list of every server with a status indicator
  derived from the live event stream. Mirrors the /servers table but
  drops everything that isn't triage-relevant (server types, project
  ID, mode). The operator comes here when something looks wrong on
  the dashboard.
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
        Header row with the same column widths as /servers so the
        operator's eye can swap between the two pages without
        re-learning the layout. Action: replaced project ID with a
        blank — it's not useful for triage.
      -->
      <div
        class="hidden border-b border-border px-4 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground sm:grid sm:grid-cols-[1fr_8rem_8rem_8rem] sm:gap-3"
      >
        <span>{m.servers_status_col_name()}</span>
        <span>{m.servers_status_col_mode()}</span>
        <span>{m.servers_status_col_top()}</span>
        <span>{m.servers_status_col_status()}</span>
      </div>
      <ul>
        {#each servers as s, i (s.id)}
          <li
            class="px-4 py-3 text-sm sm:grid sm:grid-cols-[1fr_8rem_8rem_8rem] sm:items-center sm:gap-3 {i > 0 ? 'border-t border-border' : ''}"
          >
            <a
              href="/servers/{s.id}"
              class="block truncate font-medium text-foreground hover:underline"
            >
              {s.name}
            </a>
            <span class="mt-1 block text-xs text-muted-foreground sm:mt-0">
              {s.mode}
            </span>
            <span class="mt-1 block font-mono text-xs text-muted-foreground sm:mt-0">
              {s.top_server_type}
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