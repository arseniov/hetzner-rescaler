<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Alert } from 'flowbite-svelte';
  import { api } from '$lib/api';
  import { m } from '$lib/paraglide/messages.js';
  import type { Project, Server, RescaleEvent } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import ServerCard from '$lib/components/ServerCard.svelte';

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

<div class="p-6 space-y-6 max-w-6xl mx-auto">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{m.dashboard_title()}</h1>

  {#if error}
    <Alert color="red">{error}</Alert>
  {/if}

  {#if loading}
    <p class="text-gray-600 dark:text-gray-400">{m.dashboard_loading()}</p>
  {:else}
    <Card>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.dashboard_section_projects({ count: projects.length })}
      </h2>
      {#if projects.length === 0}
        <p class="text-sm text-gray-600 dark:text-gray-400">No projects yet.</p>
      {:else}
        <ul class="space-y-2">
          {#each projects as p (p.id)}
            <li class="rounded-md border border-gray-200 dark:border-gray-700 p-3">
              <a href="/projects/{p.id}" class="font-medium hover:underline">{p.name}</a>
              <span class="ml-2 text-xs text-gray-600 dark:text-gray-400">
                {p.has_token ? m.projects_token_stored() : m.projects_no_token()}
              </span>
            </li>
          {/each}
        </ul>
      {/if}
    </Card>

    <Card>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.dashboard_section_servers({ count: servers.length })}
      </h2>
      {#if servers.length === 0}
        <p class="text-sm text-gray-600 dark:text-gray-400">{m.servers_empty()}</p>
      {:else}
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {#each servers.slice(0, 6) as s (s.id)}
            <ServerCard server={s} />
          {/each}
        </div>
      {/if}
    </Card>

    <Card>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.dashboard_section_recent_events()}
      </h2>
      <EventList {events} limit={10} />
    </Card>
  {/if}
</div>
