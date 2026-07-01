<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { api, ApiError } from '$lib/api';
  import type { Server, Window_ as Window } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let server = $state<Server | null>(null);
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let saving = $state(false);

  let serverId = $derived(Number($page.params.id));

  // Form state for the "new window" form.
  let form = $state({
    label: '',
    days_of_week: 0b00111110,
    start_time: '09:00',
    stop_time: '18:00',
    target_type: '',
    enabled: true
  });

  const dayLabels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

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
    form.days_of_week ^= (1 << d);
  }

  async function addWindow(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    error = null;
    try {
      await api.createWindow(serverId, {
        label: form.label.trim(),
        days_of_week: form.days_of_week,
        start_time: form.start_time,
        stop_time: form.stop_time,
        target_type: form.target_type.trim(),
        enabled: form.enabled
      });
      form = { label: '', days_of_week: 0b00111110, start_time: '09:00', stop_time: '18:00', target_type: '', enabled: true };
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      saving = false;
    }
  }

  async function removeWindow(id: number) {
    if (!confirm('Delete this window?')) return;
    try {
      await api.deleteWindow(id);
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    }
  }

  async function toggleEnabled(w: Window) {
    try {
      await api.updateWindow(w.id, {
        label: w.label, days_of_week: w.days_of_week,
        start_time: w.start_time, stop_time: w.stop_time,
        target_type: w.target_type, enabled: !w.enabled
      });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    }
  }
</script>

<div class="p-6 max-w-3xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold">Windows for {server?.name ?? '…'}</h1>
  {#if error}<Alert variant="destructive">{error}</Alert>{/if}

  <form onsubmit={addWindow} class="space-y-3 rounded-md border border-border p-4">
    <h2 class="font-medium">New window</h2>
    <label class="block text-sm">
      Label
      <Input bind:value={form.label} required class="mt-1" placeholder="weekday-peak" />
    </label>
    <div class="flex gap-3 text-sm">
      <label class="block">Start<Input type="time" bind:value={form.start_time} required class="mt-1" /></label>
      <label class="block">Stop<Input type="time" bind:value={form.stop_time} required class="mt-1" /></label>
      <label class="block">Target type<Input bind:value={form.target_type} required class="mt-1" placeholder="cpx31" /></label>
    </div>
    <div>
      <p class="text-sm mb-1">Days of week</p>
      <div class="flex gap-1">
        {#each dayLabels as label, i}
          <button
            type="button"
            onclick={() => toggleDay(i)}
            class="rounded px-2 py-1 text-xs border border-border"
            class:bg-primary={form.days_of_week & (1 << i)}
            class:text-primary-foreground={form.days_of_week & (1 << i)}
          >{label}</button>
        {/each}
      </div>
    </div>
    <label class="flex items-center gap-2 text-sm">
      <input type="checkbox" bind:checked={form.enabled} /> Enabled
    </label>
    <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Add window'}</Button>
  </form>

  <section>
    <h2 class="text-lg font-medium mb-3">Existing windows ({windows.length})</h2>
    {#if windows.length === 0}
      <p class="text-sm text-muted-foreground">No windows yet.</p>
    {:else}
      <ul class="space-y-2">
        {#each windows as w (w.id)}
          <li class="flex items-center justify-between rounded-md border border-border p-3 text-sm">
            <div>
              <span class="font-medium">{w.label}</span>
              <span class="ml-2 text-muted-foreground">
                {w.start_time}–{w.stop_time} → {w.target_type}
              </span>
            </div>
            <div class="flex gap-2">
              <Button size="sm" variant="outline" onclick={() => toggleEnabled(w)}>
                {w.enabled ? 'Disable' : 'Enable'}
              </Button>
              <Button size="sm" variant="destructive" onclick={() => removeWindow(w.id)}>Delete</Button>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </section>
</div>