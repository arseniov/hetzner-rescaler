<script lang="ts" module>
  import { tv, type VariantProps } from 'tailwind-variants';

  export const buttonVariants = tv({
    base: 'inline-flex items-center justify-center gap-2 whitespace-nowrap font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-50 [&_svg]:size-4 [&_svg]:shrink-0',
    variants: {
      variant: {
        // Default is restrained — outline + hover fill. The dashboard
        // never wants loud primary buttons except for hard confirmations.
        default:
          'border border-border bg-transparent text-foreground hover:bg-muted',
        // Filled muted for secondary actions (cancel, back).
        secondary:
          'bg-muted text-foreground border border-transparent hover:bg-secondary',
        // Strong primary — only for the one hard-confirmation action on
        // a screen (sign in, save, submit).
        primary:
          'bg-primary text-primary-foreground border border-primary hover:opacity-90',
        // Ghost — text-only with hover bg, for inline links.
        ghost: 'text-foreground hover:bg-muted',
        // Destructive — red, used for delete actions.
        destructive:
          'bg-destructive text-destructive-foreground border border-transparent hover:opacity-90'
      },
      size: {
        sm: 'h-8 px-3 text-sm rounded-sm',
        md: 'h-9 px-4 text-sm rounded-md',
        lg: 'h-10 px-5 text-base rounded-md',
        icon: 'size-9 rounded-md'
      }
    },
    defaultVariants: {
      variant: 'default',
      size: 'md'
    }
  });

  export type ButtonVariant = VariantProps<typeof buttonVariants>['variant'];
  export type ButtonSize = VariantProps<typeof buttonVariants>['size'];
</script>

<script lang="ts">
  import type { HTMLButtonAttributes, HTMLAnchorAttributes } from 'svelte/elements';
  import { cn } from '$lib/utils';

  type Props = HTMLButtonAttributes &
    HTMLAnchorAttributes & {
      variant?: ButtonVariant;
      size?: ButtonSize;
      class?: string;
      href?: string;
    };

  let {
    variant = 'default',
    size = 'md',
    class: className = '',
    href,
    type = 'button',
    children,
    ...rest
  }: Props = $props();

  // If `href` is provided, render as an anchor styled the same way —
  // useful for "Sign in" links and other in-page navigations.
</script>

{#if href}
  <a {href} class={cn(buttonVariants({ variant, size }), className)} {...rest}>
    {@render children?.()}
  </a>
{:else}
  <button {type} class={cn(buttonVariants({ variant, size }), className)} {...rest}>
    {@render children?.()}
  </button>
{/if}
