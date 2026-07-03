<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api, ApiError } from '$lib/api';
  import type { Server, RescaleEvent } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import RescaleButton from '$lib/components/RescaleButton.svelte';
  import Button from '$lib/components/ui/button.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let server = $state<Server | null>(null);
  let events = $state<RescaleEvent[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  let serverId = $derived(Number($page.params.id));

  async function refresh() {
    loading = true;
    error = null;
    try {
      server = await api.getServer(serverId);
      events = await api.serverEvents(serverId, 25);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  async function promote() {
    try {
      await api.promote(serverId, { confirm: true });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    }
  }

  async function demote() {
    try {
      await api.demote(serverId, { confirm: true });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    }
  }
</script>

<div class="p-6 max-w-4xl mx-auto space-y-6">
  {#if loading}
    <p class="text-muted-foreground">Loading…</p>
  {:else if !server}
    <Alert variant="destructive">Server not found.</Alert>
  {:else}
    <div class="flex items-start justify-between">
      <div>
        <h1 class="text-3xl font-semibold">{server.name}</h1>
        <p class="text-sm text-muted-foreground">
          Hetzner #{server.hcloud_server_id} · mode: {server.mode}
          {#if server.promote_state} · state: {server.promote_state}{/if}
        </p>
      </div>
      <a href="/servers/{server.id}/edit"><Button variant="outline">Edit</Button></a>
    </div>

    {#if error}<Alert variant="destructive">{error}</Alert>{/if}

    <dl class="grid grid-cols-2 gap-3 text-sm rounded-md border border-border p-4">
      <dt class="text-muted-foreground">Base type</dt><dd>{server.base_server_type}</dd>
      <dt class="text-muted-foreground">Top type</dt><dd>{server.top_server_type}</dd>
      <dt class="text-muted-foreground">Fallback chain</dt>
      <dd>{server.fallback_chain.join(' → ')}</dd>
      <dt class="text-muted-foreground">Timezone</dt><dd>{server.timezone}</dd>
    </dl>

    <div class="flex flex-wrap gap-2">
      <RescaleButton {serverId} direction="up" label="Rescale up" onComplete={refresh} />
      <RescaleButton {serverId} direction="down" label="Rescale down" onComplete={refresh} />
      {#if server.mode === 'auto_promote'}
        <Button variant="outline" onclick={promote}>Request promote</Button>
        <Button variant="outline" onclick={demote}>Request demote</Button>
      {/if}
      <a href="/servers/{server.id}/windows">
        <Button variant="outline">Edit windows</Button>
      </a>
    </div>

    <section>
      <h2 class="text-lg font-medium mb-3">Recent events</h2>
      <EventList {events} limit={20} />
    </section>
  {/if}
</div>