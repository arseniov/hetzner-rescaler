<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { Alert, Button, Card, Input, Label } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import type { Project, Server } from '$lib/types';
  import ServerCard from '$lib/components/ServerCard.svelte';

  let project = $state<Project | null>(null);
  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  let newHcloudId = $state<number | undefined>(undefined);
  let newName = $state('');

  let projectId = $derived(Number($page.params.id));

  async function refresh() {
    loading = true;
    error = null;
    try {
      const all = await api.listProjects();
      project = all.find((p) => p.id === projectId) ?? null;
      servers = await api.listServers();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  let projectServers = $derived(servers.filter((s) => s.project_id === projectId));

  async function addServerManually(e: SubmitEvent) {
    e.preventDefault();
    if (!newHcloudId) return;
    error = null;
    try {
      await api.createServer({
        project_id: projectId,
        hcloud_server_id: newHcloudId,
        name: newName,
        label: newName,
        base_server_type: 'cpx11',
        top_server_type: 'cpx31',
        fallback_chain: ['cpx31', 'cpx11'],
        mode: 'manual',
        timezone: 'UTC'
      });
      newName = '';
      newHcloudId = undefined;
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }
</script>

<div class="p-6 max-w-5xl mx-auto space-y-6">
  {#if loading}
    <p class="text-sm text-gray-600 dark:text-gray-400">{m.project_detail_loading()}</p>
  {:else if !project}
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">
        {m.project_detail_title()}
      </h1>
      <Button color="alternative" href="/projects">{m.project_detail_back()}</Button>
    </div>
    <Alert color="danger">{m.project_detail_not_found()}</Alert>
  {:else}
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{project.name}</h1>
        <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
          token: {project.has_token
            ? m.project_detail_token_stored()
            : m.project_detail_token_missing()}
          ·
          {m.project_detail_created_at({
            date: new Date(project.created_at).toLocaleDateString()
          })}
        </p>
      </div>
      <Button color="alternative" href="/projects">{m.project_detail_back()}</Button>
    </div>

    {#if error}
      <Alert color="danger">{error}</Alert>
    {/if}

    <Card class="border-0">
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.project_detail_register_title()}
      </h2>
      <form onsubmit={addServerManually} class="space-y-3">
        <Label>
          {m.project_detail_hcloud_id_label()}
          <Input type="number" bind:value={newHcloudId} required class="mt-1" />
        </Label>
        <Label>
          {m.project_detail_name_label()}
          <Input bind:value={newName} required class="mt-1" />
        </Label>
        <Button type="submit">{m.project_detail_add_submit()}</Button>
        <p class="text-xs text-gray-600 dark:text-gray-400">{m.project_detail_add_hint()}</p>
      </form>
    </Card>

    <section>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.project_detail_servers_title({ count: projectServers.length })}
      </h2>
      {#if projectServers.length === 0}
        <p class="text-sm text-gray-600 dark:text-gray-400">
          {m.project_detail_servers_empty()}
        </p>
      {:else}
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {#each projectServers as s (s.id)}
            <ServerCard server={s} />
          {/each}
        </div>
      {/if}
    </section>
  {/if}
</div>
