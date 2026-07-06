<script lang="ts">
  import type { Server } from '$lib/types';
  interface Props { server: Server }
  let { server }: Props = $props();

  const modeLabel = {
    manual: 'Manual',
    auto_promote: 'Auto-promote',
    scheduled: 'Scheduled'
  } as const;
</script>

<a
  href="/servers/{server.id}"
  class="block rounded-md border border-border bg-card p-3 text-card-foreground transition-colors hover:bg-muted"
>
  <div class="flex items-center justify-between">
    <span class="font-medium text-foreground">{server.name}</span>
    <span class="font-mono text-xs text-muted-foreground">{modeLabel[server.mode]}</span>
  </div>
  <div class="mt-1 font-mono text-xs text-muted-foreground">
    {server.base_server_type} ⇄ {server.top_server_type}
  </div>
</a>
