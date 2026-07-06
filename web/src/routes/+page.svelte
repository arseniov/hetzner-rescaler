<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api';
  import { m } from '$lib/paraglide/messages.js';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { Project, Server, MetricsResponse } from '$lib/types';
  import EventList from '$lib/components/EventList.svelte';
  import ServerCard from '$lib/components/ServerCard.svelte';
  import KpiCard from '$lib/components/KpiCard.svelte';
  import RescalingActivityChart from '$lib/components/RescalingActivityChart.svelte';
  import CostBreakdownChart from '$lib/components/CostBreakdownChart.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let projects = $state<Project[]>([]);
  let servers = $state<Server[]>([]);
  let events = $derived(eventsStream.events);
  let metrics = $state<MetricsResponse | null>(null);
  let metricsLoaded = $state(false);
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
    } finally {
      metricsLoaded = true;
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
      // Seed the SSE store with the initial REST snapshot so /events
      // and other consumers see the same context.
      eventsStream.replaceAll(e);
      await refreshMetrics();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      loading = false;
    }
  });

  function formatRangeCount(n: number | null | undefined): string {
    return n === null || n === undefined ? '—' : String(n);
  }
</script>

<svelte:head>
  <title>{m.dashboard_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  Dashboard. One job: tell the operator at a glance what the rescaler
  is doing. Pages have one page, one job, one heading — no hero
  sections, no decorative chrome. Density is the operator's right:
  when in doubt, show the data denser.
-->
<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6 flex items-end justify-between gap-3">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.dashboard_title()}
    </h1>
    {#if metricsLoaded}
      <span class="font-mono text-xs text-muted-foreground tabular">
        {new Date().toLocaleTimeString()}
      </span>
    {/if}
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <!-- KPI row: four flat panels in a responsive grid. The numeric
       values use tabular font so digits align across the row. -->
  <section
    aria-label="Key metrics"
    class="mb-6 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4"
  >
    <KpiCard
      label={m.kpi_active_servers()}
      value={formatRangeCount(metrics?.kpis.activeServerCount)}
      hint={m.kpi_active_servers_hint()}
      loading={!metricsLoaded}
    />
    <KpiCard
      label={m.kpi_projects()}
      value={formatRangeCount(metrics?.kpis.projectsWithTokenCount)}
      hint={m.kpi_projects_hint()}
      loading={!metricsLoaded}
    />
    <KpiCard
      label={m.kpi_rescales_24h_ok()}
      value={formatRangeCount(metrics?.kpis.rescales24hOk)}
      loading={!metricsLoaded}
    />
    <KpiCard
      label={m.kpi_last_error()}
      value={metrics?.kpis.lastRescaleError?.error ?? m.kpi_no_error()}
      loading={!metricsLoaded}
    />
  </section>

  <!-- Charts row. Two panels separated by a hairline; activity takes
       2/3 of the column, cost breakdown takes 1/3. The range picker is
       a compact text-button group, not a system <select>. -->
  {#if metrics}
    <section
      aria-label="Charts"
      class="mb-6 grid grid-cols-1 gap-6 lg:grid-cols-3"
    >
      <div class="rounded-md border border-border bg-card p-4 lg:col-span-2">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="font-display text-base font-semibold text-foreground">
            {m.dashboard_chart_activity()}
          </h2>
          <!-- Range picker as a segmented control. Three text buttons
               share a row; the active one inverts the foreground. -->
          <div role="radiogroup" aria-label={m.dashboard_chart_range()} class="inline-flex rounded-md border border-border bg-muted p-0.5 text-xs">
            {#each [{ v: '1d', l: m.dashboard_chart_range_1d() }, { v: '7d', l: m.dashboard_chart_range_7d() }, { v: '30d', l: m.dashboard_chart_range_30d() }] as opt}
              <button
                type="button"
                role="radio"
                aria-checked={chartRange === opt.v}
                onclick={() => { chartRange = opt.v as '1d' | '7d' | '30d'; refreshMetrics(); }}
                class="rounded-sm px-2.5 py-1 font-mono uppercase tracking-wider transition-colors {chartRange === opt.v ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'}"
              >
                {opt.l}
              </button>
            {/each}
          </div>
        </div>
        <RescalingActivityChart data={metrics.rescaleCountsByDay ?? []} />
      </div>

      <div class="rounded-md border border-border bg-card p-4">
        <h2 class="mb-3 font-display text-base font-semibold text-foreground">
          {m.dashboard_chart_cost()}
        </h2>
        {#if (metrics.hoursAtType ?? []).length === 0}
          <p class="text-sm text-muted-foreground">{m.dashboard_chart_cost_empty()}</p>
        {:else}
          <CostBreakdownChart rows={metrics.hoursAtType ?? []} />
        {/if}
      </div>
    </section>
  {/if}

  {#if loading}
    <p class="text-sm text-muted-foreground">{m.dashboard_loading()}</p>
  {:else}
    <!-- Lists: projects, servers, recent events. Each is a flat panel,
         no internal headings hierarchy, no nested cards. The list rows
         are separated by hairlines (use border-t on items past the
         first), not by cards-with-padding. -->
    <section
      aria-label="Project summary"
      class="mb-6 rounded-md border border-border bg-card"
    >
      <header class="flex items-center justify-between border-b border-border px-4 py-3">
        <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          {m.dashboard_section_projects({ count: projects.length })}
        </h2>
      </header>
      {#if projects.length === 0}
        <p class="px-4 py-4 text-sm text-muted-foreground">No projects yet.</p>
      {:else}
        <ul>
          {#each projects as p, i (p.id)}
            <li class="flex items-center justify-between px-4 py-2.5 text-sm {i > 0 ? 'border-t border-border' : ''}">
              <a href="/projects/{p.id}" class="font-medium text-foreground hover:underline">
                {p.name}
              </a>
              <span class="font-mono text-xs uppercase tracking-wider text-muted-foreground">
                {p.has_token ? m.projects_token_stored() : m.projects_no_token()}
              </span>
            </li>
          {/each}
        </ul>
      {/if}
    </section>

    <section
      aria-label="Server summary"
      class="mb-6 rounded-md border border-border bg-card"
    >
      <header class="flex items-center justify-between border-b border-border px-4 py-3">
        <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          {m.dashboard_section_servers({ count: servers.length })}
        </h2>
      </header>
      {#if servers.length === 0}
        <p class="px-4 py-4 text-sm text-muted-foreground">{m.servers_empty()}</p>
      {:else}
        <div class="grid grid-cols-1 gap-3 p-4 sm:grid-cols-2 lg:grid-cols-3">
          {#each servers.slice(0, 6) as s (s.id)}
            <ServerCard server={s} />
          {/each}
        </div>
      {/if}
    </section>

    <section
      aria-label="Recent events"
      class="rounded-md border border-border bg-card"
    >
      <header class="flex items-center justify-between border-b border-border px-4 py-3">
        <h2 class="font-display text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          {m.dashboard_section_recent_events()}
        </h2>
      </header>
      <EventList {events} limit={10} />
    </section>
  {/if}
</div>
