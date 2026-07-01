<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api';
  import type { RescaleEvent, Server } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

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

  $effect(() => { serverFilter; kindFilter; limit; refresh(); });

  let filtered = $derived(
    kindFilter ? events.filter((e) => e.kind.includes(kindFilter)) : events
  );

  const kinds = ['rescale_up', 'rescale_down', 'rescale_fallback', 'rescale_failed', 'promote', 'demote'];
</script>

<div class="p-6 max-w-4xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold">Events</h1>

  {#if error}<Alert variant="destructive">{error}</Alert>{/if}

  <div class="grid gap-3 sm:grid-cols-3 rounded-md border border-border p-4">
    <label class="block text-sm">
      Server
      <select bind:value={serverFilter} class="mt-1 flex h-10 w-full rounded-md border border-border bg-background px-3">
        <option value="">All servers</option>
        {#each servers as s (s.id)}
          <option value={s.id}>{s.name}</option>
        {/each}
      </select>
    </label>
    <label class="block text-sm">
      Kind
      <select bind:value={kindFilter} class="mt-1 flex h-10 w-full rounded-md border border-border bg-background px-3">
        <option value="">All kinds</option>
        {#each kinds as k}
          <option value={k}>{k}</option>
        {/each}
      </select>
    </label>
    <label class="block text-sm">
      Limit
      <Input type="number" bind:value={limit} min={1} max={500} class="mt-1" />
    </label>
  </div>

  {#if loading}
    <p class="text-muted-foreground">Loading…</p>
  {:else}
    <EventList events={filtered} limit={filtered.length} />
  {/if}
</div>