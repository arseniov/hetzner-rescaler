<script lang="ts">
  import { onMount } from 'svelte';
  import { serverTypes } from '$lib/stores/serverTypes.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import type { ServerType } from '$lib/types';
  import { cn } from '$lib/utils';

  /**
   * ServerTypeSelect — single-value dropdown of Hetzner server types.
   *
   * Drives three operator-facing forms (edit server base/top, add
   * window target_type) that previously took a free-text type code.
   * The free-text path was a real footgun: typing "cpxl1" instead of
   * "cpx11" silently creates a server with an unknown type that the
   * Hetzner API will reject on first rescale.
   *
   * The dropdown is a styled native <select> rather than a bits-ui
   * Combobox: there are only a few dozen types, the picker doesn't
   * need to be searchable, and a native control keeps keyboard
   * navigation and screen-reader announcements working without extra
   * ARIA wiring.
   *
   * `value` is two-way bound via `$bindable()` so consumers can
   * `bind:value` exactly like a native <select>. The component does
   * not own the form state — it just reflects it.
   */
  type Props = {
    value: string;
    id?: string;
    required?: boolean;
    disabled?: boolean;
    /** Filter to only available types (default: false; sold-out types still listed). */
    onlyAvailable?: boolean;
    /** Optional class merged onto the <select> for layout alignment. */
    class?: string;
    /** Optional aria-label when no Label sibling is provided. */
    ariaLabel?: string;
  };
  let {
    value = $bindable(),
    id,
    required = false,
    disabled = false,
    onlyAvailable = false,
    class: className = '',
    ariaLabel
  }: Props = $props();

  onMount(() => {
    // Fire-and-forget — `load()` is idempotent and shared across all
    // ServerTypeSelect instances. The empty placeholder option below
    // covers the "still loading" UX.
    serverTypes.load().catch(() => {
      // loadError is set on the store; the option list will simply
      // remain empty until the operator retries by editing again.
    });
  });

  // Recompute option list whenever the catalog changes. `$derived` over
  // the store field means Svelte's reactivity tracks it across the
  // module boundary.
  let options = $derived.by<ServerType[]>(() => {
    const all = onlyAvailable ? serverTypes.available() : serverTypes.types;
    // Sort: available types first (sensible default for "pick one"),
    // then alphabetical within each bucket so the list is stable.
    return [...all].sort((a, b) => {
      if (a.available !== b.available) return a.available ? -1 : 1;
      return a.name.localeCompare(b.name);
    });
  });

  // The current `value` may be a code the operator typed before the
  // catalog loaded (or a Hetzner code the API returned that we don't
  // list — e.g. a deprecated type). When the catalog doesn't contain
  // the value, we still want to display it in the dropdown so the
  // operator sees what's currently set.
  let hasValue = $derived(value !== '' && value !== undefined);

  // Pretty-print one option. Memory + cores + price are only included
  // when the catalog row actually has them; older Hetzner responses
  // (or the minimal `{name, available}` shape from legacy callers)
  // collapse to just the type code.
  function describe(t: ServerType): string {
    const bits: string[] = [];
    if (typeof t.cores === 'number' && typeof t.memory_gb === 'number') {
      bits.push(`${t.cores} cores · ${t.memory_gb} GB`);
    } else if (typeof t.cores === 'number') {
      bits.push(`${t.cores} cores`);
    }
    if (typeof t.price_monthly_eur === 'number' && t.price_monthly_eur > 0) {
      bits.push(`€${t.price_monthly_eur.toFixed(2)}/mo`);
    }
    if (t.description) bits.push(t.description);
    return bits.length > 0 ? bits.join(' · ') : t.name;
  }
</script>

<!--
  Native <select> styled to match the Input primitive (hairline
  border, ring-token focus, mono type code, no glow). When the catalog
  is still loading we render a single disabled placeholder so the
  field is visually identical to the other inputs — no spinners, no
  separate "loading" state to design.
-->
<select
  {id}
  {required}
  {disabled}
  bind:value
  aria-label={ariaLabel}
  class={cn(
    'flex h-9 rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground',
    'focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background',
    'disabled:cursor-not-allowed disabled:opacity-50',
    className
  )}
>
  {#if serverTypes.loadedAt === null && options.length === 0}
    <option value="" disabled>{m.server_type_select_loading()}</option>
  {:else if serverTypes.loadError && options.length === 0}
    <option value="" disabled>{m.server_type_select_unavailable()}</option>
  {:else}
    {#if !required && !hasValue}
      <option value="">—</option>
    {/if}
    {#if hasValue && !options.some((o) => o.name === value)}
      <option value={value}>{value} (unknown)</option>
    {/if}
    {#each options as t (t.name)}
      <option value={t.name}>{t.name} · {describe(t)}</option>
    {/each}
  {/if}
</select>