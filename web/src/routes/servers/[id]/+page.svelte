<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { ArrowUp, ArrowDown, ChevronsUp, ChevronsDown, Pencil, Calendar } from 'lucide-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server, RescaleEvent, Window_ as Window } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import Button from '$lib/components/ui/button.svelte';
  import Alert from '$lib/components/ui/alert.svelte';
  import Tabs from '$lib/components/ui/tabs.svelte';

  let server = $state<Server | null>(null);
  let events = $state<RescaleEvent[]>([]);
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);
  let busy = $state<'up' | 'down' | 'promote' | 'demote' | null>(null);
  let activeTab = $state<'overview' | 'windows' | 'events'>('overview');

  let serverId = $derived(Number($page.params.id));

  async function refresh() {
    loading = true;
    error = null;
    try {
      server = await api.getServer(serverId);
      events = await api.serverEvents(serverId, 25);
      windows = await api.listWindows(serverId);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  async function rescale(direction: 'up' | 'down') {
    busy = direction;
    error = null;
    try {
      await api.rescale(serverId, { direction, confirm: true });
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
</script>

<svelte:head>
  <title>{server?.name ?? m.server_detail_tab_overview()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Server detail — three tabs: Overview (server facts + actions),
  Windows (scheduled rescales), Events (live log). The tabs are
  kept shallow so the operator can land on the right context in one
  click.
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
    <!--
      Capture `server` into a non-null local alias. The tabs snippet
      is rendered as a separate closure so TypeScript's flow
      narrowing from the `{:else if !server}` guard doesn't survive
      across the snippet boundary; the `!` lets us tell the compiler
      "we've already checked above, give me a non-null view of this
      value inside the snippet".
    -->
    {@const s = server}

    <!-- Header. Title + metadata in a single row. Edit link sits on
         the right at md+ so the page doesn't grow vertical chrome on
         mobile. -->
    <header class="mb-6 flex flex-wrap items-end justify-between gap-3">
      <div class="min-w-0">
        <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
          {s.name}
        </h1>
        <p class="mt-1 font-mono text-xs text-muted-foreground">
          {m.server_detail_hcloud_id({ id: s.hcloud_server_id })}
          <span class="mx-2 text-foreground/30">·</span>
          {m.server_detail_mode({ mode: s.mode })}
          {#if s.promote_state}
            <span class="mx-2 text-foreground/30">·</span>
            {m.server_detail_state({ state: s.promote_state })}
          {/if}
        </p>
      </div>
      <Button variant="ghost" size="sm" href="/servers/{s.id}/edit">
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
            Overview: definition-list grid for the server facts, then
            a row of action buttons. The promote / demote buttons
            only render for auto_promote mode; for other modes the
            rescale pair is sufficient.
          -->
          <section class="rounded-md border border-border bg-card p-4">
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

            <div class="mt-4 flex flex-wrap gap-2 border-t border-border pt-4">
              <Button
                variant="primary"
                size="sm"
                disabled={busy !== null}
                onclick={() => rescale('up')}
              >
                <ArrowUp class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                {busy === 'up' ? '…' : m.server_detail_rescale_up()}
              </Button>
              <Button
                variant="default"
                size="sm"
                disabled={busy !== null}
                onclick={() => rescale('down')}
              >
                <ArrowDown class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                {busy === 'down' ? '…' : m.server_detail_rescale_down()}
              </Button>
              {#if s.mode === 'auto_promote'}
                <Button
                  variant="default"
                  size="sm"
                  disabled={busy !== null}
                  onclick={promote}
                >
                  <ChevronsUp class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                  {busy === 'promote' ? '…' : m.server_detail_promote()}
                </Button>
                <Button
                  variant="default"
                  size="sm"
                  disabled={busy !== null}
                  onclick={demote}
                >
                  <ChevronsDown class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                  {busy === 'demote' ? '…' : m.server_detail_demote()}
                </Button>
              {/if}
              <Button variant="default" size="sm" href="/servers/{s.id}/windows">
                <Calendar class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
                {m.server_detail_edit_windows()}
              </Button>
            </div>
          </section>
        {:else if tab === 'windows'}
          <section class="rounded-md border border-border bg-card">
            <header class="flex items-center justify-between border-b border-border px-4 py-3">
              <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                {m.server_detail_windows_count({ count: windows.length })}
              </h2>
              <Button variant="ghost" size="sm" href="/servers/{s.id}/windows">
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