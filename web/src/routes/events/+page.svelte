<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Label, Select, Input, Alert } from 'flowbite-svelte';
  import EventList from '$lib/components/EventList.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import type { RescaleEvent, Server } from '$lib/types';

  let servers = $state<Server[]>([]);
  let events = $state<RescaleEvent[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  let serverFilter = $state<number | ''>('');
  let kindFilter = $state<string>('');
  let limit = $state(50);

  async function refresh() {
    loading = true;
    error = null;
    try {
      if (!servers.length) servers = await api.listServers();
      events = await api.globalEvents({
        serverId: serverFilter === '' ? undefined : Number(serverFilter),
        limit
      });
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  // Re-fetch when server or limit changes. kindFilter is purely client-side
  // (filtered via $derived below), so it does not need to trigger a refresh.
  $effect(() => {
    serverFilter;
    limit;
    refresh();
  });

  let filtered = $derived(
    kindFilter ? events.filter((e) => e.kind.includes(kindFilter)) : events
  );

  const kinds = [
    'rescale_up',
    'rescale_down',
    'rescale_fallback',
    'rescale_failed',
    'promote',
    'demote'
  ];
</script>

<div class="p-6 max-w-4xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{m.events_title()}</h1>

  {#if error}<Alert color="danger">{error}</Alert>{/if}

  <Card>
    <div class="grid gap-3 sm:grid-cols-3">
      <Label>
        {m.events_filter_server()}
        <Select bind:value={serverFilter} class="mt-1">
          <option value="">{m.events_filter_server_all()}</option>
          {#each servers as s (s.id)}
            <option value={s.id}>{s.name}</option>
          {/each}
        </Select>
      </Label>
      <Label>
        {m.events_filter_kind()}
        <Select bind:value={kindFilter} class="mt-1">
          <option value="">{m.events_filter_kind_all()}</option>
          {#each kinds as k}
            <option value={k}>{k}</option>
          {/each}
        </Select>
      </Label>
      <Label>
        {m.events_filter_limit()}
        <Input type="number" bind:value={limit} min={1} max={500} class="mt-1" />
      </Label>
    </div>
  </Card>

  {#if loading}
    <p class="text-sm text-gray-600 dark:text-gray-400">{m.events_loading()}</p>
  {:else}
    <EventList events={filtered} limit={filtered.length} />
  {/if}
</div>
