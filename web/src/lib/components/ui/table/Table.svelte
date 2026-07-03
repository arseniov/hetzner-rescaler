<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';
  interface Props { class?: string; children?: Snippet }
  let { class: cls, children }: Props = $props();
</script>
<div class={cn('w-full overflow-auto', cls)}>
  <table class="w-full caption-bottom text-sm">{@render children?.()}</table>
</div>
