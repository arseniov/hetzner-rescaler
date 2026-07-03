<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';

  interface Props {
    variant?: 'default' | 'destructive';
    children?: Snippet;
    class?: string;
  }

  let { variant = 'default', class: cls, children }: Props = $props();
</script>

<div
  role="alert"
  class={cn(
    'rounded-md border p-4 text-sm',
    variant === 'destructive'
      ? 'border-destructive bg-destructive/10 text-destructive'
      : 'border-border bg-muted text-foreground',
    cls
  )}
>
  {@render children?.()}
</div>
