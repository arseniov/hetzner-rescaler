<script lang="ts">
  import { cn } from '$lib/utils';
  import type { HTMLInputAttributes } from 'svelte/elements';

  interface Props extends Omit<HTMLInputAttributes, 'value'> {
    class?: string;
    value?: string | number;
  }

  let { class: cls, value = $bindable(''), ...rest }: Props = $props();
</script>

<input
  bind:value
  class={cn(
    'flex h-10 w-full rounded-md border border-border bg-background px-3 py-2 text-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary',
    cls
  )}
  {...rest}
/>
