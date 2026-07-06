<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import type { HoursAtTypeRow } from '$lib/types';

  interface Props { rows: HoursAtTypeRow[] }
  let { rows }: Props = $props();

  let el: HTMLDivElement;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let chart: any = null;
  let isDark = $state(false);

  const total = $derived(rows.reduce((acc, r) => acc + r.costEur, 0));

  function token(name: string, fallback: string): string {
    if (typeof document === 'undefined') return fallback;
    const v = getComputedStyle(document.documentElement).getPropertyValue(name).trim();
    return v || fallback;
  }

  function options(rs: HoursAtTypeRow[]) {
    return {
      chart: {
        type: 'donut',
        height: 240,
        background: 'transparent',
        foreColor: token('--color-chart-axis', '#94a3b8')
      },
      series: rs.map((r) => r.costEur),
      labels: rs.map((r) => r.serverName),
      colors: [token('--color-chart-series-2', '#f5a86b')],
      legend: {
        position: 'bottom',
        labels: { colors: token('--color-chart-axis', '#94a3b8') }
      },
      dataLabels: { enabled: false },
      stroke: { width: 1, colors: [token('--color-card', '#1f2937')] },
      plotOptions: {
        pie: {
          donut: {
            size: '70%',
            labels: {
              show: true,
              total: {
                show: true,
                label: 'Total',
                color: token('--color-chart-axis', '#94a3b8'),
                fontFamily: 'JetBrains Mono, ui-monospace, monospace'
              },
              value: {
                color: token('--color-foreground', '#e2e8f0'),
                fontFamily: 'JetBrains Mono, ui-monospace, monospace'
              }
            }
          }
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
    chart = new Ctor(el, options(rows));
    chart.render();
  });

  $effect(() => {
    if (chart) chart.updateOptions(options(rows));
  });

  onDestroy(() => {
    chart?.destroy();
  });
</script>

<div bind:this={el}></div>
<p class="mt-2 font-mono text-sm text-muted-foreground tabular">Total: €{total.toFixed(2)}</p>
