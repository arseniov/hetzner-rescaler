<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Calendar } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import { serverTypes } from '$lib/stores/serverTypes.svelte';
  import { nextWindow, type WindowSpec } from '$lib/utils/windowSchedule';
  import type { Server, Window_ as Window } from '$lib/types';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import ServerTypeSelect from '$lib/components/ServerTypeSelect.svelte';
  import ServerTypeMultiSelect from '$lib/components/ServerTypeMultiSelect.svelte';

  let server = $state<Server | null>(null);
  // Server's existing windows. Fetched alongside the server so the
  // "first-window" banner can render meaningful info, not just a generic
  // warning that fires for every scheduled-mode edit (including ones with
  // dozens of well-defined windows). Populated in onMount below.
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let saving = $state(false);

  // Local form state (mirrors the loaded server, with two-way bindings).
  // Fallback chain is an ordered string[] — bound directly to the
  // drag-drop ServerTypeMultiSelect, no CSV conversion needed.
  let form = $state({
    name: '',
    label: '',
    base_server_type: '',
    top_server_type: '',
    fallback_chain: [] as string[],
    mode: 'manual' as Server['mode'],
    timezone: 'UTC'
  });

  let serverId = $derived(Number($page.params.id));

  // Default location for the catalog load. Used when /api/servers/[id]
  // doesn't return a `location` field — which DOES happen for any
  // server where Hetzner's GetServer call soft-fails (token problems,
  // the server was deleted out-of-band, Datacenter is nil in the Hetzner
  // response). In that case `liveServerState` returns the zero-value
  // LiveServerState and the JSON tag omitempty drops the field, so
  // `server.location` is undefined on the wire. Confirmed 2026-07-10
  // by inspecting a real API payload:
  //   {"id":2,"project_id":2,"hcloud_server_id":139131129,
  //    "name":"fiutaspesa-app",..., "status":"running",
  //    "current_type":"cx33","created_at":"...", "updated_at":"..."}
  // — note the missing `location` key. Without this default, the
  // /api/server-types?location=X call NEVER fires for such servers.
  // `fsn1` matches the hardcoded default used by the project page's
  // register-server dialog (see projects/[id]/+page.svelte).
  const FALLBACK_LOCATION = 'fsn1';

  onMount(async () => {
    // Fire the per-location catalog load IMMEDIATELY on mount, BEFORE
    // awaiting the server. This guarantees the network request goes
    // out on every navigation to /servers/[id]/edit — which is the
    // contract the type-availability gate relies on. The store's
    // in-flight map dedupes overlapping calls; the per-location TTL
    // cache keeps repeat visits cheap.
    serverTypes.load(FALLBACK_LOCATION).catch(() => {
      /* loadError is set on the store */
    });

    try {
      server = await api.getServer(serverId);
      // If the server resolves with a real location (the typical case),
      // refire the load with it. Store-level dedup means a back-to-back
      // load for the SAME location returns the cached result with no
      // extra network call; a different location triggers a fresh
      // fetch that's also cached for 5min. The dropdown is already
      // showing types from FALLBACK_LOCATION so it has something to
      // render immediately; the refire upgrades it to per-location
      // availability as soon as the server data arrives.
      if (server.location && server.location !== FALLBACK_LOCATION) {
        serverTypes.load(server.location).catch(() => {
          /* loadError is set on the store */
        });
      }
      // Fetch windows in parallel to the form setup — the "first-window"
      // banner for scheduled mode needs them. windowSchedule.nextWindow
      // is a pure derivation so the spec test for it covers the logic;
      // here we only wire the data flow.
      const fetched = await api.listWindows(serverId).catch(() => []);
      windows = Array.isArray(fetched) ? fetched : [];
      form = {
        name: server.name,
        label: server.label,
        base_server_type: server.base_server_type,
        top_server_type: server.top_server_type,
        fallback_chain: [...server.fallback_chain],
        mode: server.mode,
        timezone: server.timezone
      };
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  });

  // Pure derivation: collapse fetched windows into the WindowSpec shape
  // that nextWindow expects. Re-runs whenever `windows` or `form.timezone`
  // changes, so editing the timezone field re-evaluates the banner with
  // the new wall-clock reference without a remount.
  let firstWindow = $derived(
    nextWindow(
      windows.map<WindowSpec>((w) => ({
        days_of_week: w.days_of_week,
        start_time: w.start_time,
        stop_time: w.stop_time,
        target_type: w.target_type,
        enabled: w.enabled
      })),
      form.timezone || 'UTC',
      new Date()
    )
  );

  // Shared formatter: same Intl shape as ModePill so the banner copy
  // and the mode pill read as one vocabulary across pages.
  const fmts = new Map<string, Intl.DateTimeFormat>();
  function fmtWhen(d: Date, tz: string): string {
    let f = fmts.get(tz);
    if (!f) {
      f = new Intl.DateTimeFormat('en-US', {
        timeZone: tz || 'UTC',
        weekday: 'short',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      });
      fmts.set(tz, f);
    }
    return f.format(d);
  }

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    saving = true;
    error = null;
    try {
      const updated = await api.updateServer(serverId, {
        name: form.name.trim(),
        label: form.label.trim(),
        base_server_type: form.base_server_type.trim(),
        top_server_type: form.top_server_type.trim(),
        fallback_chain: form.fallback_chain,
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

  <!--
    First-window banner — only renders when the operator has picked
    `scheduled` mode. Heads up about the gap between "scheduled is
    selected" and "a window is configured and live":

      - no windows            → warning: a scheduled server with no
                                windows never rescales; shortcut to the
                                Windows tab.
      - currently in a window → info: state change is happening; no edit
                                needed.
      - next window in future → info: this is when the first automatic
                                rescale will fire.

    The shortcut is always present (even when no windows exist) so
    operators can reach the Windows tab in one click from any scheduled
    edit.
  -->
  {#if form.mode === 'scheduled'}
    <Alert
      variant={firstWindow.kind === 'none' ? 'warning' : 'default'}
      class="mb-6"
    >
      <div class="flex flex-wrap items-start justify-between gap-3">
        <div class="min-w-0 flex-1 space-y-1">
          {#if firstWindow.kind === 'none'}
            <p class="text-sm text-foreground">
              {m.server_edit_scheduled_no_windows()}
            </p>
          {:else if firstWindow.kind === 'in_window'}
            <p class="text-sm text-foreground">
              {m.server_edit_scheduled_in_window({
                target: firstWindow.target,
                until: fmtWhen(firstWindow.endsAt, form.timezone)
              })}
            </p>
          {:else}
            <p class="text-sm text-foreground">
              {m.server_edit_scheduled_next_window({
                when: fmtWhen(firstWindow.startsAt, form.timezone),
                target: firstWindow.target
              })}
            </p>
          {/if}
        </div>
        <Button variant="default" size="sm" href="/servers/{serverId}/windows">
          <Calendar class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
          {m.server_edit_scheduled_edit_windows()}
        </Button>
      </div>
    </Alert>
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
      /api/server-types) with role chips. The fallback chain is a
      drag-drop ServerTypeMultiSelect — chips can be reordered by
      pointer or touch. Base and top are excluded from the add list
      because a fallback that's the same as base or top is a config
      contradiction.
    -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <div class="flex flex-col gap-1.5">
        <Label for="f-base">{m.server_edit_field_base()}</Label>
        <ServerTypeSelect id="f-base" bind:value={form.base_server_type} server={server} location={server?.location} required />
      </div>
      <div class="flex flex-col gap-1.5">
        <Label for="f-top">{m.server_edit_field_top()}</Label>
        <ServerTypeSelect id="f-top" bind:value={form.top_server_type} server={server} location={server?.location} required />
      </div>
    </div>

    <div class="flex flex-col gap-1.5">
      <Label for="f-chain">{m.server_edit_field_fallback()}</Label>
      <ServerTypeMultiSelect
        id="f-chain"
        bind:value={form.fallback_chain}
        excluded={[form.base_server_type, form.top_server_type].filter(Boolean)}
        server={server}
        location={server?.location}
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