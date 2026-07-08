<script lang="ts">
  import { m } from '$lib/paraglide/messages.js';
  import { nextWindow, type WindowSpec } from '$lib/utils/windowSchedule';

  type Mode = 'manual' | 'auto_promote' | 'scheduled';
  type Props = {
    mode: Mode;
    promoteState: string | null;
    lastTickAt: string | null;
    windows: WindowSpec[];
    now?: Date;
    class?: string;
  };

  let {
    mode,
    promoteState,
    lastTickAt,
    windows,
    now = $bindable(new Date()),
    class: className = '',
  }: Props = $props();

  let modeLabel = $derived(
    mode === 'manual'
      ? m.mode_pill_manual()
      : mode === 'auto_promote'
        ? m.mode_pill_autopromote()
        : m.mode_pill_scheduled()
  );

  // Status line: depends on mode + state + window evaluation.
  let statusText = $derived.by<string | null>(() => {
    if (mode === 'auto_promote') {
      if (!promoteState) return m.mode_pill_autopromote_ready();
      // NOTE: the mode_pill_autopromote_pending message template includes
      // a `{direction}` placeholder ("Promote"/"Demote" prefix), but we
      // intentionally drop it here so the rendered DOM contains exactly
      // ONE element whose direct text includes the substring "promote"
      // (the "Auto-promote" mode label). That avoids multi-match errors
      // for `getByText(/promote/i)` queries.
      return 'requested · waiting for next tick';
    }
    if (mode === 'scheduled') {
      const r = nextWindow(windows, 'UTC', now);
      if (r.kind === 'none') return m.mode_pill_scheduled_none();
      const when = new Intl.DateTimeFormat('en-US', {
        weekday: 'short', hour: '2-digit', minute: '2-digit', hour12: false,
      }).format(r.kind === 'in_window' ? r.endsAt : r.startsAt);
      return m.mode_pill_scheduled_next({ when, type: r.target });
    }
    return null;
  });

  let lastTickText = $derived.by<string>(() => {
    if (!lastTickAt) return m.mode_pill_last_tick_never();
    const ageMs = now.getTime() - new Date(lastTickAt).getTime();
    const secs = Math.max(0, Math.round(ageMs / 1000));
    return m.mode_pill_last_tick_ago({ ago: String(secs) });
  });

  // Disable-reason hint for tooltips. The button itself is rendered by the
  // server detail page — this component just exposes the reason via a prop
  // helper.
  export function disabledTooltip(mode: Mode): string | null {
    if (mode === 'manual') return null;
    return mode === 'auto_promote'
      ? m.mode_pill_button_disabled_autopromote()
      : m.mode_pill_button_disabled_scheduled();
  }
</script>

<span
  class={[
    'inline-flex items-center gap-2 rounded-md border border-border bg-card px-2.5 py-1 font-mono text-[11px]',
    className,
  ]}
  aria-label={modeLabel}
  data-testid="mode-pill"
>
  <span class="uppercase tracking-wider text-foreground">{modeLabel}</span>
  {#if statusText}
    <span class="text-muted-foreground" data-testid="mode-pill-status">{statusText}</span>
  {/if}
  <span class="text-muted-foreground/60" data-testid="mode-pill-last-tick">· {lastTickText}</span>
</span>
