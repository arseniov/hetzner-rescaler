<script lang="ts">
  import { onMount } from 'svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { Server } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import Label from '$lib/components/ui/label.svelte';

  let servers = $state<Server[]>([]);
  // Live events stream — read from the SSE-backed store. The store is
  // process-wide (seeded by the Dashboard's initial REST fetch) and
  // fills in with new events as they arrive. The server-id filter and
  // kind filter below are applied client-side only.
  let events = $derived(eventsStream.events);
  let error = $state<string | null>(null);
  let loading = $state(true);

  // Selects are stateful native <select>s styled with our border /
  // radius tokens. No need for a shadcn Select primitive — the
  // native control already matches the dashboard's restrained tone.
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

  // The kinds the backend emits — keep in sync with store.Event.Kind.
  const kinds = [
    'rescale_up',
    'rescale_down',
    'rescale_fallback',
    'rescale_failed',
    'promote',
    'demote'
  ];
</script>

<svelte:head>
  <title>{m.events_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Events — full event log, filterable by server and kind. The live
  store keeps the list fresh via SSE; the filters apply on the
  client. Empty state and loading state are inline so the page never
  blanks out.
-->
<div class="mx-auto max-w-4xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.events_title()}
    </h1>
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <!-- Filter row. Native <select> styled to match the inputs. The
       fields are full-width on mobile and split 50/50 from sm up. -->
  <section
    aria-label="Filters"
    class="mb-6 grid grid-cols-1 gap-4 rounded-md border border-border bg-card p-4 sm:grid-cols-2"
  >
    <div class="flex flex-col gap-1.5">
      <Label for="filter-server">{m.events_filter_server()}</Label>
      <select
        id="filter-server"
        bind:value={serverFilter}
        class="flex h-9 rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background"
      >
        <option value="">{m.events_filter_server_all()}</option>
        {#each servers as s (s.id)}
          <option value={s.id}>{s.name}</option>
        {/each}
      </select>
    </div>
    <div class="flex flex-col gap-1.5">
      <Label for="filter-kind">{m.events_filter_kind()}</Label>
      <select
        id="filter-kind"
        bind:value={kindFilter}
        class="flex h-9 rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background"
      >
        <option value="">{m.events_filter_kind_all()}</option>
        {#each kinds as k}
          <option value={k}>{k}</option>
        {/each}
      </select>
    </div>
  </section>

  {#if loading}
    <p class="text-sm text-muted-foreground">{m.events_loading()}</p>
  {:else}
    <EventList events={filtered} limit={filtered.length} />
  {/if}
</div>