<script lang="ts">
  import type { Snippet } from 'svelte';
  interface Props { children?: Snippet }
  let { children }: Props = $props();
</script>
<tbody class="[&_tr:last-child]:border-0">{@render children?.()}</tbody>
