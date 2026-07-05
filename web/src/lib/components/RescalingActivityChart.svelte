<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import type { RescaleCountsByDayRow } from '$lib/types';

  interface Props { data: RescaleCountsByDayRow[] }
  let { data }: Props = $props();

  let el: HTMLDivElement;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  let chart: any = null;
  let isDark = $state(false);

  function options(rows: RescaleCountsByDayRow[]) {
    return {
      chart: { type: 'area', height: 240, toolbar: { show: false }, background: 'transparent' },
      series: [
        { name: 'OK', data: rows.map((r) => ({ x: r.date, y: r.ok })) },
        { name: 'Failed', data: rows.map((r) => ({ x: r.date, y: r.failed })) }
      ],
      xaxis: { type: 'category' },
      stroke: { curve: 'smooth' },
      dataLabels: { enabled: false },
      theme: { mode: isDark ? 'dark' : 'light' }
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