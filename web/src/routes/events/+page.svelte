<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Label, Select, Input, Alert } from 'flowbite-svelte';
  import EventList from '$lib/components/EventList.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { RescaleEvent, Server } from '$lib/types';

  let servers = $state<Server[]>([]);
  // Live events stream — read from the SSE-backed store. The store is
  // process-wide (seeded by the Dashboard's initial REST fetch) and
  // fills in with new events as they arrive. The server-id filter and
  // kind filter below are applied client-side only.
  let events = $derived(eventsStream.events);
  let error = $state<string | null>(null);
  let loading = $state(true);

  let serverFilter = $state<number | ''>('');
  let kindFilter = $state<string>('');

  onMount(async () => {
    try {
      // Fetch the server list for the filter dropdown, and seed the
      // SSE store with a recent snapshot if it doesn't already have
      // anything. This keeps direct-navigation to /events working
      // when the user hasn't visited the Dashboard first.
      const [s, e] = await Promise.all([api.listServers(), api.globalEvents({ limit: 50 })]);
      servers = s;
      if (eventsStream.events.length === 0) {
        eventsStream.replaceAll(e);
      }
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  });

  let filtered = $derived(
    events
      .filter((e) => serverFilter === '' || e.server_id === Number(serverFilter))
      .filter((e) => !kindFilter || e.kind.includes(kindFilter))
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

  <Card class="border-0">
    <div class="grid gap-3 sm:grid-cols-2">
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
    </div>
  </Card>

  {#if loading}
    <p class="text-sm text-gray-600 dark:text-gray-400">{m.events_loading()}</p>
  {:else}
    <EventList events={filtered} limit={filtered.length} />
  {/if}
</div>
