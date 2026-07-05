<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Card, Alert, Badge, Spinner, Button } from 'flowbite-svelte';
  import { api } from '$lib/api';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import { m } from '$lib/paraglide/messages.js';

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

  function badgeColor(bucket: Bucket): 'green' | 'yellow' | 'red' | 'gray' {
    switch (bucket) {
      case 'ok':
        return 'green';
      case 'warn':
        return 'yellow';
      case 'fail':
        return 'red';
      default:
        return 'gray';
    }
  }

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

<div class="p-6 max-w-5xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{m.health_title()}</h1>

  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
    <Card class="border-0">
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.health_card_api_label()}</p>
      {#if apiOk == null}
        <Spinner class="mt-2" />
      {:else if apiOk}
        <Alert color="green" class="mt-2">{m.health_ok_below()} ({apiMs}ms)</Alert>
      {:else}
        <Alert color="danger" class="mt-2">{m.health_fail_above()}</Alert>
      {/if}
    </Card>

    <Card class="border-0">
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.health_card_last_event_label()}</p>
      {#if eventAgeSec == null}
        <Spinner class="mt-2" />
      {:else}
        <Badge color={badgeColor(ageBucket(eventAgeSec))} class="mt-2">
          {eventAgeSec}s
        </Badge>
      {/if}
    </Card>

    <Card class="border-0">
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.health_card_recent_errors_label()}</p>
      <p class="mt-2 text-3xl font-semibold text-gray-900 dark:text-white">{recentErrors}</p>
    </Card>
  </div>

  <Button color="alternative" onclick={refresh}>{m.health_checking()}</Button>
</div>
