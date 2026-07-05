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

  function options(rs: HoursAtTypeRow[]) {
    return {
      chart: { type: 'donut', height: 240, background: 'transparent' },
      series: rs.map((r) => r.costEur),
      labels: rs.map((r) => r.serverName),
      legend: { position: 'bottom' },
      dataLabels: { enabled: true },
      theme: { mode: isDark ? 'dark' : 'light' }
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
<p class="text-sm text-gray-600 dark:text-gray-400">Total: €{total.toFixed(2)}</p>