<script lang="ts">
  import { cn } from '$lib/utils';
  import type { Snippet } from 'svelte';

  type Props = {
    /** Identifier of the currently active tab. */
    value: string;
    /** List of tabs to render. Order matters — render order = tab order. */
    tabs: { value: string; label: string }[];
    /** Snippet rendered when a tab is selected. Receives the active tab value. */
    children: Snippet<[string]>;
  };

  let { value = $bindable(), tabs, children }: Props = $props();
</script>

<!--
  Tabs — restrained underline indicator. The active tab gets a
  foreground-coloured underline; inactive tabs are muted but readable
  on hover. No animation, no chrome. The tab list is a single row of
  text buttons separated by a hairline bottom border.

  Implementation note: `value` is bindable so a page can store the
  active tab in its own $state and persist it across re-renders or
  deep-link it from the URL if needed.
-->
<div>
  <div
    role="tablist"
    class="flex border-b border-border"
  >
    {#each tabs as tab (tab.value)}
      {@const active = tab.value === value}
      <button
        type="button"
        role="tab"
        aria-selected={active}
        onclick={() => (value = tab.value)}
        class={cn(
          'relative -mb-px px-4 py-2 text-sm font-medium transition-colors',
          active
            ? 'text-foreground'
            : 'text-muted-foreground hover:text-foreground'
        )}
      >
        {tab.label}
        {#if active}
          <span
            class="absolute inset-x-0 bottom-0 h-px bg-foreground"
            aria-hidden="true"
          ></span>
        {/if}
      </button>
    {/each}
  </div>

  <div class="pt-4">
    {@render children(value)}
  </div>
</div>