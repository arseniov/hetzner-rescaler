<script lang="ts">
  import { Dialog } from 'bits-ui';
  import { X } from 'lucide-svelte';
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';

  type Props = {
    open: boolean;
    title: string;
    description?: string;
    /** Tailwind width class; defaults to `max-w-md`. */
    size?: 'sm' | 'md' | 'lg';
    children?: Snippet;
    footer?: Snippet;
  };

  let {
    open = $bindable(),
    title,
    description,
    size = 'md',
    children,
    footer
  }: Props = $props();

  // `size` maps to a max-width. The panel itself stays inside the
  // viewport with px-4 padding so the modal never touches the edge on
  // narrow screens.
  const sizeClass = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-2xl'
  } as const;
</script>

<!--
  Dialog — bits-ui Dialog wrapped in our tokens. The overlay is a
  soft black wash (no blur — that would be too decorative for this
  aesthetic). The panel is a flat hairline-bordered card on a solid
  background. No drop shadow.

  The close button sits in the top-right corner; the title sits next
  to it. The body content is rendered from a snippet passed via
  `children`; the footer (cancel / save actions) is rendered from
  `footer` so the consumer controls its own button row.

  The outer `Dialog.Root` is bindable on `open` so consumers can
  close it after a successful submit by setting `open = false`.
-->
<Dialog.Root bind:open>
  <Dialog.Portal>
    <Dialog.Overlay
      class="fixed inset-0 z-40 bg-foreground/40"
    />
    <Dialog.Content
      class={cn(
        'fixed left-1/2 top-1/2 z-50 w-[calc(100vw-2rem)] -translate-x-1/2 -translate-y-1/2',
        'rounded-md border border-border bg-card p-5',
        sizeClass[size]
      )}
    >
      <header class="mb-4 flex items-start justify-between gap-4">
        <div class="min-w-0">
          <Dialog.Title class="font-display text-base font-semibold tracking-tight text-foreground">
            {title}
          </Dialog.Title>
          {#if description}
            <Dialog.Description class="mt-1 text-sm text-muted-foreground">
              {description}
            </Dialog.Description>
          {/if}
        </div>
        <Dialog.Close
          aria-label="Close"
          class="inline-flex size-7 shrink-0 items-center justify-center rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
        >
          <X class="size-4" strokeWidth={1.5} aria-hidden="true" />
        </Dialog.Close>
      </header>

      {#if children}
        <div class="text-sm text-foreground">
          {@render children()}
        </div>
      {/if}

      {#if footer}
        <footer class="mt-5 flex justify-end gap-2 border-t border-border pt-4">
          {@render footer()}
        </footer>
      {/if}
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>