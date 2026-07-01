<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';
  interface Props { class?: string; children?: Snippet }
  let { class: cls, children }: Props = $props();
</script>
<tr class={cn('border-b transition-colors hover:bg-muted/50', cls)}>{@render children?.()}</tr>
