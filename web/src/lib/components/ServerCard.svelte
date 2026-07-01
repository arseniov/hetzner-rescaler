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
  class="block rounded-md border border-border p-4 hover:bg-muted/50 transition-colors"
>
  <div class="flex items-center justify-between">
    <span class="font-medium">{server.name}</span>
    <span class="text-xs text-muted-foreground">{modeLabel[server.mode]}</span>
  </div>
  <div class="mt-1 text-xs text-muted-foreground">
    {server.base_server_type} ⇄ {server.top_server_type}
  </div>
</a>
