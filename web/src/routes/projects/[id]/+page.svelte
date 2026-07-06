<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { ArrowLeft } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Project, Server } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import ServerCard from '$lib/components/ServerCard.svelte';

  let project = $state<Project | null>(null);
  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  // Inline register-server form state. We keep this on the page
  // rather than behind a dialog because it's a short, primary
  // action for the project view.
  let newHcloudId = $state<string>('');
  let newName = $state('');
  let registering = $state(false);

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

  async function registerServer(e: SubmitEvent) {
    e.preventDefault();
    if (!newHcloudId) return;
    error = null;
    registering = true;
    try {
      await api.createServer({
        project_id: projectId,
        hcloud_server_id: Number(newHcloudId),
        name: newName,
        label: newName,
        // Manual registration: defaults match the project's
        // fallback chain (cpx11 ↔ cpx31) — operators can adjust
        // from the server's edit page after registering.
        base_server_type: 'cpx11',
        top_server_type: 'cpx31',
        fallback_chain: ['cpx31', 'cpx11'],
        mode: 'manual',
        timezone: 'UTC'
      });
      newName = '';
      newHcloudId = '';
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : err instanceof Error ? err.message : String(err);
    } finally {
      registering = false;
    }
  }
</script>

<svelte:head>
  <title>{project?.name ?? m.project_detail_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Project detail — one project, three sections: identity header,
  register-server form, list of servers in this project. The header
  carries the token status and creation date as a single muted line;
  we don't repeat the name as a separate badge.
-->
<div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
  {#if loading}
    <p class="text-sm text-muted-foreground">{m.project_detail_loading()}</p>
  {:else if !project}
    <header class="mb-6 flex items-center justify-between gap-3">
      <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
        {m.project_detail_title()}
      </h1>
      <Button variant="ghost" size="sm" href="/projects">
        <ArrowLeft class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
        {m.project_detail_back()}
      </Button>
    </header>
    <Alert variant="destructive">{m.project_detail_not_found()}</Alert>
  {:else}
    <header class="mb-6 flex flex-wrap items-end justify-between gap-3">
      <div class="min-w-0">
        <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
          {project.name}
        </h1>
        <p class="mt-1 font-mono text-xs text-muted-foreground">
          token: {project.has_token
            ? m.project_detail_token_stored()
            : m.project_detail_token_missing()}
          <span class="mx-2 text-foreground/30">·</span>
          {m.project_detail_created_at({
            date: new Date(project.created_at).toLocaleDateString()
          })}
        </p>
      </div>
      <Button variant="ghost" size="sm" href="/projects">
        <ArrowLeft class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
        {m.project_detail_back()}
      </Button>
    </header>

    {#if error}
      <Alert variant="destructive" class="mb-6">{error}</Alert>
    {/if}

    <!-- Register-server form. Inline; same hairline panel vocabulary
         as the rest of the page. The hcloud_server_id input is
         number-coerced (string ↔ number) so the field can be empty
         during typing without "NaN" appearing in the value. -->
    <section
      aria-label="Register a server manually"
      class="mb-6 rounded-md border border-border bg-card p-4"
    >
      <h2 class="mb-3 font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
        {m.project_detail_register_title()}
      </h2>
      <form onsubmit={registerServer} class="grid grid-cols-1 gap-3 sm:grid-cols-[10rem_1fr_auto] sm:items-end">
        <div class="flex flex-col gap-1.5">
          <Label for="hcloud-id">{m.project_detail_hcloud_id_label()}</Label>
          <Input
            id="hcloud-id"
            type="number"
            bind:value={newHcloudId}
            required
            placeholder="12345678"
          />
        </div>
        <div class="flex flex-col gap-1.5">
          <Label for="hcloud-name">{m.project_detail_name_label()}</Label>
          <Input id="hcloud-name" bind:value={newName} required placeholder="web-1" />
        </div>
        <Button variant="primary" type="submit" disabled={registering}>
          {registering ? '…' : m.project_detail_add_submit()}
        </Button>
      </form>
      <p class="mt-2 text-xs text-muted-foreground">{m.project_detail_add_hint()}</p>
    </section>

    <!-- Servers in this project. The list mirrors the dashboard's
         server list vocabulary (ServerCard) so a server looks the
         same everywhere. -->
    <section aria-label="Servers in this project">
      <h2 class="mb-3 font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
        {m.project_detail_servers_title({ count: projectServers.length })}
      </h2>
      {#if projectServers.length === 0}
        <p class="text-sm text-muted-foreground">{m.project_detail_servers_empty()}</p>
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