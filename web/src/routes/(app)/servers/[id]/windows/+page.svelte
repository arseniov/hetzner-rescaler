<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { ArrowLeft, Plus, Trash2 } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server, Window_ as Window } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import Dialog from '$lib/components/ui/dialog.svelte';
  import ServerTypeSelect from '$lib/components/ServerTypeSelect.svelte';

  let server = $state<Server | null>(null);
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let saving = $state(false);

  let serverId = $derived(Number($page.params.id));

  // Add-window dialog state.
  let openModal = $state(false);

  // New-window form state (flat to match plan variable names).
  let newLabel = $state('');
  let newDays = $state(0b00111110);
  let newStart = $state('09:00');
  let newStop = $state('18:00');
  let newTarget = $state('');
  let newEnabled = $state(true);

  // Pending-deletion state. Two-tap pattern (same as /projects):
  // first tap arms the action; a 3 s timer or cancel disarms.
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

  const dayLabels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

  function resetForm() {
    newLabel = '';
    newDays = 0b00111110;
    newStart = '09:00';
    newStop = '18:00';
    newTarget = '';
    newEnabled = true;
  }

  async function refresh() {
    try {
      server = await api.getServer(serverId);
      windows = await api.listWindows(serverId);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  onMount(refresh);

  function toggleDay(d: number) {
    newDays ^= (1 << d);
  }

  async function create(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    error = null;
    try {
      await api.createWindow(serverId, {
        label: newLabel.trim(),
        days_of_week: newDays,
        start_time: newStart,
        stop_time: newStop,
        target_type: newTarget.trim(),
        enabled: newEnabled
      });
      resetForm();
      openModal = false;
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : err instanceof Error ? err.message : String(err);
    } finally {
      saving = false;
    }
  }

  async function commitDelete(id: number) {
    pendingDeleteId = null;
    if (deleteTimer) {
      clearTimeout(deleteTimer);
      deleteTimer = null;
    }
    try {
      await api.deleteWindow(id);
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : err instanceof Error ? err.message : String(err);
    }
  }

  async function toggleEnabled(w: Window) {
    try {
      await api.updateWindow(w.id, {
        label: w.label,
        days_of_week: w.days_of_week,
        start_time: w.start_time,
        stop_time: w.stop_time,
        target_type: w.target_type,
        enabled: !w.enabled
      });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : err instanceof Error ? err.message : String(err);
    }
  }

  // Format the day-of-week bitmask for display.
  function daysLabel(mask: number): string {
    const on = dayLabels.filter((_, i) => mask & (1 << i));
    return on.join(', ') || '—';
  }
</script>

<svelte:head>
  <title>{m.windows_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Windows — scheduled rescale rules for a server. The header carries
  the server name (mono caption) plus a back link; the +Add button
  sits on the right. The list itself is a flat panel with row
  separators; the add-form lives behind a dialog.
-->
<div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6 flex flex-wrap items-end justify-between gap-3">
    <div>
      <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
        {m.windows_title()}
      </h1>
      <p class="mt-1 font-mono text-xs text-muted-foreground">
        {server?.name ?? '…'}
      </p>
    </div>
    <div class="flex gap-2">
      <Button variant="ghost" size="sm" href="/servers/{serverId}">
        <ArrowLeft class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
        {m.server_detail_edit()}
      </Button>
      <Button variant="primary" size="md" onclick={() => (openModal = true)}>
        <Plus class="size-4" strokeWidth={1.75} aria-hidden="true" />
        {m.windows_add()}
      </Button>
    </div>
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <section aria-label="Windows" class="rounded-md border border-border bg-card">
    {#if windows.length === 0}
      <p class="px-4 py-6 text-sm text-muted-foreground">{m.windows_empty()}</p>
    {:else}
      <!--
        Header row aligned with the data grid below. Days / start /
        stop / target / enabled are smaller columns; the label takes
        the wider left edge.
      -->
      <div
        class="hidden border-b border-border px-4 py-2 text-[11px] font-medium uppercase tracking-wider text-muted-foreground lg:grid lg:grid-cols-[1.5fr_2fr_5rem_5rem_6rem_5rem_7rem] lg:gap-3"
      >
        <span>{m.windows_col_label()}</span>
        <span>{m.windows_col_days()}</span>
        <span>{m.windows_col_start()}</span>
        <span>{m.windows_col_stop()}</span>
        <span>{m.windows_col_target()}</span>
        <span>{m.windows_col_enabled()}</span>
        <span class="text-right">Actions</span>
      </div>
      <ul>
        {#each windows as w, i (w.id)}
          {@const armed = pendingDeleteId === w.id}
          <li
            class="flex flex-wrap items-center gap-2 px-4 py-3 text-sm lg:grid lg:grid-cols-[1.5fr_2fr_5rem_5rem_6rem_5rem_7rem] lg:items-center lg:gap-3 {i > 0 ? 'border-t border-border' : ''}"
          >
            <span class="min-w-0 flex-1 truncate font-medium text-foreground lg:flex-none">
              {w.label}
            </span>
            <span class="font-mono text-xs text-muted-foreground lg:block">
              {daysLabel(w.days_of_week)}
            </span>
            <span class="font-mono text-xs text-foreground">{w.start_time}</span>
            <span class="font-mono text-xs text-foreground">{w.stop_time}</span>
            <span class="font-mono text-xs text-foreground">{w.target_type}</span>
            <span class="inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-wider {w.enabled ? 'text-success' : 'text-muted-foreground'}">
              <span
                class="inline-block size-1.5 rounded-full {w.enabled
                  ? 'bg-success'
                  : 'bg-muted-foreground/40'}"
              ></span>
              {w.enabled ? m.windows_col_yes() : m.windows_col_no()}
            </span>
            <div class="ml-auto flex items-center gap-2 lg:ml-0 lg:justify-end">
              <button
                type="button"
                onclick={() => toggleEnabled(w)}
                class="font-mono text-xs uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
              >
                {w.enabled ? m.windows_disable() : m.windows_enable()}
              </button>
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
                  onclick={() => commitDelete(w.id)}
                  class="inline-flex h-7 items-center gap-1.5 rounded-sm border border-destructive/30 bg-destructive/10 px-2 text-xs font-medium text-destructive transition-colors hover:bg-destructive hover:text-destructive-foreground"
                >
                  <Trash2 class="size-3" strokeWidth={1.75} aria-hidden="true" />
                  Confirm
                </button>
              {:else}
                <button
                  type="button"
                  onclick={() => armDelete(w.id)}
                  aria-label="{m.windows_delete()} {w.label}"
                  class="inline-flex size-7 items-center justify-center rounded-sm text-muted-foreground hover:bg-destructive/10 hover:text-destructive transition-colors"
                >
                  <Trash2 class="size-3.5" strokeWidth={1.5} aria-hidden="true" />
                </button>
              {/if}
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>

<!--
  Add-window dialog. Form is contained inside the dialog body; the
  footer holds the cancel / save actions so they're always visible
  regardless of how tall the form grows.
-->
<Dialog bind:open={openModal} title={m.windows_modal_title()} size="lg">
  <form id="add-window-form" onsubmit={create} class="space-y-4">
    <div class="flex flex-col gap-1.5">
      <Label for="w-label">{m.windows_field_label()}</Label>
      <Input id="w-label" bind:value={newLabel} required placeholder="weekday-peak" />
    </div>

    <div class="grid grid-cols-2 gap-3">
      <div class="flex flex-col gap-1.5">
        <Label for="w-start">{m.windows_field_start()}</Label>
        <Input id="w-start" type="time" bind:value={newStart} required />
      </div>
      <div class="flex flex-col gap-1.5">
        <Label for="w-stop">{m.windows_field_stop()}</Label>
        <Input id="w-stop" type="time" bind:value={newStop} required />
      </div>
    </div>

    <div class="flex flex-col gap-1.5">
      <Label for="w-target">{m.windows_field_target()}</Label>
      <!--
        ServerTypeSelect replaces the old free-text target_type input
        so the operator can't typo a Hetzner type code into the form.
        `required` enforces a non-empty selection before submit.
      -->
      <ServerTypeSelect id="w-target" bind:value={newTarget} server={server} location={server?.location} required />
    </div>

    <!--
      Day-of-week selector. Bitmask under the hood (Sunday = bit 0);
      each button toggles its bit. We use the "outlined + filled"
      pattern: an inactive day is a hairline outline; an active day
      inverts to filled foreground — the same vocabulary used for the
      segmented control on the dashboard.
    -->
    <div class="flex flex-col gap-1.5">
      <span class="text-sm font-medium text-foreground">{m.windows_field_days()}</span>
      <div class="flex gap-1">
        {#each dayLabels as lbl, i}
          {@const on = !!(newDays & (1 << i))}
          <button
            type="button"
            onclick={() => toggleDay(i)}
            class="h-8 w-10 rounded-sm font-mono text-xs uppercase tracking-wider transition-colors {on
              ? 'bg-primary text-primary-foreground border border-primary'
              : 'border border-border bg-transparent text-muted-foreground hover:text-foreground'}"
            aria-pressed={on}
          >
            {lbl}
          </button>
        {/each}
      </div>
    </div>

    <label class="flex items-center gap-2 text-sm text-foreground">
      <input
        type="checkbox"
        bind:checked={newEnabled}
        class="size-4 rounded-sm border border-border bg-input text-primary focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
      />
      {m.windows_field_enabled()}
    </label>
  </form>

  {#snippet footer()}
    <Button variant="ghost" onclick={() => (openModal = false)} disabled={saving}>
      {m.windows_modal_cancel()}
    </Button>
    <Button variant="primary" type="submit" form="add-window-form" disabled={saving}>
      {saving ? m.windows_modal_saving() : m.windows_modal_save()}
    </Button>
  {/snippet}
</Dialog>