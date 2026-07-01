<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';
  interface Props { class?: string; children?: Snippet }
  let { class: cls, children }: Props = $props();
</script>
<td class={cn('p-2 align-middle', cls)}>{@render children?.()}</td>
