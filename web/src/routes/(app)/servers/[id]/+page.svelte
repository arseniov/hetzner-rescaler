<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { page } from '$app/stores';
  import { ArrowUp, ArrowDown, ChevronsUp, ChevronsDown, Pencil, Calendar } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server, RescaleEvent, Window_ as Window } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import Button from '$lib/components/ui/button.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import Tabs from '$lib/components/ui/tabs.svelte';
  import Dialog from '$lib/components/ui/dialog.svelte';
  import StatusBadge, { type Status } from '$lib/components/StatusBadge.svelte';
  import PendingRescaleBadge from '$lib/components/PendingRescaleBadge.svelte';
  import ModePill from '$lib/components/ModePill.svelte';
  import { pendingRescale } from '$lib/stores/pendingRescale.svelte';

  let server = $state<Server | null>(null);
  let events = $state<RescaleEvent[]>([]);
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);
  let busy = $state<'up' | 'down' | 'promote' | 'demote' | null>(null);
  let activeTab = $state<'overview' | 'windows' | 'events'>('overview');

  // Rescale-confirmation modal. Rescaling a server costs money (up)
  // or can cause outages (down), so the action button only stages the
  // confirmation — the actual API call fires from the modal's
  // "Rescale" button. The direction is held in a separate $state so
  // the dialog can render the right copy and we never have to peek
  // at the DOM to know which path we're on.
  let confirmOpen = $state(false);
  let confirmDirection = $state<'up' | 'down'>('up');

  function askRescale(direction: 'up' | 'down') {
    confirmDirection = direction;
    confirmOpen = true;
  }

  let serverId = $derived(Number($page.params.id));

  async function refresh() {
    loading = true;
    error = null;
    try {
      server = await api.getServer(serverId);
      events = await api.serverEvents(serverId, 25);
      windows = await api.listWindows(serverId);
      // Seed the pending store from the embedded event so the badge
      // appears even if the SSE stream was missed during page load
      // (e.g. rescale started before the operator opened the page).
      pendingRescale.setFromServer(server?.pending_event);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  // Background poll so the live status + current_type badges track a
  // server through shutdown / poweron cycles without requiring a
  // manual reload. The API fetches fresh state from Hetzner per call
  // (no server-side caching); the previous behaviour only re-fetched
  // on mount, so a server toggled while you had the page open stayed
  // visually "OFF" forever.
  //
  // 30s is a compromise between freshness and Hetzner API quota —
  // operators on the dashboard typically aren't staring at a single
  // server, so a half-minute cadence is responsive without hammering
  // the upstream. The interval is paused while a rescale is in flight
  // — `refresh()` already runs on every phase via the SSE-fed
  // pendingRescale store, so polling would be redundant churn.
  const LIVE_POLL_MS = 30_000;
  let pollTimer: ReturnType<typeof setInterval> | null = null;
  function startPolling() {
    stopPolling();
    pollTimer = setInterval(() => {
      if (pendingEvent === undefined) refresh().catch(() => {});
    }, LIVE_POLL_MS);
  }
  function stopPolling() {
    if (pollTimer !== null) {
      clearInterval(pollTimer);
      pollTimer = null;
    }
  }
  // React to rescale-state transitions: pause polling while a rescale
  // is active (events drive updates), resume once it clears so the
  // final running/off state settles in.
  $effect(() => {
    if (pendingEvent === undefined) startPolling();
    else stopPolling();
  });

  onDestroy(stopPolling);

  async function commitRescale() {
    busy = confirmDirection;
    confirmOpen = false;
    error = null;
    try {
      await api.rescale(serverId, { direction: confirmDirection, confirm: true });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      busy = null;
    }
  }

  async function promote() {
    busy = 'promote';
    error = null;
    try {
      await api.promote(serverId, { confirm: true });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      busy = null;
    }
  }

  async function demote() {
    busy = 'demote';
    error = null;
    try {
      await api.demote(serverId, { confirm: true });
      await refresh();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      busy = null;
    }
  }

  function fmtDays(mask: number): string {
    const labels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    const on: string[] = [];
    for (let i = 0; i < 7; i++) {
      if (mask & (1 << i)) on.push(labels[i]);
    }
    return on.length === 7 ? m.server_detail_window_every_day() : on.join(', ') || '—';
  }

  // Live state, mirrored from ServerCard so the server page can render
  // it without re-deriving twice. `current_type` is the Hetzner
  // authoritative view (might be null on a freshly-imported server
  // before the engine has polled it); we fall back to top_server_type
  // when null. `status` is the Hetzner lifecycle state; we fall back
  // to 'unknown' when missing so the badge always has a label.
  let currentType = $derived(server?.current_type ?? server?.top_server_type ?? null);
  let liveStatus = $derived<Status>((server?.status as Status | undefined) ?? 'unknown');

  // Compact last-events row: click any row (or the header's "View all"
  // link) to jump to the Events tab. The Overview tab is now a single
  // screen of facts + actions + activity, with detail one click away.
  function showEvents() {
    activeTab = 'events';
  }

  // Live pending-rescale state. SSE pushes these into `pendingRescale`;
  // the derived `isRescaling` drives the badge visibility and disables
  // every action button while a rescale is in flight so the operator
  // can't stage two rescales against the same server.
  let pendingEvent = $derived(serverId ? pendingRescale.get(serverId) : undefined);
  let isRescaling = $derived(pendingEvent !== undefined);
</script>

<svelte:head>
  <title>{server?.name ?? m.server_detail_tab_overview()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Server detail — three tabs: Overview (server facts, action buttons,
  current state, last events), Windows (scheduled rescales), Events
  (live log). The tabs are kept shallow so the operator can land on
  the right context in one click.
-->
<div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
  {#if error && !server}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  {#if loading}
    <p class="text-sm text-muted-foreground">{m.server_detail_loading()}</p>
  {:else if !server}
    <Alert variant="destructive">{m.server_detail_not_found()}</Alert>
  {:else}
    {@const s = server}

    <!-- Header. Title + metadata in a single row. Edit link sits on
         the right at md+ so the page doesn't grow vertical chrome on
         mobile. -->
    <header class="mb-6 flex flex-wrap items-end justify-between gap-3">
      <div class="min-w-0">
        <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
          {s.name}
        </h1>
        {#if pendingEvent}
          <div class="mt-2">
            <PendingRescaleBadge event={pendingEvent} />
          </div>
        {/if}
        <div class="mt-2">
          <ModePill
            mode={s.mode}
            promoteState={s.promote_state ?? null}
            lastTickAt={(events.find((e) => e.kind === 'scheduler_tick')?.started_at) ?? null}
            timezone={s.timezone}
            windows={windows.map((w) => ({
              days_of_week: w.days_of_week,
              start_time: w.start_time,
              stop_time: w.stop_time,
              target_type: w.target_type,
              enabled: w.enabled,
            }))}
          />
        </div>
      </div>
      <Button variant="ghost" size="sm" href="/servers/{s.id}/edit" disabled={isRescaling}>
        <Pencil class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
        {m.server_detail_edit()}
      </Button>
    </header>

    {#if error}
      <Alert variant="destructive" class="mb-6">{error}</Alert>
    {/if}

    <Tabs
      bind:value={activeTab}
      tabs={[
        { value: 'overview', label: m.server_detail_tab_overview() },
        { value: 'windows', label: m.server_detail_tab_windows() },
        { value: 'events', label: m.server_detail_tab_events() }
      ]}
    >
      {#snippet children(tab)}
        {#if tab === 'overview'}
          <!--
            Overview: three stacked panels —
              1. Current state (live Hetzner snapshot)
              2. Definition list + action buttons
              3. Compact recent events list (clicking jumps to Events)

            The action buttons are MODE-AWARE: each mode exposes the
            affordance that actually does something in that mode.
              - manual:       direct rescale up/down to configured base/top.
              - auto_promote: request promote/demote; the engine moves the
                              server at the next tick. Rescale up/down
                              are NOT a valid affordance under auto-promote
                              (they wouldn't carry the auto-promote
                              semantics), so they're hidden rather than
                              disabled-with-a-tooltip.
              - scheduled:    windows drive rescales. There is no
                              on-demand affordance here — edit windows
                              to change behaviour.
            This avoids "disabled buttons with tooltips" which is
            visual noise and forces operators to read why each button
            is grey before they can place the cursor.
          -->

          <!-- Current state panel. Shows the live Hetzner-reported
               status and current type. Falls back to event-derived or
               configured values when the live fields are null so the
               panel never blanks. -->
          <section
            aria-label={m.server_detail_current_state()}
            class="mb-4 rounded-md border border-border bg-card p-4"
          >
            <h2 class="mb-3 font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
              {m.server_detail_current_state()}
            </h2>
            <dl class="grid grid-cols-1 gap-x-6 gap-y-3 text-sm sm:grid-cols-2">
              <div>
                <dt class="text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_current_status()}
                </dt>
                <dd class="mt-1.5">
                  <StatusBadge status={liveStatus} />
                </dd>
              </div>
              <div>
                <dt class="text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_current_type()}
                </dt>
                <dd class="mt-1.5 font-mono text-foreground">
                  {currentType ?? '—'}
                </dd>
              </div>
            </dl>
          </section>

          <section class="mb-4 rounded-md border border-border bg-card p-4">
            <dl class="grid grid-cols-1 gap-x-6 gap-y-3 text-sm sm:grid-cols-2">
              <div>
                <dt class="text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_base_type()}
                </dt>
                <dd class="mt-0.5 font-mono text-foreground">{s.base_server_type}</dd>
              </div>
              <div>
                <dt class="text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_top_type()}
                </dt>
                <dd class="mt-0.5 font-mono text-foreground">{s.top_server_type}</dd>
              </div>
              <div class="sm:col-span-2">
                <dt class="text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_fallback_chain()}
                </dt>
                <dd class="mt-0.5 font-mono text-foreground">
                  {s.fallback_chain.join(' → ')}
                </dd>
              </div>
              <div>
                <dt class="text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_timezone()}
                </dt>
                <dd class="mt-0.5 font-mono text-foreground">{s.timezone}</dd>
              </div>
            </dl>

            <div class="mt-4 flex flex-wrap items-center gap-2 border-t border-border pt-4">
              {#if s.mode === 'manual'}
                <Button
                  variant="primary"
                  size="sm"
                  disabled={busy !== null || isRescaling}
                  onclick={() => askRescale('up')}
                >
                  <ArrowUp class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                  {busy === 'up' ? '…' : m.server_detail_rescale_up()}
                </Button>
                <Button
                  variant="default"
                  size="sm"
                  disabled={busy !== null || isRescaling}
                  onclick={() => askRescale('down')}
                >
                  <ArrowDown class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                  {busy === 'down' ? '…' : m.server_detail_rescale_down()}
                </Button>
              {:else if s.mode === 'auto_promote'}
                <Button
                  variant="primary"
                  size="sm"
                  disabled={busy !== null || isRescaling}
                  onclick={promote}
                >
                  <ChevronsUp class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                  {busy === 'promote' ? '…' : m.server_detail_promote()}
                </Button>
                <Button
                  variant="default"
                  size="sm"
                  disabled={busy !== null || isRescaling}
                  onclick={demote}
                >
                  <ChevronsDown class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                  {busy === 'demote' ? '…' : m.server_detail_demote()}
                </Button>
              {:else}
                <!--
                  Scheduled mode: windows drive rescales. There is no
                  on-demand rescale affordance here, so the buttons row
                  shows a quiet caption and the only action — Edit
                  windows — stays visible but icon-only.
                -->
                <span class="font-mono text-xs uppercase tracking-wider text-muted-foreground">
                  {m.server_detail_scheduled_caption()}
                </span>
              {/if}
              <!-- Icon-only "Edit windows" button. The label lives in
                   aria-label so screen readers still get the full
                   context; the visible glyph alone is enough for
                   operators who already know the affordance. -->
              <Button
                variant="default"
                size="icon"
                href="/servers/{s.id}/windows"
                aria-label={m.server_detail_edit_windows()}
                disabled={isRescaling}
              >
                <Calendar class="size-4" strokeWidth={1.75} aria-hidden="true" />
              </Button>
            </div>
          </section>

          <!-- Compact recent events. Five rows is enough to glance at
               "anything blowing up lately?" without scrolling. The
               entire panel header is clickable and switches to the
               Events tab, where the full log lives. -->
          <section class="rounded-md border border-border bg-card p-4">
            <header class="mb-3 flex items-center justify-between gap-3">
              <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                {m.server_detail_recent_events_compact()}
              </h2>
              <button
                type="button"
                onclick={showEvents}
                class="font-mono text-xs uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
              >
                {m.server_detail_view_all_events()} →
              </button>
            </header>
            {#if events.length === 0}
              <p class="text-sm text-muted-foreground">{m.server_detail_events_empty()}</p>
            {:else}
              <ul class="divide-y divide-border rounded-md border border-border">
                {#each events.slice(0, 5) as e (e.id)}
                  <li class="flex items-center justify-between px-3 py-2 text-sm">
                    <button
                      type="button"
                      onclick={showEvents}
                      class="min-w-0 flex-1 truncate text-left {e.ok ? 'text-foreground' : 'text-destructive'} hover:underline"
                    >
                      {e.kind}
                      {#if e.from_type && e.to_type}
                        ({e.from_type} → {e.to_type})
                      {/if}
                    </button>
                    <span class="ml-3 font-mono text-xs text-muted-foreground">
                      {new Date(e.started_at).toLocaleString()}
                    </span>
                  </li>
                {/each}
              </ul>
            {/if}
          </section>
        {:else if tab === 'windows'}
          <section class="rounded-md border border-border bg-card">
            <header class="flex items-center justify-between border-b border-border px-4 py-3">
              <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                {m.server_detail_windows_count({ count: windows.length })}
              </h2>
              <Button variant="ghost" size="sm" href="/servers/{s.id}/windows" disabled={isRescaling}>
                <Pencil class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                {m.server_detail_edit()}
              </Button>
            </header>
            {#if windows.length === 0}
              <p class="px-4 py-6 text-sm text-muted-foreground">
                {m.server_detail_windows_empty()}
              </p>
            {:else}
              <ul>
                {#each windows as w, i (w.id)}
                  <li
                    class="flex flex-wrap items-center justify-between gap-2 px-4 py-3 text-sm {i > 0 ? 'border-t border-border' : ''}"
                  >
                    <div class="min-w-0">
                      <div class="font-medium text-foreground">{w.label}</div>
                      <div class="mt-0.5 font-mono text-xs text-muted-foreground">
                        {w.start_time}–{w.stop_time} <span class="mx-1 text-foreground/30">→</span> {w.target_type}
                      </div>
                      <div class="mt-0.5 font-mono text-[11px] text-muted-foreground">
                        {fmtDays(w.days_of_week)}
                      </div>
                    </div>
                    <span
                      class="inline-flex items-center gap-1.5 font-mono text-xs uppercase tracking-wider {w.enabled ? 'text-success' : 'text-muted-foreground'}"
                    >
                      <span
                        class="inline-block size-1.5 rounded-full {w.enabled
                          ? 'bg-success'
                          : 'bg-muted-foreground/40'}"
                      ></span>
                      {w.enabled ? m.server_detail_window_enabled() : m.server_detail_window_disabled()}
                    </span>
                  </li>
                {/each}
              </ul>
            {/if}
          </section>
        {:else}
          <section class="rounded-md border border-border bg-card p-4">
            <h2 class="mb-3 font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
              {m.server_detail_recent_events()}
            </h2>
            <EventList {events} limit={20} />
          </section>
        {/if}
      {/snippet}
    </Tabs>
  {/if}
</div>

<!--
  Rescale-confirmation modal. Rescaling a server costs money (up) or
  can cause outages (down); an accidental click would be expensive.
  The modal title and description flip based on the staged direction
  so the operator sees the exact consequence they're about to commit.
-->
<Dialog
  bind:open={confirmOpen}
  title={confirmDirection === 'up'
    ? m.server_detail_rescale_up_confirm_title({ name: server?.name ?? '' })
    : m.server_detail_rescale_down_confirm_title({ name: server?.name ?? '' })}
  description={confirmDirection === 'up'
    ? m.server_detail_rescale_up_confirm_description()
    : m.server_detail_rescale_down_confirm_description()}
  size="md"
>
  <p class="text-sm text-muted-foreground">
    {#if confirmDirection === 'up'}
      <span class="font-mono">{server?.top_server_type ?? '—'}</span>
    {:else}
      <span class="font-mono">{server?.base_server_type ?? '—'}</span>
    {/if}
  </p>

  {#snippet footer()}
    <Button variant="ghost" onclick={() => (confirmOpen = false)} disabled={busy !== null}>
      {m.server_detail_cancel()}
    </Button>
    <Button variant="primary" onclick={commitRescale} disabled={busy !== null}>
      {m.server_detail_confirm_action()}
    </Button>
  {/snippet}
</Dialog>