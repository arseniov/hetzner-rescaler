<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { Tabs, TabItem, Card, Button, Alert, Badge } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api, ApiError } from '$lib/api';
  import type { Server, RescaleEvent, Window_ as Window } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';

  let server = $state<Server | null>(null);
  let events = $state<RescaleEvent[]>([]);
  let windows = $state<Window[]>([]);
  let error = $state<string | null>(null);
  let loading = $state(true);
  let busy = $state<'up' | 'down' | 'promote' | 'demote' | null>(null);

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

<div class="p-6 max-w-5xl mx-auto space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">
        {server?.name ?? '…'}
      </h1>
      {#if server}
        <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
          {m.server_detail_hcloud_id({ id: server.hcloud_server_id })}
          · {m.server_detail_mode({ mode: server.mode })}
          {#if server.promote_state}
            · {m.server_detail_state({ state: server.promote_state })}
          {/if}
        </p>
      {/if}
    </div>
    {#if server}
      <Button color="alternative" href="/servers/{server.id}/edit">
        {m.server_detail_edit()}
      </Button>
    {/if}
  </div>

  {#if error}<Alert color="danger">{error}</Alert>{/if}

  {#if loading}
    <p class="text-sm text-gray-600 dark:text-gray-400">{m.server_detail_loading()}</p>
  {:else if !server}
    <Alert color="danger">{m.server_detail_not_found()}</Alert>
  {:else}
    <Tabs tabStyle="underline">
      <TabItem title={m.server_detail_tab_overview()}>
        <Card>
          <dl class="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
            <div>
              <dt class="text-gray-600 dark:text-gray-400">{m.server_detail_base_type()}</dt>
              <dd class="font-medium text-gray-900 dark:text-white">{server.base_server_type}</dd>
            </div>
            <div>
              <dt class="text-gray-600 dark:text-gray-400">{m.server_detail_top_type()}</dt>
              <dd class="font-medium text-gray-900 dark:text-white">{server.top_server_type}</dd>
            </div>
            <div class="sm:col-span-2">
              <dt class="text-gray-600 dark:text-gray-400">{m.server_detail_fallback_chain()}</dt>
              <dd class="font-medium text-gray-900 dark:text-white">
                {server.fallback_chain.join(' → ')}
              </dd>
            </div>
            <div>
              <dt class="text-gray-600 dark:text-gray-400">{m.server_detail_timezone()}</dt>
              <dd class="font-medium text-gray-900 dark:text-white">{server.timezone}</dd>
            </div>
          </dl>

          <div class="mt-4 flex flex-wrap gap-2">
            <Button
              color="brand"
              disabled={busy !== null}
              onclick={() => rescale('up')}
            >
              {busy === 'up' ? '…' : m.server_detail_rescale_up()}
            </Button>
            <Button
              color="alternative"
              disabled={busy !== null}
              onclick={() => rescale('down')}
            >
              {busy === 'down' ? '…' : m.server_detail_rescale_down()}
            </Button>
            {#if server.mode === 'auto_promote'}
              <Button
                color="alternative"
                disabled={busy !== null}
                onclick={promote}
              >
                {busy === 'promote' ? '…' : m.server_detail_promote()}
              </Button>
              <Button
                color="alternative"
                disabled={busy !== null}
                onclick={demote}
              >
                {busy === 'demote' ? '…' : m.server_detail_demote()}
              </Button>
            {/if}
            <Button color="alternative" href="/servers/{server.id}/windows">
              {m.server_detail_edit_windows()}
            </Button>
          </div>
        </Card>
      </TabItem>

      <TabItem title={m.server_detail_tab_windows()}>
        <Card>
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-lg font-medium text-gray-900 dark:text-white">
              {m.server_detail_windows_count({ count: windows.length })}
            </h2>
            <Button color="alternative" href="/servers/{server.id}/windows">
              {m.server_detail_edit()}
            </Button>
          </div>
          {#if windows.length === 0}
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {m.server_detail_windows_empty()}
            </p>
          {:else}
            <ul class="divide-y divide-gray-200 dark:divide-gray-700">
              {#each windows as w (w.id)}
                <li class="flex items-center justify-between py-3 text-sm">
                  <div>
                    <div class="font-medium text-gray-900 dark:text-white">{w.label}</div>
                    <div class="text-gray-600 dark:text-gray-400">
                      {w.start_time}–{w.stop_time} → {w.target_type}
                    </div>
                    <div class="text-xs text-gray-500 dark:text-gray-500">
                      {fmtDays(w.days_of_week)}
                    </div>
                  </div>
                  <Badge color={w.enabled ? 'green' : 'gray'}>
                    {w.enabled ? m.server_detail_window_enabled() : m.server_detail_window_disabled()}
                  </Badge>
                </li>
              {/each}
            </ul>
          {/if}
        </Card>
      </TabItem>

      <TabItem title={m.server_detail_tab_events()}>
        <Card>
          <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
            {m.server_detail_recent_events()}
          </h2>
          <EventList {events} limit={20} />
        </Card>
      </TabItem>
    </Tabs>
  {/if}
</div>