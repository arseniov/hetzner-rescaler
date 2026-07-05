<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Alert } from 'flowbite-svelte';
  import { api } from '$lib/api';
  import { m } from '$lib/paraglide/messages.js';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { Project, Server, RescaleEvent, MetricsResponse } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import ServerCard from '$lib/components/ServerCard.svelte';
  import KpiCard from '$lib/components/KpiCard.svelte';
  import RescalingActivityChart from '$lib/components/RescalingActivityChart.svelte';
  import CostBreakdownChart from '$lib/components/CostBreakdownChart.svelte';

  let projects = $state<Project[]>([]);
  let servers = $state<Server[]>([]);
  // Live events stream (REST seed + SSE updates). New events from the
  // SSE store are prepended to this array automatically.
  let events = $derived(eventsStream.events);
  let metrics = $state<MetricsResponse | null>(null);
  let error = $state<string | null>(null);
  let loading = $state(true);
  let chartRange = $state<'1d' | '7d' | '30d'>('7d');

  async function refreshMetrics() {
    try {
      metrics = await api.metrics(chartRange);
    } catch (err) {
      // Non-fatal: KPI cards render placeholder values if metrics
      // cannot be loaded. Surface the error in the console for now.
      console.warn('metrics refresh failed:', err);
    }
  }

  onMount(async () => {
    try {
      const [p, s, e] = await Promise.all([
        api.listProjects(),
        api.listServers(),
        api.globalEvents({ limit: 20 })
      ]);
      projects = p;
      servers = s;
      // Seed the SSE store with the initial REST snapshot so other pages
      // (e.g. /events) that read from `eventsStream.events` see the same
      // context.
      eventsStream.replaceAll(e);
      // Fetch metrics independently so a slow metrics endpoint does not
      // block the rest of the dashboard from rendering.
      await refreshMetrics();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  });
</script>

<div class="p-6 space-y-6 max-w-6xl mx-auto">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{m.dashboard_title()}</h1>

  {#if error}
    <Alert color="red">{error}</Alert>
  {/if}

  <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
    <KpiCard
      label={m.kpi_active_servers()}
      value={metrics?.kpis.activeServerCount ?? m.kpi_loading()}
      hint={m.kpi_active_servers_hint()}
    />
    <KpiCard
      label={m.kpi_projects()}
      value={metrics?.kpis.projectsWithTokenCount ?? m.kpi_loading()}
      hint={m.kpi_projects_hint()}
    />
    <KpiCard
      label={m.kpi_rescales_24h_ok()}
      value={metrics?.kpis.rescales24hOk ?? m.kpi_loading()}
    />
    <KpiCard
      label={m.kpi_last_error()}
      value={metrics?.kpis.lastRescaleError?.error ?? m.kpi_no_error()}
    />
  </div>

  {#if metrics}
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <Card class="lg:col-span-2">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white">{m.dashboard_chart_activity()}</h2>
          <select
            bind:value={chartRange}
            onchange={() => refreshMetrics()}
            class="rounded border-gray-300 dark:bg-gray-700 dark:text-white text-sm"
          >
            <option value="1d">{m.dashboard_chart_range_1d()}</option>
            <option value="7d">{m.dashboard_chart_range_7d()}</option>
            <option value="30d">{m.dashboard_chart_range_30d()}</option>
          </select>
        </div>
        <RescalingActivityChart data={metrics.rescaleCountsByDay ?? []} />
      </Card>

      <Card>
        <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
          {m.dashboard_chart_cost()}
        </h2>
        {#if (metrics.hoursAtType ?? []).length === 0}
          <p class="text-sm text-gray-600 dark:text-gray-400">{m.dashboard_chart_cost_empty()}</p>
        {:else}
          <CostBreakdownChart rows={metrics.hoursAtType ?? []} />
        {/if}
      </Card>
    </div>
  {/if}

  {#if loading}
    <p class="text-gray-600 dark:text-gray-400">{m.dashboard_loading()}</p>
  {:else}
    <Card>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.dashboard_section_projects({ count: projects.length })}
      </h2>
      {#if projects.length === 0}
        <p class="text-sm text-gray-600 dark:text-gray-400">No projects yet.</p>
      {:else}
        <ul class="space-y-2">
          {#each projects as p (p.id)}
            <li class="rounded-md border border-gray-200 dark:border-gray-700 p-3">
              <a href="/projects/{p.id}" class="font-medium hover:underline">{p.name}</a>
              <span class="ml-2 text-xs text-gray-600 dark:text-gray-400">
                {p.has_token ? m.projects_token_stored() : m.projects_no_token()}
              </span>
            </li>
          {/each}
        </ul>
      {/if}
    </Card>

    <Card>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.dashboard_section_servers({ count: servers.length })}
      </h2>
      {#if servers.length === 0}
        <p class="text-sm text-gray-600 dark:text-gray-400">{m.servers_empty()}</p>
      {:else}
        <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {#each servers.slice(0, 6) as s (s.id)}
            <ServerCard server={s} />
          {/each}
        </div>
      {/if}
    </Card>

    <Card>
      <h2 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
        {m.dashboard_section_recent_events()}
      </h2>
      <EventList {events} limit={10} />
    </Card>
  {/if}
</div>