<script lang="ts">
  import type { RescaleEvent } from '$lib/types';
  interface Props { events: RescaleEvent[]; limit?: number }
  let { events, limit = 10 }: Props = $props();
  let shown = $derived(events.slice(0, limit));

  function fmtDate(iso: string) {
    return new Date(iso).toLocaleString();
  }
</script>

<ul class="divide-y divide-border rounded-md border border-border">
  {#each shown as e (e.id)}
    <li class="flex items-center justify-between px-3 py-2 text-sm">
      <span class={e.ok ? 'text-foreground' : 'text-destructive'}>
        {e.kind}
        {#if e.from_type && e.to_type}
          ({e.from_type} → {e.to_type})
        {/if}
      </span>
      <span class="text-muted-foreground">{fmtDate(e.started_at)}</span>
    </li>
  {:else}
    <li class="px-3 py-2 text-sm text-muted-foreground">No events yet.</li>
  {/each}
</ul>
