<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import ServerTypeSelect from '$lib/components/ServerTypeSelect.svelte';

  let server = $state<Server | null>(null);
  let error = $state<string | null>(null);
  let saving = $state(false);

  // Local form state (mirrors the loaded server, with two-way bindings).
  // Fallback chain is stored as a CSV string for input ergonomics;
  // split on submit.
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
      error = err instanceof ApiError ? err.message : err instanceof Error ? err.message : String(err);
    } finally {
      saving = false;
    }
  }
</script>

<svelte:head>
  <title>{m.server_edit_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Edit server — single-column form. The fallback chain field is a CSV
  input rather than a chip-input; CSV matches the way operators
  actually copy/paste from documentation and avoids dragging in a
  combobox primitive.
-->
<div class="mx-auto max-w-2xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.server_edit_title()}
    </h1>
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <form onsubmit={submit} class="space-y-4">
    <div class="flex flex-col gap-1.5">
      <Label for="f-name">{m.server_edit_field_name()}</Label>
      <Input id="f-name" bind:value={form.name} required />
    </div>

    <div class="flex flex-col gap-1.5">
      <Label for="f-label">{m.server_edit_field_label()}</Label>
      <Input id="f-label" bind:value={form.label} />
    </div>

    <!--
      Base / top type are now ServerTypeSelect dropdowns (driven by
      /api/server-types) instead of free-text inputs. The fallback
      chain stays CSV because the design vocabulary deliberately
      prefers copy/paste from documentation over a chip combobox.
    -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <div class="flex flex-col gap-1.5">
        <Label for="f-base">{m.server_edit_field_base()}</Label>
        <ServerTypeSelect id="f-base" bind:value={form.base_server_type} required />
      </div>
      <div class="flex flex-col gap-1.5">
        <Label for="f-top">{m.server_edit_field_top()}</Label>
        <ServerTypeSelect id="f-top" bind:value={form.top_server_type} required />
      </div>
    </div>

    <div class="flex flex-col gap-1.5">
      <Label for="f-chain">{m.server_edit_field_fallback()}</Label>
      <Input
        id="f-chain"
        bind:value={form.fallback_chain_csv}
        required
        placeholder={m.server_edit_field_fallback_placeholder()}
      />
    </div>

    <div class="flex flex-col gap-1.5">
      <Label for="f-mode">{m.server_edit_field_mode()}</Label>
      <select
        id="f-mode"
        bind:value={form.mode}
        class="flex h-9 rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background"
      >
        <option value="manual">{m.servers_mode_manual()}</option>
        <option value="auto_promote">{m.servers_mode_auto_promote()}</option>
        <option value="scheduled">{m.servers_mode_scheduled()}</option>
      </select>
    </div>

    <div class="flex flex-col gap-1.5">
      <Label for="f-tz">{m.server_edit_field_timezone()}</Label>
      <Input
        id="f-tz"
        bind:value={form.timezone}
        required
        placeholder={m.server_edit_field_timezone_placeholder()}
      />
    </div>

    <div class="flex gap-2 border-t border-border pt-4">
      <Button variant="primary" type="submit" disabled={saving}>
        {saving ? m.server_edit_saving() : m.server_edit_save()}
      </Button>
      <Button variant="ghost" onclick={() => goto(`/servers/${serverId}`)}>
        {m.server_edit_cancel()}
      </Button>
    </div>
  </form>
</div>