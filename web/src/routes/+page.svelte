<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api';
  import { m } from '$lib/paraglide/messages.js';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { MetricsResponse } from '$lib/types';
  import KpiCard from '$lib/components/KpiCard.svelte';
  import RescalingActivityChart from '$lib/components/RescalingActivityChart.svelte';
  import CostBreakdownChart from '$lib/components/CostBreakdownChart.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let metrics = $state<MetricsResponse | null>(null);
  let metricsLoaded = $state(false);
  let error = $state<string | null>(null);
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
      // Seed the SSE store with a recent snapshot so /events and other
      // consumers see the same context. We don't keep the projects or
      // servers list here — each page owns its own fetch — so a
      // dashboard visit is just one API call + the events seed.
      const e = await api.globalEvents({ limit: 20 });
      eventsStream.replaceAll(e);
      await refreshMetrics();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
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
  Dashboard. Compact overview: the four KPIs each deep-link to the page
  that owns the underlying data, and the charts summarize the recent
  activity. Detail (per-project, per-server, full event log) lives on
  /projects, /servers and /events — the dashboard does not duplicate
  those lists.
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

  <!-- KPI row: four flat panels. The first three are deep-links into
       the page that owns the underlying list — the dashboard never
       duplicates those lists. The last-error card stays a passive
       indicator: clicking it would be ambiguous (which server?). -->
  <section
    aria-label="Key metrics"
    class="mb-6 grid grid-cols-2 gap-3 sm:grid-cols-4"
  >
    <a
      href="/servers"
      class="block rounded-md transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
      aria-label="Open servers"
    >
      <KpiCard
        label={m.kpi_active_servers()}
        value={formatRangeCount(metrics?.kpis.activeServerCount)}
        hint={m.kpi_active_servers_hint()}
        loading={!metricsLoaded}
      />
    </a>
    <a
      href="/projects"
      class="block rounded-md transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
      aria-label="Open projects"
    >
      <KpiCard
        label={m.kpi_projects()}
        value={formatRangeCount(metrics?.kpis.projectsWithTokenCount)}
        hint={m.kpi_projects_hint()}
        loading={!metricsLoaded}
      />
    </a>
    <a
      href="/events"
      class="block rounded-md transition-colors hover:bg-muted focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background"
      aria-label="Open events"
    >
      <KpiCard
        label={m.kpi_rescales_24h_ok()}
        value={formatRangeCount(metrics?.kpis.rescales24hOk)}
        loading={!metricsLoaded}
      />
    </a>
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
      class="grid grid-cols-1 gap-6 lg:grid-cols-3"
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
                class="rounded-sm px-2.5 py-1 font-mono uppercase tracking-wider transition-colors {chartRange === opt.v ? 'bg-card text-foreground' : 'text-muted-foreground hover:text-foreground'}"
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
</div>