<script lang="ts">
  import { Select } from 'bits-ui';
  import { dndzone, type DndEvent } from 'svelte-dnd-action';
  import { onMount } from 'svelte';
  import { X, GripVertical } from 'lucide-svelte';
  import { serverTypes } from '$lib/stores/serverTypes.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import type { Server, ServerType } from '$lib/types';
  import { cn, SELECT_TRIGGER_CLASSES } from '$lib/utils';
  import { roleFor, type ServerTypeRole } from '$lib/utils/serverTypeRoles';

  /**
   * ServerTypeMultiSelect — drag-drop multiselect for the fallback
   * chain. The bound value is an ordered array of type code strings;
   * the operator can:
   *   - add a type from the custom dropdown (matches ServerTypeSelect
   *     visually so base/top/fallback selectors read as a set)
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
   * `server` (optional) is passed so each option can be marked with
   * the role it plays for this specific server (current / base / top /
   * fallback), matching ServerTypeSelect's affordance.
   */
  type Chip = { id: string; name: string };
  type Props = {
    value: string[];
    excluded?: string[];
    server?: Server | null;
    /**
     * Server location used to drive `serverTypes.load(location)`. When
     * omitted (purely-presentation uses that don't render badges), the
     * onMount load is skipped. Defaults to undefined so callers that
     * don't yet know the server's location (Add-server dialog while
     * typing) don't trigger a load for the wrong region.
     */
    location?: string;
    id?: string;
    disabled?: boolean;
    class?: string;
    ariaLabel?: string;
    emptyLabel?: string;
  };
  let {
    value = $bindable(),
    excluded = [],
    server = null,
    location,
    id,
    disabled = false,
    class: className = '',
    ariaLabel,
    emptyLabel = m.server_type_multiselect_empty()
  }: Props = $props();

  onMount(() => {
    if (location) {
      serverTypes.load(location).catch(() => { /* loadError is set on the store */ });
    }
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
      .sort((a, b) => {
        if (a.available !== b.available) return a.available ? -1 : 1;
        return a.name.localeCompare(b.name);
      })
  );

  // bits-ui Select needs a non-empty current value to drive the
  // bind:value callback. We hold a transient "picked" value: when an
  // item is selected, the setter forwards it to addFromDropdown and
  // resets the transient value so the trigger re-renders with the
  // placeholder for the next pick.
  let pickerValue = $state('');

  // Role classification: reuse the pure helper so the multi-select
  // shows identical chips to ServerTypeSelect.
  const ROLE_FOR = roleFor;

  function chipClass(role: ServerTypeRole): string {
    if (role === 'current') {
      return 'border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-300';
    }
    return 'border-border bg-muted text-muted-foreground';
  }

  function chipLabel(role: ServerTypeRole): string {
    switch (role) {
      case 'current':  return m.server_type_role_current();
      case 'base':     return m.server_type_role_base();
      case 'top':      return m.server_type_role_top();
      case 'fallback': return m.server_type_role_fallback();
    }
  }

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

  const flipDurationMs = 120;
</script>

<!--
  Two-row layout: the chip list (drag-drop) on top and an "add"
  dropdown below. The chip list flips to a draggable style as soon as
  the user picks up a chip; the flip duration is set low so reordering
  feels snappy. The add dropdown uses the same bits-ui Select as
  ServerTypeSelect so all three selectors in a form read as one set.
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
      aria-label={ariaLabel ?? m.server_type_multiselect_chip_list_label()}
    >
      {#each chips as chip (chip.id)}
        <li class="inline-flex items-center gap-1.5 rounded-md border border-border bg-card px-2 py-1 text-sm">
          <GripVertical class="size-3.5 cursor-grab text-muted-foreground" strokeWidth={1.75} aria-hidden="true" />
          <span class="font-mono">{chip.name}</span>
          <button
            type="button"
            class="ml-0.5 rounded-sm p-0.5 text-muted-foreground hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
            aria-label={m.server_type_multiselect_remove_chip({ name: chip.name })}
            onclick={() => remove(chip.name)}
            {disabled}
          >
            <X class="size-3" strokeWidth={1.75} aria-hidden="true" />
          </button>
        </li>
      {/each}
    </ul>
  {/if}

  <Select.Root
    type="single"
    bind:value={
      () => pickerValue,
      (v) => {
        // Reset the transient value first so the trigger re-renders
        // with the placeholder for the next pick. Forward the picked
        // value (which may be '' if the menu was dismissed without
        // selecting) to addFromDropdown.
        pickerValue = '';
        if (v) addFromDropdown(v);
      }
    }
    {disabled}
  >
    <Select.Trigger
      {id}
      {disabled}
      aria-label={m.server_type_multiselect_add_label()}
      class={cn(SELECT_TRIGGER_CLASSES, className)}
    >
      <Select.Value placeholder={m.server_type_multiselect_add_placeholder()}>
        {#snippet child({ props })}
          <span {...props}>{m.server_type_multiselect_add_placeholder()}</span>
        {/snippet}
      </Select.Value>
      <span class="ml-2 text-muted-foreground" aria-hidden="true">▾</span>
    </Select.Trigger>
    <Select.Portal>
      <Select.Content
        class="z-50 max-h-72 overflow-y-auto rounded-md border border-border bg-popover p-1 text-popover-foreground shadow-md"
      >
        <Select.Viewport>
          {#if addableOptions.length === 0}
            <p class="px-2 py-1.5 text-sm text-muted-foreground">
              {m.server_type_multiselect_no_addable()}
            </p>
          {:else}
            {#each addableOptions as t (t.name)}
              {@const role = ROLE_FOR(t, server)}
              <Select.Item
                value={t.name}
                label={`${t.name} · ${describe(t)}`}
                class="flex w-full cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm data-[highlighted]:bg-muted data-[state=checked]:bg-muted/60"
              >
                {#if role}
                  <span
                    class={cn(
                      'inline-flex shrink-0 items-center rounded-sm border px-1 py-0.5 font-mono text-[9px] uppercase tracking-wider',
                      chipClass(role)
                    )}
                  >
                    {chipLabel(role)}
                  </span>
                {/if}
                <span class="font-mono">{t.name}</span>
                <span class="text-xs text-muted-foreground">· {describe(t)}</span>
              </Select.Item>
            {/each}
          {/if}
        </Select.Viewport>
      </Select.Content>
    </Select.Portal>
  </Select.Root>
</div>