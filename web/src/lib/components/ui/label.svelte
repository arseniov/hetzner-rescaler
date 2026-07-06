<script lang="ts">
  import { cn } from '$lib/utils';
  import type { HTMLLabelAttributes } from 'svelte/elements';

  type Props = HTMLLabelAttributes & { class?: string };
  let { class: className = '', children, ...rest }: Props = $props();
</script>

<!--
  Label — peer-styled (peer-disabled / peer-focus-visible) so it picks
  up the input's disabled / focus state without per-component wiring.
-->
<label
  class={cn(
    'text-sm font-medium leading-none text-foreground',
    'peer-disabled:cursor-not-allowed peer-disabled:opacity-50',
    className
  )}
  {...rest}
>
  {@render children?.()}
</label>
