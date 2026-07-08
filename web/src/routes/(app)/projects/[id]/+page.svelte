<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { ArrowLeft, Plus } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Project, Server } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import Dialog from '$lib/components/ui/dialog.svelte';
  import ServerCard from '$lib/components/ServerCard.svelte';
  import ServerTypeSelect from '$lib/components/ServerTypeSelect.svelte';
  import ServerTypeMultiSelect from '$lib/components/ServerTypeMultiSelect.svelte';

  let project = $state<Project | null>(null);
  let servers = $state<Server[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);

  // Register-server dialog. The form used to live inline on the page;
  // it was a constant visual tax on a project that already has its
  // servers. Now it opens behind an "Add server" button next to the
  // Servers section title — same affordance as "Add project" on /projects.
  // Form fields live in component state and are reset on close so the
  // next invocation starts blank.
  let registerOpen = $state(false);
  let newHcloudId = $state<string>('');
  let newName = $state('');
  let newBase = $state('');
  let newTop = $state('');
  let newFallback = $state<string[]>([]);
  let registering = $state(false);
  let registerError = $state<string | null>(null);

  function resetRegisterForm() {
    newHcloudId = '';
    newName = '';
    newBase = '';
    newTop = '';
    newFallback = [];
    registerError = null;
  }

  function openRegisterDialog() {
    resetRegisterForm();
    registerOpen = true;
  }

  function closeRegisterDialog() {
    registerOpen = false;
    resetRegisterForm();
  }

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
    registerError = null;
    registering = true;
    try {
      await api.createServer({
        project_id: projectId,
        hcloud_server_id: Number(newHcloudId),
        name: newName,
        label: newName,
        // Operator-driven base/top/fallback — no magic defaults. We
        // require the operator to pick because silently registering
        // a server with an empty type would create a row the API
        // can never rescale.
        base_server_type: newBase,
        top_server_type: newTop,
        fallback_chain: newFallback,
        mode: 'manual',
        timezone: 'UTC'
      });
      registerOpen = false;
      resetRegisterForm();
      await refresh();
    } catch (err) {
      registerError = err instanceof ApiError ? err.message : err instanceof Error ? err.message : String(err);
    } finally {
      registering = false;
    }
  }
</script>

<svelte:head>
  <title>{project?.name ?? m.project_detail_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Project detail — three sections: identity header, servers list (with
  an "Add server" affordance that opens the register dialog), and the
  register dialog itself. The header carries the token status and
  creation date as a single muted line; we don't repeat the name as a
  separate badge.
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

    <!--
      Servers in this project. The list mirrors the dashboard's server
      list vocabulary (ServerCard) so a server looks the same
      everywhere. The section header carries the count on the left and
      the "Add server" button on the right — same idiom as "Add project"
      on /projects. Clicking it opens the register dialog below.
    -->
    <section aria-label="Servers in this project">
      <header class="mb-3 flex items-center justify-between gap-3">
        <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          {m.project_detail_servers_title({ count: projectServers.length })}
        </h2>
        <Button variant="default" size="sm" onclick={openRegisterDialog}>
          <Plus class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
          {m.project_detail_add_server_button()}
        </Button>
      </header>
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

<!--
  Register-server dialog. The hcloud_server_id input is number-coerced
  (string ↔ number) so the field can be empty during typing without
  "NaN" appearing in the value. The base/top/fallback chain fields are
  operator-driven dropdowns (was: hardcoded 'cpx11'/'cpx31' defaults
  that the operator had to fix on the server's edit page later). We
  intentionally leave them empty rather than auto-pick — silent
  defaults hide mistakes, and the dropdown makes the choice cheap.

  The dialog is sized `lg` so all five fields fit on one screen without
  scrolling on desktop; on mobile the content area scrolls naturally.
-->
<Dialog
  bind:open={registerOpen}
  title={m.project_detail_register_dialog_title()}
  description={m.project_detail_register_dialog_description()}
  size="lg"
>
  <form id="register-server-form" onsubmit={registerServer} class="space-y-4">
    {#if registerError}
      <Alert variant="destructive">{registerError}</Alert>
    {/if}

    <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
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
    </div>

    <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
      <div class="flex flex-col gap-1.5">
        <Label for="reg-base">{m.project_detail_field_base()}</Label>
        <ServerTypeSelect id="reg-base" bind:value={newBase} required />
      </div>
      <div class="flex flex-col gap-1.5">
        <Label for="reg-top">{m.project_detail_field_top()}</Label>
        <ServerTypeSelect id="reg-top" bind:value={newTop} required />
      </div>
      <div class="flex flex-col gap-1.5">
        <Label for="reg-fallback">{m.project_detail_field_fallback()}</Label>
        <ServerTypeMultiSelect
          id="reg-fallback"
          bind:value={newFallback}
          excluded={[newBase, newTop].filter(Boolean)}
        />
      </div>
    </div>

    <p class="text-xs text-muted-foreground">{m.project_detail_add_hint()}</p>
  </form>

  {#snippet footer()}
    <Button variant="ghost" onclick={closeRegisterDialog} disabled={registering}>
      Cancel
    </Button>
    <Button variant="primary" type="submit" form="register-server-form" disabled={registering}>
      {registering ? '…' : m.project_detail_add_submit()}
    </Button>
  {/snippet}
</Dialog>