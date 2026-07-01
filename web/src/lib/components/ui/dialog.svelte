<script lang="ts">
  import { Dialog as DialogPrimitive } from 'bits-ui';
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';

  interface Props {
    open: boolean;
    title: string;
    description?: string;
    children?: Snippet;
    footer?: Snippet;
    onOpenChange?: (open: boolean) => void;
  }

  let { open = $bindable(false), title, description, children, footer, onOpenChange }: Props = $props();

  function handleOpenChange(next: boolean) {
    open = next;
    onOpenChange?.(next);
  }
</script>

<DialogPrimitive.Root bind:open onOpenChange={handleOpenChange}>
  <DialogPrimitive.Portal>
    <DialogPrimitive.Overlay class="fixed inset-0 z-50 bg-black/50" />
    <DialogPrimitive.Content
      class={cn(
        'fixed left-1/2 top-1/2 z-50 grid w-full max-w-lg -translate-x-1/2 -translate-y-1/2 gap-4 border border-border bg-background p-6 shadow-lg rounded-lg'
      )}
    >
      <DialogPrimitive.Title class="text-lg font-semibold">{title}</DialogPrimitive.Title>
      {#if description}
        <DialogPrimitive.Description class="text-sm text-muted-foreground">{description}</DialogPrimitive.Description>
      {/if}
      <div>{@render children?.()}</div>
      {#if footer}
        <div class="flex justify-end gap-2">{@render footer()}</div>
      {/if}
    </DialogPrimitive.Content>
  </DialogPrimitive.Portal>
</DialogPrimitive.Root>
