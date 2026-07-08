<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { RefreshCw } from 'lucide-svelte';
  import { api } from '$lib/api';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import Button from '$lib/components/ui/button.svelte';

  type Bucket = 'ok' | 'warn' | 'fail' | 'unknown';

  let apiOk = $state<null | boolean>(null);
  let apiMs = $state<number | null>(null);
  let lastEventAt = $state<Date | null>(null);
  let recentErrors = $state(0);
  // Tick on every interval so $derived below can recompute from
  // Date.now(), which isn't reactive on its own.
  let nowMs = $state(Date.now());
  let eventAgeSec = $derived(
    lastEventAt ? Math.floor((nowMs - lastEventAt.getTime()) / 1000) : null
  );

  function ageBucket(sec: number | null): Bucket {
    if (sec == null) return 'unknown';
    if (sec > 86400) return 'fail';
    if (sec > 3600) return 'warn';
    return 'ok';
  }

  // Map a status bucket to a tailwind token class so the dot stays in
  // the same vocabulary as StatusBadge and the design system.
  const dotClass: Record<Bucket, string> = {
    ok: 'bg-success',
    warn: 'bg-warning',
    fail: 'bg-destructive',
    unknown: 'bg-muted-foreground'
  };

  async function refresh() {
    const t0 = Date.now();
    try {
      await api.healthz();
      apiOk = true;
      apiMs = Date.now() - t0;
    } catch {
      apiOk = false;
      apiMs = null;
    }

    if (eventsStream.events.length > 0) {
      const ev = eventsStream.events[0];
      lastEventAt = new Date(ev.started_at);
      recentErrors = eventsStream.events.filter((e) => !e.ok).length;
    }
  }

  onMount(refresh);
  // Tick once a second so the age label updates; the network call has
  // its own 10s interval below.
  const ageTimer = setInterval(() => {
    nowMs = Date.now();
  }, 1000);
  const timer = setInterval(refresh, 10000);
  onDestroy(() => {
    clearInterval(ageTimer);
    clearInterval(timer);
  });
</script>

<svelte:head>
  <title>{m.health_title()} · Hetzner Rescaler</title>
</svelte:head>

<!--
  System health — three flat panels summarising the rescaler's
  connectivity, recency, and recent errors. Each panel uses the same
  hairline border + monospaced figure vocabulary as the dashboard
  KPIs. Polling: the network check runs every 10s, the age label
  re-renders every 1s.
-->
<div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
  <header class="mb-6 flex items-end justify-between gap-3">
    <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
      {m.health_title()}
    </h1>
    <Button variant="ghost" size="sm" onclick={refresh}>
      <RefreshCw class="size-3.5" strokeWidth={1.75} aria-hidden="true" />
      {m.health_checking()}
    </Button>
  </header>

  <section aria-label="Health checks" class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <!-- API reachable -->
    <article class="rounded-md border border-border bg-card p-4">
      <p class="text-sm text-muted-foreground">{m.health_card_api_label()}</p>
      {#if apiOk === null}
        <div class="mt-2 flex items-center gap-2 text-sm text-muted-foreground">
          <span class="inline-block size-1.5 rounded-full bg-muted-foreground/40"></span>
          <span class="font-mono">…</span>
        </div>
      {:else if apiOk}
        <div class="mt-2 flex items-center gap-2">
          <span class="inline-block size-1.5 rounded-full bg-success"></span>
          <span class="font-mono text-2xl font-semibold tabular text-foreground">
            {apiMs}<span class="ml-0.5 text-sm font-normal text-muted-foreground">ms</span>
          </span>
        </div>
      {:else}
        <div class="mt-2 flex items-center gap-2">
          <span class="inline-block size-1.5 rounded-full bg-destructive"></span>
          <span class="text-sm text-foreground">{m.health_fail_above()}</span>
        </div>
      {/if}
    </article>

    <!-- Last event age -->
    <article class="rounded-md border border-border bg-card p-4">
      <p class="text-sm text-muted-foreground">{m.health_card_last_event_label()}</p>
      {#if eventAgeSec === null}
        <div class="mt-2 flex items-center gap-2 text-sm text-muted-foreground">
          <span class="inline-block size-1.5 rounded-full bg-muted-foreground/40"></span>
          <span class="font-mono">…</span>
        </div>
      {:else}
        {@const bucket = ageBucket(eventAgeSec)}
        <div class="mt-2 flex items-center gap-2">
          <span class="inline-block size-1.5 rounded-full {dotClass[bucket]}"></span>
          <span class="font-mono text-2xl font-semibold tabular text-foreground">
            {eventAgeSec}<span class="ml-0.5 text-sm font-normal text-muted-foreground">s</span>
          </span>
        </div>
      {/if}
    </article>

    <!-- Recent errors -->
    <article class="rounded-md border border-border bg-card p-4">
      <p class="text-sm text-muted-foreground">{m.health_card_recent_errors_label()}</p>
      <div class="mt-2 flex items-center gap-2">
        <span
          class="inline-block size-1.5 rounded-full {recentErrors > 0
            ? 'bg-destructive'
            : 'bg-success'}"
        ></span>
        <span class="font-mono text-2xl font-semibold tabular text-foreground">
          {recentErrors}
        </span>
      </div>
    </article>
  </section>

  <!--
    Threshold legend in the same muted treatment as the dashboard's
    hint lines. Lower-case mono keeps it visually subordinate.
  -->
  <p class="mt-6 font-mono text-xs text-muted-foreground">
    {m.health_ok_below()} · {m.health_warn_above()} · {m.health_fail_above()}
  </p>
</div>