<script lang="ts" module>
  import { tv, type VariantProps } from 'tailwind-variants';

  export const alertVariants = tv({
    base: 'flex flex-col gap-1 rounded-md border border-border bg-muted px-4 py-3 text-sm',
    variants: {
      variant: {
        // Default: muted background, hairline border. No state change.
        default: '',
        // Destructive: reserved for actual failures (rescale error,
        // failed fetch). Uses the destructive token directly.
        destructive:
          'border-destructive/30 bg-destructive/10 text-foreground',
        // Warning: signup-disabled, soft warnings. Subtle, not amber —
        // the amber accent is reserved for "this just changed" state.
        warning: 'border-warning/30 bg-warning/10 text-foreground'
      }
    },
    defaultVariants: { variant: 'default' }
  });
  export type AlertVariant = VariantProps<typeof alertVariants>['variant'];
</script>

<script lang="ts">
  import { cn } from '$lib/utils';
  import type { HTMLAttributes } from 'svelte/elements';

  type Props = HTMLAttributes<HTMLDivElement> & {
    variant?: AlertVariant;
    class?: string;
  };
  let { variant = 'default', class: className = '', children, ...rest }: Props = $props();
</script>

<div
  role={variant === 'destructive' ? 'alert' : undefined}
  class={cn(alertVariants({ variant }), className)}
  {...rest}
>
  {@render children?.()}
</div>
