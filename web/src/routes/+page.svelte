<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api';
  import type { Project, Server, RescaleEvent } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import ServerCard from '$lib/components/ServerCard.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let projects = $state<Project[]>([]);
  let servers = $state<Server[]>([]);
  let events = $state<RescaleEvent[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  onMount(async () => {
    try {
      const [p, s, e] = await Promise.all([
        api.listProjects(),
        api.listServers(),
        api.globalEvents({ limit: 20 })
      ]);
      projects = p;
      servers = s;
      events = e;
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  });
</script>

<div class="p-6 space-y-6 max-w-5xl mx-auto">
  <h1 class="text-3xl font-semibold">Dashboard</h1>

  {#if error}
    <Alert variant="destructive">{error}</Alert>
  {/if}

  {#if loading}
    <p class="text-muted-foreground">Loading…</p>
  {:else}
    <section>
      <h2 class="text-lg font-medium mb-3">Projects ({projects.length})</h2>
      {#if projects.length === 0}
        <p class="text-sm text-muted-foreground">
          No projects yet. <a href="/projects" class="underline">Add one</a>.
        </p>
      {:else}
        <ul class="space-y-2">
          {#each projects as p (p.id)}
            <li class="rounded-md border border-border p-3">
              <a href="/projects/{p.id}" class="font-medium hover:underline">{p.name}</a>
              <span class="ml-2 text-xs text-muted-foreground">
                {p.has_token ? 'token stored' : 'no token'}
              </span>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <section>
      <h2 class="text-lg font-medium mb-3">Servers ({servers.length})</h2>
      {#if servers.length === 0}
        <p class="text-sm text-muted-foreground">No servers registered.</p>
      {:else}
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {#each servers.slice(0, 6) as s (s.id)}
            <ServerCard server={s} />
          {/each}
        </div>
      {/if}
    </section>

    <section>
      <h2 class="text-lg font-medium mb-3">Recent events</h2>
      <EventList {events} limit={10} />
    </section>
  {/if}
</div>
