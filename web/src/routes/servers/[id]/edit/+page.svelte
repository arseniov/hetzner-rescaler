<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Alert, Button, Input, Label, Select } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server } from '$lib/types';

  let server = $state<Server | null>(null);
  let error = $state<string | null>(null);
  let saving = $state(false);

  // Local form state (mirrors the loaded server, with two-way bindings).
  let form = $state({
    name: '',
    label: '',
    base_server_type: '',
    top_server_type: '',
    fallback_chain_csv: '',
    mode: 'manual' as Server['mode'],
    timezone: 'UTC'
  });

  let serverId = $derived(Number($page.params.id));

  onMount(async () => {
    try {
      server = await api.getServer(serverId);
      form = {
        name: server.name,
        label: server.label,
        base_server_type: server.base_server_type,
        top_server_type: server.top_server_type,
        fallback_chain_csv: server.fallback_chain.join(','),
        mode: server.mode,
        timezone: server.timezone
      };
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  });

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    error = null;
    try {
      const chain = form.fallback_chain_csv
        .split(',')
        .map((s) => s.trim())
        .filter(Boolean);
      const updated = await api.updateServer(serverId, {
        name: form.name.trim(),
        label: form.label.trim(),
        base_server_type: form.base_server_type.trim(),
        top_server_type: form.top_server_type.trim(),
        fallback_chain: chain,
        mode: form.mode,
        timezone: form.timezone.trim()
      });
      server = updated;
      await goto(`/servers/${serverId}`);
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      saving = false;
    }
  }
</script>

<div class="p-6 max-w-2xl mx-auto space-y-4">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">
    {m.server_edit_title()}
  </h1>
  {#if error}
    <Alert color="danger">{error}</Alert>
  {/if}

  <form onsubmit={submit} class="space-y-3">
    <Label class="space-y-1">
      <span>{m.server_edit_field_name()}</span>
      <Input bind:value={form.name} required />
    </Label>
    <Label class="space-y-1">
      <span>{m.server_edit_field_label()}</span>
      <Input bind:value={form.label} />
    </Label>
    <div class="grid grid-cols-2 gap-3">
      <Label class="space-y-1">
        <span>{m.server_edit_field_base()}</span>
        <Input bind:value={form.base_server_type} required />
      </Label>
      <Label class="space-y-1">
        <span>{m.server_edit_field_top()}</span>
        <Input bind:value={form.top_server_type} required />
      </Label>
    </div>
    <Label class="space-y-1">
      <span>{m.server_edit_field_fallback()}</span>
      <Input
        bind:value={form.fallback_chain_csv}
        required
        placeholder={m.server_edit_field_fallback_placeholder()}
      />
    </Label>
    <Label class="space-y-1">
      <span>{m.server_edit_field_mode()}</span>
      <Select bind:value={form.mode}>
        <option value="manual">{m.servers_mode_manual()}</option>
        <option value="auto_promote">{m.servers_mode_auto_promote()}</option>
        <option value="scheduled">{m.servers_mode_scheduled()}</option>
      </Select>
    </Label>
    <Label class="space-y-1">
      <span>{m.server_edit_field_timezone()}</span>
      <Input bind:value={form.timezone} required placeholder={m.server_edit_field_timezone_placeholder()} />
    </Label>

    <div class="flex gap-2">
      <Button type="submit" color="brand" disabled={saving}>
        {saving ? m.server_edit_saving() : m.server_edit_save()}
      </Button>
      <Button color="alternative" onclick={() => goto(`/servers/${serverId}`)}>
        {m.server_edit_cancel()}
      </Button>
    </div>
  </form>
</div>