<script lang="ts">
  import type { Snippet } from 'svelte';
  interface Props { children?: Snippet }
  let { children }: Props = $props();
</script>
<thead class="[&_tr]:border-b">{@render children?.()}</thead>
