<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api } from '$lib/api';
  import type { Project, Server } from '$lib/types';
  import ServerCard from '$lib/components/ServerCard.svelte';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let project = $state<Project | null>(null);
  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  let newHcloudId = $state<number | null>(null);
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
      newHcloudId = null;
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }
</script>

<div class="p-6 max-w-4xl mx-auto space-y-6">
  {#if loading}
    <p class="text-muted-foreground">Loading…</p>
  {:else if !project}
    <Alert variant="destructive">Project not found.</Alert>
  {:else}
    <h1 class="text-3xl font-semibold">{project.name}</h1>
    <p class="text-sm text-muted-foreground">
      token: {project.has_token ? 'stored' : 'missing'} ·
      created {new Date(project.created_at).toLocaleDateString()}
    </p>

    {#if error}<Alert variant="destructive">{error}</Alert>{/if}

    <form onsubmit={addServerManually} class="space-y-3 rounded-md border border-border p-4">
      <h2 class="font-medium">Register a server manually</h2>
      <label class="block text-sm">
        Hetzner server ID
        <Input type="number" bind:value={newHcloudId} required class="mt-1" />
      </label>
      <label class="block text-sm">
        Display name
        <Input bind:value={newName} required class="mt-1" />
      </label>
      <Button type="submit">Add server</Button>
      <p class="text-xs text-muted-foreground">
        Default base/top types are filled in. Edit them on the server detail page.
      </p>
    </form>

    <section>
      <h2 class="text-lg font-medium mb-3">Servers ({projectServers.length})</h2>
      {#if projectServers.length === 0}
        <p class="text-sm text-muted-foreground">No servers registered.</p>
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
