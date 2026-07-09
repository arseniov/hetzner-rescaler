<script lang="ts">
  import { Select } from 'bits-ui';
  import { onMount } from 'svelte';
  import { serverTypes } from '$lib/stores/serverTypes.svelte';
  import { m } from '$lib/paraglide/messages.js';
  import type { Server, ServerType } from '$lib/types';
  import { cn, SELECT_TRIGGER_CLASSES } from '$lib/utils';
  import { roleFor, type ServerTypeRole } from '$lib/utils/serverTypeRoles';

  /**
   * ServerTypeSelect — single-value dropdown of Hetzner server types
   * with per-option role chips (CURRENT / BASE / TOP / FALLBACK).
   *
   * Built on bits-ui Select so the same component is used on desktop
   * (mouse) and mobile (touch) without losing keyboard nav or screen
   * reader semantics. The native <select> we had previously was
   * fine for the simple dropdown case but had no path to role markers;
   * a custom dropdown was unavoidable for that affordance.
   *
   * `server` is passed so each option can be marked with the role it
   * plays for this specific server (current / base / top / fallback).
   * The role chips are visual only — they never appear in the bound
   * value, which is just the type code string.
   */
  type Props = {
    value: string;
    /**
     * Server whose role chip set to render. Optional because the
     * dropdown is also used on the "Add server" form, where no server
     * exists yet. When omitted (or null while loading), no role chips
     * are rendered and the dropdown behaves as a plain picker.
     */
    server?: Server | null;
    /**
     * Server location used to drive `serverTypes.load(location)`. When
     * omitted (purely-presentation uses that don't render badges), the
     * onMount load is skipped.
     */
    location?: string;
    id?: string;
    required?: boolean;
    disabled?: boolean;
    onlyAvailable?: boolean;
    class?: string;
    ariaLabel?: string;
    placeholder?: string;
    /**
     * Whether the dropdown is open. Bound from the parent so tests can
     * pre-mount the dropdown contents without simulating a click.
     * Has no effect on user-facing behavior — the trigger toggles this
     * transparently.
     */
    open?: boolean;
  };
  let {
    value = $bindable(),
    server,
    location,
    id,
    required = false,
    disabled = false,
    onlyAvailable = false,
    class: className = '',
    ariaLabel,
    placeholder,
    open = $bindable(false)
  }: Props = $props();

  onMount(() => {
    if (location) {
      serverTypes.load(location).catch(() => { /* loadError is set on the store */ });
    }
  });

  let options = $derived.by<ServerType[]>(() => {
    const all = onlyAvailable ? serverTypes.types.filter((t) => t.available) : serverTypes.types;
    return [...all].sort((a, b) => {
      if (a.available !== b.available) return a.available ? -1 : 1;
      return a.name.localeCompare(b.name);
    });
  });

  // Role classification lives in a pure utility so it can be tested
  // without going through bits-ui's portal-based dropdown (which
  // jsdom can't fully exercise). See web/src/lib/utils/serverTypeRoles.ts.
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

  // bits-ui Select needs a non-empty value for the trigger label.
  let selected = $derived(
    options.find((o) => o.name === value) ??
      (value ? { name: value, available: false } as ServerType : null)
  );
</script>

<Select.Root type="single" bind:value={() => value, (v) => (value = v ?? '')} bind:open {disabled}>
  <Select.Trigger
    {id}
    {disabled}
    aria-label={ariaLabel}
    class={cn(SELECT_TRIGGER_CLASSES, className)}
  >
    <Select.Value placeholder={placeholder ?? (required ? '' : '—')}>
      {#snippet child({ props })}
        {#if selected}
          <span {...props}>{selected.name}</span>
        {:else}
          <span {...props}>{placeholder ?? (required ? '' : '—')}</span>
        {/if}
      {/snippet}
    </Select.Value>
    <span class="ml-2 text-muted-foreground" aria-hidden="true">▾</span>
  </Select.Trigger>
  <Select.Portal>
    <Select.Content
      class="z-50 max-h-72 overflow-y-auto rounded-md border border-border bg-popover p-1 text-popover-foreground shadow-md"
    >
      <Select.Viewport>
        {#if !required && !value}
          <Select.Item value="" label="—" class="px-2 py-1.5 text-sm data-[highlighted]:bg-muted">—</Select.Item>
        {/if}
        {#each options as t (t.name)}
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
            {#if !t.available}
              <span class="ml-auto inline-flex shrink-0 items-center rounded-sm border border-border bg-muted px-1 py-0.5 font-mono text-[9px] uppercase tracking-wider text-muted-foreground">
                {m.server_type_unavailable()}
              </span>
            {/if}
            <span class="font-mono">{t.name}</span>
            <span class="text-xs text-muted-foreground">· {describe(t)}</span>
          </Select.Item>
        {/each}
      </Select.Viewport>
    </Select.Content>
  </Select.Portal>
</Select.Root>
