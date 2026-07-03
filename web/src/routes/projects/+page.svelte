<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api';
  import type { Project } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let projects = $state<Project[]>([]);
  let newName = $state('');
  let newToken = $state('');
  let error = $state<string | null>(null);
  let submitting = $state(false);
  let loading = $state(true);

  async function refresh() {
    loading = true;
    try {
      projects = await api.listProjects();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    error = null;
    submitting = true;
    try {
      await api.createProject({ name: newName.trim(), hcloud_token: newToken.trim() });
      newName = '';
      newToken = '';
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      submitting = false;
    }
  }

  async function remove(id: number) {
    if (!confirm('Delete this project and all its servers?')) return;
    try {
      await api.deleteProject(id);
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function refreshHetzner(id: number) {
    try {
      const r = await api.refreshProject(id);
      alert(`Added ${r.added.length} new server(s); ${r.skipped.length} already registered.`);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }
</script>

<div class="p-6 max-w-3xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold">Projects</h1>

  {#if error}<Alert variant="destructive">{error}</Alert>{/if}

  <form onsubmit={submit} class="space-y-3 rounded-md border border-border p-4">
    <h2 class="font-medium">Add a project</h2>
    <label class="block text-sm">
      Name
      <Input bind:value={newName} required placeholder="production" class="mt-1" />
    </label>
    <label class="block text-sm">
      Hetzner Cloud API token
      <Input type="password" bind:value={newToken} required placeholder="abc123…" class="mt-1" />
    </label>
    <Button type="submit" disabled={submitting}>{submitting ? 'Saving…' : 'Add project'}</Button>
  </form>

  {#if loading}
    <p class="text-muted-foreground">Loading…</p>
  {:else if projects.length === 0}
    <p class="text-sm text-muted-foreground">No projects yet.</p>
  {:else}
    <ul class="space-y-2">
      {#each projects as p (p.id)}
        <li class="flex items-center justify-between rounded-md border border-border p-3">
          <div>
            <a href="/projects/{p.id}" class="font-medium hover:underline">{p.name}</a>
            <span class="ml-2 text-xs text-muted-foreground">
              {p.has_token ? 'token stored' : 'no token'}
              {#if p.last_error} · error: {p.last_error}{/if}
            </span>
          </div>
          <div class="flex gap-2">
            <Button size="sm" variant="outline" onclick={() => refreshHetzner(p.id)}>Refresh from Hetzner</Button>
            <Button size="sm" variant="destructive" onclick={() => remove(p.id)}>Delete</Button>
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</div>
