<script lang="ts">
  import { cn } from '$lib/utils';
  import type { HTMLInputAttributes } from 'svelte/elements';

  // `value` is bindable so consumers can `bind:value={email}` the way
  // they would on a native <input>. Svelte 5 requires the explicit
  // `$bindable()` annotation for props that are two-way bound.
  type Props = HTMLInputAttributes & { class?: string };
  let { class: className = '', value = $bindable(), ...rest }: Props = $props();
</script>

<!--
  Input — single hairline border, tight radius, no glow on focus.
  Focus state uses `ring` token (cool muted, not accent). The accent
  is reserved for state changes — focus on a form input is too common
  to deserve it.
-->
<input
  class={cn(
    'flex h-9 w-full rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground',
    'placeholder:text-muted-foreground',
    'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background',
    'disabled:cursor-not-allowed disabled:opacity-50',
    className
  )}
  bind:value
  {...rest}
/>
