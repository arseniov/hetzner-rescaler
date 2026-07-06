<script lang="ts">
  import { cn } from '$lib/utils';
  import type { HTMLAttributes } from 'svelte/elements';

  type Props = HTMLAttributes<HTMLDivElement> & { class?: string };
  let { class: className = '', ...rest }: Props = $props();
</script>

<!--
  Skeleton — muted rectangle with subtle pulse. Use to indicate
  loading state for KPIs and chart panels. Avoid amber or accent
  here: an operator-facing loading state should be quiet.
-->
<div
  class={cn(
    'rounded-md bg-muted',
    // Subtle pulse — slower than animate-pulse, less distracting.
    'animate-pulse',
    className
  )}
  {...rest}
></div>
