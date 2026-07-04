<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import {
    Alert,
    Button,
    Card,
    Checkbox,
    Input,
    Label,
    Modal,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server, Window_ as Window } from '$lib/types';

  let server = $state<Server | null>(null);
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let saving = $state(false);

  let serverId = $derived(Number($page.params.id));

  // Modal state.
  let openModal = $state(false);

  // New-window form state (flat to match plan variable names).
  let newLabel = $state('');
  let newDays = $state(0b00111110);
  let newStart = $state('09:00');
  let newStop = $state('18:00');
  let newTarget = $state('');
  let newEnabled = $state(true);

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
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      saving = false;
    }
  }

  async function remove(id: number) {
    if (!confirm(m.windows_delete_confirm())) return;
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
        label: w.label,
        days_of_week: w.days_of_week,
        start_time: w.start_time,
        stop_time: w.stop_time,
        target_type: w.target_type,
        enabled: !w.enabled
      });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    }
  }
</script>

<div class="p-6 max-w-5xl mx-auto space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">
      {m.windows_title()} — {server?.name ?? '…'}
    </h1>
    <Button color="brand" onclick={() => (openModal = true)}>{m.windows_add()}</Button>
  </div>

  {#if error}
    <Alert color="danger">{error}</Alert>
  {/if}

  <Card>
    {#if windows.length === 0}
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.windows_empty()}</p>
    {:else}
      <Table hoverable>
        <TableHead>
          <TableHeadCell>{m.windows_col_label()}</TableHeadCell>
          <TableHeadCell>{m.windows_col_days()}</TableHeadCell>
          <TableHeadCell>{m.windows_col_start()}</TableHeadCell>
          <TableHeadCell>{m.windows_col_stop()}</TableHeadCell>
          <TableHeadCell>{m.windows_col_target()}</TableHeadCell>
          <TableHeadCell>{m.windows_col_enabled()}</TableHeadCell>
          <TableHeadCell><span class="sr-only">Actions</span></TableHeadCell>
        </TableHead>
        <TableBody>
          {#each windows as w (w.id)}
            <TableBodyRow>
              <TableBodyCell class="font-medium">{w.label}</TableBodyCell>
              <TableBodyCell class="text-xs">
                {dayLabels
                  .map((lbl, i) => (w.days_of_week & (1 << i) ? lbl : null))
                  .filter(Boolean)
                  .join(', ')}
              </TableBodyCell>
              <TableBodyCell>{w.start_time}</TableBodyCell>
              <TableBodyCell>{w.stop_time}</TableBodyCell>
              <TableBodyCell>{w.target_type}</TableBodyCell>
              <TableBodyCell>{w.enabled ? m.windows_col_yes() : m.windows_col_no()}</TableBodyCell>
              <TableBodyCell>
                <div class="flex gap-2">
                  <Button size="xs" color="alternative" onclick={() => toggleEnabled(w)}>
                    {w.enabled ? m.windows_disable() : m.windows_enable()}
                  </Button>
                  <Button size="xs" color="danger" onclick={() => remove(w.id)}
                    >{m.windows_delete()}</Button
                  >
                </div>
              </TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    {/if}
  </Card>
</div>

<Modal title={m.windows_modal_title()} bind:open={openModal} size="lg">
  <form onsubmit={create} class="space-y-3">
    <Label class="space-y-1">
      <span>{m.windows_field_label()}</span>
      <Input bind:value={newLabel} required placeholder="weekday-peak" />
    </Label>

    <div class="grid grid-cols-2 gap-3">
      <Label class="space-y-1">
        <span>{m.windows_field_start()}</span>
        <Input type="time" bind:value={newStart} required />
      </Label>
      <Label class="space-y-1">
        <span>{m.windows_field_stop()}</span>
        <Input type="time" bind:value={newStop} required />
      </Label>
    </div>

    <Label class="space-y-1">
      <span>{m.windows_field_target()}</span>
      <Input bind:value={newTarget} required placeholder="cpx31" />
    </Label>

    <div>
      <p class="text-sm mb-1 text-gray-900 dark:text-white">{m.windows_field_days()}</p>
      <div class="flex gap-1">
        {#each dayLabels as lbl, i}
          <button
            type="button"
            onclick={() => toggleDay(i)}
            class="rounded px-2 py-1 text-xs border border-gray-300 dark:border-gray-600"
            class:bg-blue-600={newDays & (1 << i)}
            class:text-white={newDays & (1 << i)}
            class:bg-white={!(newDays & (1 << i))}
            class:text-gray-700={!(newDays & (1 << i))}
            class:dark:bg-gray-800={!(newDays & (1 << i))}
            class:dark:text-gray-300={!(newDays & (1 << i))}
          >
            {lbl}
          </button>
        {/each}
      </div>
    </div>

    <Checkbox bind:checked={newEnabled}>{m.windows_field_enabled()}</Checkbox>

    <div class="flex justify-end gap-2">
      <Button color="alternative" onclick={() => (openModal = false)}
        >{m.windows_modal_cancel()}</Button
      >
      <Button type="submit" color="brand" disabled={saving}>
        {saving ? m.windows_modal_saving() : m.windows_modal_save()}
      </Button>
    </div>
  </form>
</Modal>