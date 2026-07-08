<script lang="ts">
  import { dndzone, type DndEvent } from 'svelte-dnd-action';
  import { onMount } from 'svelte';
  import { X, GripVertical } from 'lucide-svelte';
  import { serverTypes } from '$lib/stores/serverTypes.svelte';
  import { cn } from '$lib/utils';

  /**
   * ServerTypeMultiSelect — drag-drop multiselect for the fallback
   * chain. The bound value is an ordered array of type code strings;
   * the operator can:
   *   - add a type from the dropdown (native <select> → append to array)
   *   - drag chips to reorder them
   *   - remove chips via the X button
   *
   * svelte-dnd-action handles both pointer and touch via a single
   * code path; no per-device wiring required.
   *
   * `excluded` lists types that should NOT appear in the add dropdown
   * (e.g. the server's base and top types, which are configured in
   * their own dropdowns and would be contradictory in a fallback).
   *
   * The add dropdown stays a native <select> on purpose: it is a
   * single-tap "pick from a list" action, not a multi-select picker,
   * so the bits-ui Select affordances (search, role chips) would be
   * overkill. The chips themselves are the area that needs dnd.
   */
  type Chip = { id: string; name: string };
  type Props = {
    value: string[];
    excluded?: string[];
    id?: string;
    disabled?: boolean;
    class?: string;
    ariaLabel?: string;
    emptyLabel?: string;
  };
  let {
    value = $bindable(),
    excluded = [],
    id,
    disabled = false,
    class: className = '',
    ariaLabel,
    emptyLabel = 'No fallback types yet. Add one below.'
  }: Props = $props();

  onMount(() => {
    serverTypes.load().catch(() => { /* loadError is set on the store */ });
  });

  // Mirror the bound value into a chips list. The chip `id` is the
  // type code; svelte-dnd-action needs a stable id per item.
  let chips = $state<Chip[]>([]);
  $effect(() => {
    chips = value.map((name) => ({ id: name, name }));
  });

  function onConsider(e: CustomEvent<DndEvent<Chip>>) {
    chips = e.detail.items;
  }

  function onFinalize(e: CustomEvent<DndEvent<Chip>>) {
    chips = e.detail.items;
    value = chips.map((c) => c.name);
  }

  function remove(name: string) {
    value = value.filter((v) => v !== name);
  }

  function addFromDropdown(name: string) {
    if (!name) return;
    if (value.includes(name)) return;
    value = [...value, name];
  }

  let addableOptions = $derived(
    serverTypes.types
      .filter((t) => !value.includes(t.name) && !excluded.includes(t.name))
      .sort((a, b) => a.name.localeCompare(b.name))
  );

  const flipDurationMs = 120;
</script>

<!--
  Two-row layout: the chip list (drag-drop) on top and an "add"
  dropdown below. The chip list flips to a draggable style as soon as
  the user picks up a chip; the flip duration is set low so reordering
  feels snappy.
-->
<div class={cn('flex flex-col gap-3', className)}>
  {#if chips.length === 0}
    <p class="text-xs text-muted-foreground">{emptyLabel}</p>
  {:else}
    <ul
      use:dndzone={{ items: chips, flipDurationMs, type: 'fallback-chips' }}
      onconsider={onConsider}
      onfinalize={onFinalize}
      class="flex flex-wrap gap-2"
      aria-label={ariaLabel ?? 'Fallback chain (drag to reorder)'}
    >
      {#each chips as chip (chip.id)}
        <li class="inline-flex items-center gap-1.5 rounded-md border border-border bg-card px-2 py-1 text-sm">
          <GripVertical class="size-3.5 cursor-grab text-muted-foreground" strokeWidth={1.75} aria-hidden="true" />
          <span class="font-mono">{chip.name}</span>
          <button
            type="button"
            class="ml-0.5 rounded-sm p-0.5 text-muted-foreground hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            aria-label={`Remove ${chip.name}`}
            onclick={() => remove(chip.name)}
            {disabled}
          >
            <X class="size-3" strokeWidth={1.75} aria-hidden="true" />
          </button>
        </li>
      {/each}
    </ul>
  {/if}

  <select
    {id}
    {disabled}
    aria-label="Add a fallback type"
    onchange={(e) => {
      const sel = e.currentTarget as HTMLSelectElement;
      addFromDropdown(sel.value);
      sel.value = '';
    }}
    class="flex h-9 rounded-md border border-border bg-input px-3 py-1 text-sm text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50"
  >
    <option value="">+ Add fallback type…</option>
    {#each addableOptions as t (t.name)}
      <option value={t.name}>{t.name}</option>
    {/each}
  </select>
</div>