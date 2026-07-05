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
  class="block rounded-md border border-gray-200 p-4 hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-gray-800 transition-colors"
>
  <div class="flex items-center justify-between">
    <span class="font-medium text-gray-900 dark:text-white">{server.name}</span>
    <span class="text-xs text-gray-500 dark:text-gray-400">{modeLabel[server.mode]}</span>
  </div>
  <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
    {server.base_server_type} ⇄ {server.top_server_type}
  </div>
</a>
