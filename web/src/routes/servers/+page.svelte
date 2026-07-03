<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api';
  import type { Server } from '$lib/types';
  import ServerCard from '$lib/components/ServerCard.svelte';

  let servers = $state<Server[]>([]);
  let loading = $state(true);

  onMount(async () => {
    try { servers = await api.listServers(); } finally { loading = false; }
  });
</script>

<div class="p-6 max-w-5xl mx-auto space-y-4">
  <h1 class="text-3xl font-semibold">Servers</h1>
  {#if loading}
    <p class="text-muted-foreground">Loading…</p>
  {:else if servers.length === 0}
    <p class="text-sm text-muted-foreground">No servers registered.</p>
  {:else}
    <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
      {#each servers as s (s.id)}
        <ServerCard server={s} />
      {/each}
    </div>
  {/if}
</div>