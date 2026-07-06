<script lang="ts" module>
  import { cn } from '$lib/utils';
  import type { HTMLAttributes } from 'svelte/elements';

  /**
   * Card — flat hairline-bordered panel. NEVER has a drop shadow on
   * the dashboard. Use Card.Header / Card.Title / Card.Description /
   * Card.Content / Card.Footer to compose. Sections are spaced by
   * `gap`, not by individual padding.
   */
  export function cardClasses(extra?: string) {
    return cn(
      'flex flex-col gap-4 rounded-md border border-border bg-card text-card-foreground p-4',
      extra
    );
  }
</script>

<script lang="ts">
  type Props = HTMLAttributes<HTMLDivElement> & { class?: string };
  let { class: className = '', children, ...rest }: Props = $props();
</script>

<div class={cn(cardClasses(), className)} {...rest}>
  {@render children?.()}
</div>
