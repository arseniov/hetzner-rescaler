<script lang="ts">
  import { onMount } from 'svelte';
  import { m } from '$lib/paraglide/messages.js';
  import type { RescaleEvent } from '$lib/types';
  import { cn } from '$lib/utils';

  type Props = { event: RescaleEvent; class?: string };
  let { event, class: className = '' }: Props = $props();

  // Re-tick every second so the elapsed counter advances without the
  // SSE stream pushing a fresh event. Pure 1Hz timer; we don't depend
  // on browser animation frames.
  let now = $state(Date.now());
  onMount(() => {
    const id = setInterval(() => { now = Date.now(); }, 1000);
    return () => clearInterval(id);
  });

  // Phase → human label. The four server-side phases are mapped to
  // their paraglide strings; any other value falls back to "Working…"
  // so the badge is never empty even mid-transition.
  let phaseLabel = $derived.by<string>(() => {
    switch (event.phase) {
      case 'shutting_down': return m.pending_rescale_phase_shutting_down();
      case 'changing_type': return m.pending_rescale_phase_changing_type();
      case 'powering_on':   return m.pending_rescale_phase_powering_on();
      case 'done':          return m.pending_rescale_phase_done();
      default:              return m.pending_rescale_unknown_phase();
    }
  });

  let elapsedSeconds = $derived(
    Math.max(0, Math.floor((now - new Date(event.started_at).getTime()) / 1000))
  );

  let elapsedLabel = $derived.by<string>(() => {
    if (elapsedSeconds < 60) return m.pending_rescale_elapsed_seconds({ seconds: elapsedSeconds });
    const m1 = Math.floor(elapsedSeconds / 60);
    const s = elapsedSeconds % 60;
    return m.pending_rescale_elapsed_minutes({ minutes: m1, seconds: s });
  });
</script>

<!--
  Amber accent pill — warm amber reserved for state-only moments.
  The pill lives inline with the page header; the tooltip on the
  wrapping <span> shows the full event details for operators who
  hover.
-->
<span
  class={cn(
    'inline-flex items-center gap-2 rounded-md border border-amber-500/30 bg-amber-500/10 px-2 py-1 font-mono text-[10px] uppercase tracking-wider text-amber-700 dark:text-amber-300',
    className
  )}
  title={m.pending_rescale_unknown_phase() + ' · ' + elapsedLabel}
  data-testid="pending-rescale-badge"
  data-phase={event.phase ?? 'unknown'}
>
  <span class="inline-block size-1.5 animate-pulse rounded-full bg-amber-500"></span>
  <span>{phaseLabel}</span>
  <span class="text-amber-700/60 dark:text-amber-300/60">·</span>
  <span>{elapsedLabel}</span>
</span>
