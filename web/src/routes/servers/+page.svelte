<script lang="ts">
  import { onMount } from 'svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import type { Server } from '$lib/types';
  import Alert from '$lib/components/ui/alert.svelte';
  import StatusBadge, { type Status } from '$lib/components/StatusBadge.svelte';

  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);

  // The mode column is a fixed vocabulary: manual / auto-promote /
  // scheduled. Map each to a localised label so the column is uniform
  // across rows.
  const modeLabel = {
    manual: m.servers_mode_manual,
    auto_promote: m.servers_mode_auto_promote,
    scheduled: m.servers_mode_scheduled
  } as const;

  // /servers is the inventory list, not the triage list (that's
  // /status/servers). We show the API's `current_type` so the
  // operator can see at a glance which size each server is actually
  // running, with the configured base → top range alongside.
  function statusFor(s: Server): Status {
    return (s.status as Status) ?? 'unknown';
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
  <title>{m.servers_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Servers — flat list of every server the rescaler knows about. Each
  row deep-links to the server's detail page. The status column now
  reads the API's `server.status` (populated live from Hetzner); the
  Current column shows the live `server.current_type` so the operator
  sees what size is actually running, not just the configured top.
-->
<div class="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.servers_title()}
    </h1>
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <section aria-label="Servers" class="rounded-md border border-border bg-card">
    {#if servers.length === 0}
      <p class="px-4 py-6 text-sm text-muted-foreground">{m.servers_empty()}</p>
    {:else}
      <!--
        Header row is a separate block from the data rows so the
        column boundaries line up across both. The header uses the
        same hairline underline as the dashboard's KPI panel.
      -->
      <div
        class="hidden border-b border-border px-4 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground sm:grid sm:grid-cols-[1fr_5rem_7rem_7rem_7rem] sm:gap-3"
      >
        <span>{m.servers_col_name()}</span>
        <span class="text-right">{m.servers_col_project()}</span>
        <span>{m.servers_col_types()}</span>
        <span>Current</span>
        <span>{m.servers_col_status()}</span>
      </div>
      <ul>
        {#each servers as s, i (s.id)}
          {@const status = statusFor(s)}
          <li
            class="px-4 py-3 text-sm sm:grid sm:grid-cols-[1fr_5rem_7rem_7rem_7rem] sm:items-center sm:gap-3 {i > 0 ? 'border-t border-border' : ''}"
          >
            <a
              href="/servers/{s.id}"
              class="block truncate font-medium text-foreground hover:underline"
            >
              {s.name}
            </a>
            <span class="mt-1 block font-mono text-xs text-muted-foreground sm:mt-0 sm:text-right">
              #{s.project_id}
            </span>
            <span class="mt-1 block truncate font-mono text-xs text-muted-foreground sm:mt-0">
              {s.base_server_type} <span class="text-foreground/40">→</span> {s.top_server_type}
            </span>
            <span class="mt-1 block font-mono text-xs font-semibold tabular text-foreground sm:mt-0">
              {s.current_type ?? s.top_server_type}
            </span>
            <span class="mt-1 block sm:mt-0">
              <StatusBadge {status} />
            </span>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>