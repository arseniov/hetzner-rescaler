<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';
  import type { HTMLButtonAttributes } from 'svelte/elements';

  type Variant = 'default' | 'destructive' | 'outline' | 'ghost';
  type Size = 'sm' | 'md' | 'lg';

  interface Props extends HTMLButtonAttributes {
    variant?: Variant;
    size?: Size;
    children?: Snippet;
  }

  let { variant = 'default', size = 'md', class: cls, children, ...rest }: Props = $props();

  const variants: Record<Variant, string> = {
    default: 'bg-primary text-primary-foreground hover:opacity-90',
    destructive: 'bg-destructive text-destructive-foreground hover:opacity-90',
    outline: 'border border-border bg-transparent hover:bg-muted',
    ghost: 'bg-transparent hover:bg-muted'
  };
  const sizes: Record<Size, string> = {
    sm: 'h-8 px-3 text-sm',
    md: 'h-10 px-4',
    lg: 'h-12 px-6 text-lg'
  };
</script>

<button
  class={cn(
    'inline-flex items-center justify-center rounded-md font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed',
    variants[variant],
    sizes[size],
    cls
  )}
  {...rest}
>
  {@render children?.()}
</button>
