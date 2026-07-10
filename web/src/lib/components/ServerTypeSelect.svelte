<script lang="ts">
  import { Select } from 'bits-ui';
  import { serverTypes } from '$lib/stores/serverTypes.svelte';
  import type { Server, ServerType } from '$lib/types';
  import { cn, SELECT_TRIGGER_CLASSES } from '$lib/utils';
  import ServerTypeOption from './ServerTypeOption.svelte';

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
   *
   * The dropdown row layout (name + description left, role chip +
   * Unavailable badge right) is rendered by `ServerTypeOption` so
   * this component and `ServerTypeMultiSelect` stay byte-identical
   * across the 3 use sites (base, top, fallback). Don't fork it.
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
     * load is skipped. The component uses `$effect`, so this prop can
     * transition from undefined → a real value after mount.
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

  // Load the catalog whenever `location` becomes a non-empty string.
  // Using `$effect` (not `onMount`) is intentional: on pages like the
  // server-edit form, the parent renders this component *before* the
  // server has been fetched from the API, so `location` is initially
  // undefined. `onMount` would fire only once with `undefined` and
  // never re-fire when the server finally resolves. `$effect`
  // re-runs whenever `location` changes, so it picks up the resolved
  // server and triggers the catalog load.
  $effect(() => {
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
          {@const describeLabel = `${t.name}${typeof t.cores === 'number' && typeof t.memory_gb === 'number' ? ` · ${t.cores} cores · ${t.memory_gb} GB` : ''}`}
          <Select.Item
            value={t.name}
            label={describeLabel}
            class="flex w-full cursor-pointer items-center gap-2 rounded-sm px-2 py-1.5 text-sm data-[highlighted]:bg-muted data-[state=checked]:bg-muted/60"
          >
            <ServerTypeOption type={t} {server} />
          </Select.Item>
        {/each}
      </Select.Viewport>
    </Select.Content>
  </Select.Portal>
</Select.Root>
