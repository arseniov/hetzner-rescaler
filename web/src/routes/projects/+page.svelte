<script lang="ts">
  import { onMount } from 'svelte';
  import { Plus, Trash2 } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Project } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import Dialog from '$lib/components/ui/dialog.svelte';

  let projects = $state<Project[]>([]);
  let error = $state<string | null>(null);

  // Add-project dialog state. The dialog is bound on `open` so the
  // form can dismiss itself after a successful submit.
  let addOpen = $state(false);
  let newName = $state('');
  let newToken = $state('');
  let busy = $state(false);

  // Pending-deletion state. We don't use the browser `confirm()` —
  // it blocks the thread, has no styling, and breaks with browsers
  // that throttle background tabs. Instead, the row's delete button
  // arms an inline confirmation; pressing it again within 3s commits.
  let pendingDeleteId = $state<number | null>(null);
  let deleteTimer: ReturnType<typeof setTimeout> | null = null;

  function armDelete(id: number) {
    pendingDeleteId = id;
    if (deleteTimer) clearTimeout(deleteTimer);
    deleteTimer = setTimeout(() => {
      pendingDeleteId = null;
      deleteTimer = null;
    }, 3000);
  }

  function cancelDelete() {
    pendingDeleteId = null;
    if (deleteTimer) {
      clearTimeout(deleteTimer);
      deleteTimer = null;
    }
  }

  async function refresh() {
    try {
      projects = await api.listProjects();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }
  }

  async function create(e: SubmitEvent) {
    e.preventDefault();
    busy = true;
    error = null;
    try {
      await api.createProject({ name: newName.trim(), hcloud_token: newToken.trim() });
      newName = '';
      newToken = '';
      addOpen = false;
      await refresh();
    } catch (e) {
      error = e instanceof ApiError ? e.message : e instanceof Error ? e.message : String(e);
    } finally {
      busy = false;
    }
  }

  async function commitDelete(p: Project) {
    pendingDeleteId = null;
    if (deleteTimer) {
      clearTimeout(deleteTimer);
      deleteTimer = null;
    }
    try {
      await api.deleteProject(p.id);
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }
  }
</script>

<svelte:head>
  <title>{m.projects_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Projects — list of registered Hetzner Cloud projects plus an "Add
  project" entry point that opens a dialog with the create form.

  Layout: a single flat panel. Rows are separated by hairline borders
  (border-t on items past the first), not by per-row cards. The header
  holds the title and the add button, justified to opposite edges.
-->
<div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6 flex items-center justify-between gap-3">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.projects_title()}
    </h1>
    <Button variant="primary" size="md" onclick={() => (addOpen = true)}>
      <Plus class="size-4" strokeWidth={1.75} aria-hidden="true" />
      {m.projects_new_submit()}
    </Button>
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <section aria-label="Projects" class="rounded-md border border-border bg-card">
    {#if projects.length === 0}
      <p class="px-4 py-6 text-sm text-muted-foreground">{m.projects_empty()}</p>
    {:else}
      <ul>
        {#each projects as p, i (p.id)}
          {@const armed = pendingDeleteId === p.id}
          <li
            class="flex items-center gap-3 px-4 py-3 text-sm {i > 0 ? 'border-t border-border' : ''}"
          >
            <a
              href="/projects/{p.id}"
              class="min-w-0 flex-1 truncate font-medium text-foreground hover:underline"
            >
              {p.name}
            </a>
            <span class="hidden font-mono text-xs uppercase tracking-wider text-muted-foreground sm:inline">
              {p.has_token ? m.projects_token_stored() : m.projects_no_token()}
            </span>
            {#if p.last_error}
              <span class="hidden truncate text-xs text-destructive md:inline" title={p.last_error}>
                {p.last_error}
              </span>
            {/if}
            <span class="hidden font-mono text-xs text-muted-foreground md:inline">
              {new Date(p.updated_at).toLocaleDateString()}
            </span>
            {#if armed}
              <button
                type="button"
                onclick={cancelDelete}
                class="font-mono text-xs uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
              >
                Cancel
              </button>
              <button
                type="button"
                onclick={() => commitDelete(p)}
                class="inline-flex h-8 items-center gap-1.5 rounded-sm border border-destructive/30 bg-destructive/10 px-3 text-xs font-medium text-destructive transition-colors hover:bg-destructive hover:text-destructive-foreground"
              >
                <Trash2 class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                Confirm
              </button>
            {:else}
              <button
                type="button"
                onclick={() => armDelete(p.id)}
                aria-label="{m.projects_delete_submit()} {p.name}"
                class="inline-flex size-8 items-center justify-center rounded-sm text-muted-foreground hover:bg-destructive/10 hover:text-destructive transition-colors"
              >
                <Trash2 class="size-4" strokeWidth={1.5} aria-hidden="true" />
              </button>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>

<!--
  Add-project dialog. The form collects the human-friendly name and
  the Hetzner Cloud API token. The token field is `type="password"` so
  the value isn't echoed back; better-auth-friendly server-side
  validation catches malformed tokens.
-->
<Dialog
  bind:open={addOpen}
  title={m.projects_new_submit()}
  size="md"
>
  <form id="add-project-form" onsubmit={create} class="space-y-4">
    <div class="flex flex-col gap-1.5">
      <Label for="project-name">{m.projects_new_name_label()}</Label>
      <Input
        id="project-name"
        bind:value={newName}
        required
        placeholder="production-east"
        autocomplete="off"
      />
    </div>
    <div class="flex flex-col gap-1.5">
      <Label for="project-token">{m.projects_new_token_label()}</Label>
      <Input
        id="project-token"
        type="password"
        bind:value={newToken}
        required
        autocomplete="off"
        placeholder="hc_••••••••••••••••"
      />
      <p class="text-xs text-muted-foreground">
        Token is stored encrypted; used to call Hetzner Cloud for server refresh.
      </p>
    </div>
  </form>

  {#snippet footer()}
    <Button variant="ghost" onclick={() => (addOpen = false)} disabled={busy}>
      Cancel
    </Button>
    <Button variant="primary" type="submit" form="add-project-form" disabled={busy}>
      {busy ? '…' : m.projects_new_submit()}
    </Button>
  {/snippet}
</Dialog>