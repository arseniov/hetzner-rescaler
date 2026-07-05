<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Alert, Badge } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import type { Server } from '$lib/types';
  import StatusBadge from '$lib/components/StatusBadge.svelte';

  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);

  const modeLabel = {
    manual: m.servers_mode_manual,
    auto_promote: m.servers_mode_auto_promote,
    scheduled: m.servers_mode_scheduled
  } as const;

  onMount(async () => {
    try { servers = await api.listServers(); }
    catch (e) { error = e instanceof Error ? e.message : String(e); }
  });
</script>

<div class="p-6 max-w-6xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{m.servers_title()}</h1>
  {#if error}<Alert color="danger">{error}</Alert>{/if}

  <Card class="border-0">
    {#if servers.length === 0}
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.servers_empty()}</p>
    {:else}
      <Table hoverable>
        <TableHead>
          <TableHeadCell>{m.servers_col_name()}</TableHeadCell>
          <TableHeadCell>{m.servers_col_project()}</TableHeadCell>
          <TableHeadCell>{m.servers_col_types()}</TableHeadCell>
          <TableHeadCell>{m.servers_col_mode()}</TableHeadCell>
          <TableHeadCell>{m.servers_col_status()}</TableHeadCell>
        </TableHead>
        <TableBody>
          {#each servers as s (s.id)}
            {@const status = 'ok'}
            <TableBodyRow>
              <TableBodyCell><a href="/servers/{s.id}" class="font-medium text-blue-600 dark:text-blue-400 hover:underline">{s.name}</a></TableBodyCell>
              <TableBodyCell>{s.project_id}</TableBodyCell>
              <TableBodyCell>{s.base_server_type} → {s.top_server_type}</TableBodyCell>
              <TableBodyCell><Badge color="gray">{modeLabel[s.mode]()}</Badge></TableBodyCell>
              <TableBodyCell><StatusBadge {status} /></TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    {/if}
  </Card>
</div>