<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Alert } from 'flowbite-svelte';
  import StatusBadge from '$lib/components/StatusBadge.svelte';
  import type { Status } from '$lib/components/StatusBadge.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { Server } from '$lib/types';

  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);

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

<div class="p-6 max-w-6xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">
    {m.servers_status_title()}
  </h1>

  {#if error}
    <Alert color="danger">{error}</Alert>
  {/if}

  <Card>
    {#if servers.length === 0}
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.servers_status_empty()}</p>
    {:else}
      <Table hoverable>
        <TableHead>
          <TableHeadCell>{m.servers_status_col_name()}</TableHeadCell>
          <TableHeadCell>{m.servers_status_col_mode()}</TableHeadCell>
          <TableHeadCell>{m.servers_status_col_top()}</TableHeadCell>
          <TableHeadCell>{m.servers_status_col_window()}</TableHeadCell>
          <TableHeadCell>{m.servers_status_col_status()}</TableHeadCell>
        </TableHead>
        <TableBody>
          {#each servers as s (s.id)}
            <TableBodyRow>
              <TableBodyCell>
                <a href="/servers/{s.id}" class="font-medium hover:underline text-blue-600 dark:text-blue-400">
                  {s.name}
                </a>
              </TableBodyCell>
              <TableBodyCell>{s.mode}</TableBodyCell>
              <TableBodyCell>{s.top_server_type}</TableBodyCell>
              <TableBodyCell>—</TableBodyCell>
              <TableBodyCell><StatusBadge status={statusFor(s)} /></TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    {/if}
  </Card>
</div>
