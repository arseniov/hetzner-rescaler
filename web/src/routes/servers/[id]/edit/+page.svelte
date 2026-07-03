<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api, ApiError } from '$lib/api';
  import type { Server } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

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
        .split(',').map((s) => s.trim()).filter(Boolean);
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
  <h1 class="text-3xl font-semibold">Edit server</h1>
  {#if error}<Alert variant="destructive">{error}</Alert>{/if}

  <form onsubmit={submit} class="space-y-3">
    <label class="block text-sm">
      Name
      <Input bind:value={form.name} required class="mt-1" />
    </label>
    <label class="block text-sm">
      Label
      <Input bind:value={form.label} class="mt-1" />
    </label>
    <div class="grid grid-cols-2 gap-3">
      <label class="block text-sm">
        Base server type
        <Input bind:value={form.base_server_type} required class="mt-1" />
      </label>
      <label class="block text-sm">
        Top server type
        <Input bind:value={form.top_server_type} required class="mt-1" />
      </label>
    </div>
    <label class="block text-sm">
      Fallback chain (comma-separated, top first)
      <Input bind:value={form.fallback_chain_csv} required class="mt-1" placeholder="cpx31, cpx21, cpx11" />
    </label>
    <label class="block text-sm">
      Mode
      <select bind:value={form.mode} class="mt-1 flex h-10 w-full rounded-md border border-border bg-background px-3">
        <option value="manual">Manual</option>
        <option value="auto_promote">Auto-promote</option>
        <option value="scheduled">Scheduled</option>
      </select>
    </label>
    <label class="block text-sm">
      Timezone (IANA)
      <Input bind:value={form.timezone} required class="mt-1" placeholder="Europe/Rome" />
    </label>

    <div class="flex gap-2">
      <Button type="submit" disabled={saving}>{saving ? 'Saving…' : 'Save'}</Button>
      <Button variant="outline" onclick={() => goto(`/servers/${serverId}`)}>Cancel</Button>
    </div>
  </form>
</div>