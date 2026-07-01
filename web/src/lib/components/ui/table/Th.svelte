<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';
  interface Props { class?: string; children?: Snippet }
  let { class: cls, children }: Props = $props();
</script>
<th class={cn('h-10 px-2 text-left align-middle font-medium text-muted-foreground', cls)}>{@render children?.()}</th>
