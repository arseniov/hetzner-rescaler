<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { api } from '$lib/api';
  import { m } from '$lib/paraglide/messages.js';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import type { MetricsResponse, Project, Server, RescaleEvent } from '$lib/types';
  import KpiCard from '$lib/components/KpiCard.svelte';
  import RescalingActivityChart from '$lib/components/RescalingActivityChart.svelte';
  import CostBreakdownChart from '$lib/components/CostBreakdownChart.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let metrics = $state<MetricsResponse | null>(null);
  let metricsLoaded = $state(false);
  let error = $state<string | null>(null);
  let chartRange = $state<'1d' | '7d' | '30d'>('7d');

  // Lookup tables for the "Last events" panel: project by id, server
  // by id (carries project_id), so a single event row can show both
  // server name and the parent project name. Populated in onMount —
  // before fetch completes the row falls back to placeholder strings.
  let projects = $state<Project[]>([]);
  let servers = $state<Server[]>([]);

  // Live clock — updates every second so the time chip in the header
  // is genuinely "now", not "the moment the page rendered". The previous
  // version captured `new Date()` once, so the displayed time froze as
  // soon as the page settled.
  let now = $state(new Date());
  let clockTimer: ReturnType<typeof setInterval> | null = null;

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
    clockTimer = setInterval(() => {
      now = new Date();
    }, 1000);
    try {
      // Seed the SSE store with a recent snapshot so /events and other
      // consumers see the same context. In parallel we fetch the
      // project + server lists so the "Last events" panel can render
      // human-readable names next to each event row without a second
      // round trip.
      const [e, p, s] = await Promise.all([
        api.globalEvents({ limit: 20 }),
        api.listProjects(),
        api.listServers()
      ]);
      eventsStream.replaceAll(e);
      projects = p;
      servers = s;
      await refreshMetrics();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  });

  onDestroy(() => {
    if (clockTimer) clearInterval(clockTimer);
  });

  function formatRangeCount(n: number | null | undefined): string {
    return n === null || n === undefined ? '—' : String(n);
  }

  // Two most recent events — same order as the SSE store (newest
  // first). Capped at two so the panel stays compact and never out-
  // grows the surrounding KPI cards in the grid.
  let lastEvents = $derived(eventsStream.events.slice(0, 2));

  // O(1) server_id → context lookup. Built fresh from the latest
  // projects/servers lists; if either list is empty the field shows
  // a neutral placeholder so the row still renders.
  let serverCtx = $derived.by(() => {
    const projById = new Map(projects.map((p) => [p.id, p]));
    const ctx = new Map<number, { serverName: string; projectName: string }>();
    for (const s of servers) {
      ctx.set(s.id, {
        serverName: s.name,
        projectName: projById.get(s.project_id)?.name ?? '—'
      });
    }
    return ctx;
  });

  function ctxFor(e: RescaleEvent) {
    return serverCtx.get(e.server_id) ?? { serverName: `#${e.server_id}`, projectName: '—' };
  }

  // Compact "12:34:56" timestamp — the panel is dense; the full
  // datetime lives on /events. We include seconds so the live time
  // visibly ticks even when seconds is the only thing that changes.
  function fmtTime(iso: string): string {
    const d = new Date(iso);
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
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
      <span class="font-mono text-xs tabular text-muted-foreground tabular-nums">
        {now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
      </span>
    {/if}
  </header>

  {#if error}
    <Alert variant="destructive" class="mb-6">{error}</Alert>
  {/if}

  <!-- KPI row: four flat panels. The first three are deep-links into
       the page that owns the underlying list — the dashboard never
       duplicates those lists. The fourth slot is a passive "Last
       events" panel: deep-linking it would be ambiguous (which
       server?), but each row inside links to its own server. -->
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
        value={formatRangeCount(metrics?.kpis.active_server_count)}
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
        value={formatRangeCount(metrics?.kpis.projects_with_token_count)}
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
        value={formatRangeCount(metrics?.kpis.rescales_24h_ok)}
        loading={!metricsLoaded}
      />
    </a>

    <!--
      Last events panel. Sits in the KPI grid in place of the
      one-off "Last rescale error" indicator so the operator sees the
      most recent activity at a glance. Three compact rows, each
      linking to the server that emitted the event. Each row carries:
      status dot (success / destructive), server name + kind, parent
      project name, and the HH:MM timestamp. Deeper detail (full
      datetime, error message, from/to type) lives on /events.
    -->
    <section
      aria-label={m.dashboard_last_events()}
      class="flex flex-col gap-1.5 rounded-md border border-border bg-card px-4 py-3"
    >
      <p class="text-sm text-muted-foreground">{m.dashboard_last_events()}</p>
      {#if lastEvents.length === 0}
        <p class="text-xs text-muted-foreground">{m.dashboard_no_events()}</p>
      {:else}
        <ul class="divide-y divide-border">
          {#each lastEvents as e (e.id)}
            {@const ctx = ctxFor(e)}
            <li class="flex items-center gap-2 py-1 text-sm">
              <span
                class="inline-block size-1.5 shrink-0 rounded-full {e.ok ? 'bg-success' : 'bg-destructive'}"
                aria-hidden="true"
              ></span>
              <a
                href="/servers/{e.server_id}"
                class="min-w-0 flex-1 truncate text-foreground hover:underline"
              >
                <span class="font-medium">{ctx.serverName}</span>
                <span class="mx-1 text-foreground/30">·</span>
                <span class="font-mono text-xs text-muted-foreground">{e.kind}</span>
              </a>
              <span class="hidden font-mono text-xs text-muted-foreground sm:inline">
                {ctx.projectName}
              </span>
              <span class="shrink-0 font-mono text-xs tabular text-muted-foreground">
                {fmtTime(e.started_at)}
              </span>
            </li>
          {/each}
        </ul>
      {/if}
    </section>
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
        <RescalingActivityChart data={metrics.rescale_counts_by_day ?? []} />
      </div>

      <div class="rounded-md border border-border bg-card p-4">
        <h2 class="mb-3 font-display text-base font-semibold text-foreground">
          {m.dashboard_chart_cost()}
        </h2>
        {#if (metrics.hours_at_type ?? []).length === 0}
          <p class="text-sm text-muted-foreground">{m.dashboard_chart_cost_empty()}</p>
        {:else}
          <CostBreakdownChart rows={metrics.hours_at_type ?? []} />
        {/if}
      </div>
    </section>
  {/if}
</div>