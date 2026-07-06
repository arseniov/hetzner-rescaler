<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import type { RescaleCountsByDayRow } from '$lib/types';

  interface Props { data: RescaleCountsByDayRow[] }
  let { data }: Props = $props();

  let el: HTMLDivElement;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let chart: any = null;
  let isDark = $state(false);

  // ApexCharts accepts color strings but not CSS variables directly.
  // We read the relevant tokens off :root via getComputedStyle and pass
  // the resolved value so the chart honours dark / light switching.
  function token(name: string, fallback: string): string {
    if (typeof document === 'undefined') return fallback;
    const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
    return v || fallback;
  }

  function options(rows: RescaleCountsByDayRow[]) {
    return {
      chart: {
        type: 'area',
        height: 240,
        toolbar: { show: false },
        background: 'transparent',
        foreColor: token('--color-chart-axis', '#94a3b8'),
        fontFamily: 'JetBrains Mono, ui-monospace, monospace'
      },
      colors: [
        token('--color-chart-series-1', '#94a3b8'),
        token('--color-chart-series-2', '#f5a86b')
      ],
      series: [
        { name: 'OK', data: rows.map((r) => ({ x: r.date, y: r.ok })) },
        { name: 'Failed', data: rows.map((r) => ({ x: r.date, y: r.failed })) }
      ],
      xaxis: {
        type: 'category',
        labels: { style: { colors: token('--color-chart-axis', '#94a3b8') } },
        axisBorder: { color: token('--color-chart-grid', '#475569') },
        axisTicks: { color: token('--color-chart-grid', '#475569') }
      },
      yaxis: {
        labels: { style: { colors: token('--color-chart-axis', '#94a3b8') } }
      },
      grid: { borderColor: token('--color-chart-grid', '#475569'), strokeDashArray: 3 },
      stroke: { curve: 'smooth', width: 1.5 },
      dataLabels: { enabled: false },
      legend: {
        labels: { colors: token('--color-chart-axis', '#94a3b8') }
      },
      fill: {
        type: 'gradient',
        gradient: {
          shadeIntensity: 0.4,
          opacityFrom: 0.18,
          opacityTo: 0,
          stops: [0, 90, 100]
        }
      },
      tooltip: { theme: isDark ? 'dark' : 'light' }
    };
  }

  onMount(async () => {
    if (typeof document !== 'undefined') {
      isDark = document.documentElement.classList.contains('dark');
    }
    const mod = await import('apexcharts');
    const Ctor = (mod as any).default ?? (mod as any);
    chart = new Ctor(el, options(data));
    chart.render();
  });

  $effect(() => {
    if (chart) chart.updateOptions(options(data));
  });

  onDestroy(() => {
    chart?.destroy();
  });
</script>

<div bind:this={el}></div>
