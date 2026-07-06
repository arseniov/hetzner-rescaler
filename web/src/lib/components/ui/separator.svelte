<script lang="ts">
  import { cn } from '$lib/utils';
  import type { HTMLAttributes } from 'svelte/elements';

  type Props = HTMLAttributes<HTMLHRElement> & {
    orientation?: 'horizontal' | 'vertical';
    class?: string;
  };
  let { orientation = 'horizontal', class: className = '', ...rest }: Props = $props();
</script>

<!--
  Separator — a hairline rule. The dashboard uses these to break the
  sidebar into groups (Status, Operations) without side-stripe borders.
-->
<hr
  aria-orientation={orientation}
  class={cn(
    'shrink-0 border-0 bg-border',
    orientation === 'horizontal' ? 'h-px w-full' : 'h-full w-px',
    className
  )}
  {...rest}
/>
