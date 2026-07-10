<!--
  ServerTypeOption — the inner content of a single dropdown row used
  by both ServerTypeSelect (single value, base/top selects) and
  ServerTypeMultiSelect (the fallback chain's add dropdown).

  Layout: the type name and a short description sit on the LEFT; the
  role chip (CURRENT / BASE / TOP / FALLBACK) and the Unavailable
  badge sit on the RIGHT, glued to the far edge. The description
  gets `flex-1 truncate` so it takes the remaining space and the
  badges stack against the right wall regardless of how long the
  description is.

  Both ServerTypeSelect and ServerTypeMultiSelect render this inside
  a bits-ui <Select.Item class="flex w-full …"> wrapper, so the outer
  flex context already exists — this component contributes the
  children, not the row container.

  The role classification is computed via the shared `roleFor` helper
  in `$lib/utils/serverTypeRoles`, kept pure so unit tests can cover
  the priority rules without going through a bits-ui portal.
-->
<script lang="ts">
  import type { Server, ServerType } from '$lib/types';
  import { cn } from '$lib/utils';
  import { m } from '$lib/paraglide/messages.js';
  import { roleFor, type ServerTypeRole } from '$lib/utils/serverTypeRoles';

  type Props = {
    type: ServerType;
    /** Server whose role this type plays; null/undefined omits the role chip. */
    server?: Server | null;
  };
  let { type: t, server }: Props = $props();

  const role = $derived(roleFor(t, server ?? null));

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
</script>

<!--
  Row layout (left to right):
    [type name]  [· description (flex-1, truncated)]
                                 [role chip — ml-auto]   [Unavailable badge]
  The `ml-auto` rides on whichever badge is first; if only one
  badge is present, it still gets pushed to the far right.
-->
<span class="font-mono">{t.name}</span>
<span class="flex-1 truncate text-xs text-muted-foreground">· {describe(t)}</span>
{#if role}
  <span
    class={cn(
      'ml-auto inline-flex shrink-0 items-center rounded-sm border px-1 py-0.5 font-mono text-[9px] uppercase tracking-wider',
      chipClass(role)
    )}
  >
    {chipLabel(role)}
  </span>
{/if}
{#if !t.available}
  <span
    class={cn(
      !role && 'ml-auto',
      'inline-flex shrink-0 items-center rounded-sm border border-border bg-muted px-1 py-0.5 font-mono text-[9px] uppercase tracking-wider text-muted-foreground'
    )}
  >
    {m.server_type_unavailable()}
  </span>
{/if}