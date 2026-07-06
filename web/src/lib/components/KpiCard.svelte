<script lang="ts">
  interface Props {
    label: string;
    value: string | number | null | undefined;
    hint?: string;
    loading?: boolean;
  }
  let { label, value, hint, loading = false }: Props = $props();

  // While loading, show an ellipsis placeholder. Otherwise show the value
  // if present, or an em-dash when null/undefined (i.e. fetch completed
  // but no data was returned, or a KPI is genuinely absent).
  let display = $derived(loading ? '…' : value ?? '—');
</script>

<!-- KPI panel — flat card, hairline border, no shadow. The numeric
     figure uses JetBrains Mono with tabular numerals so digits stay
     column-aligned across cards and don't jitter as values change. -->
<div
  class="flex flex-col gap-1 rounded-md border border-border bg-card px-4 py-3 text-card-foreground"
>
  <p class="text-sm text-muted-foreground">{label}</p>
  {#if loading}
    <p class="tabular animate-pulse text-3xl text-muted-foreground">…</p>
  {:else if value === null || value === undefined}
    <p class="tabular text-3xl text-muted-foreground">—</p>
  {:else}
    <p class="tabular text-3xl font-semibold text-foreground">{display}</p>
  {/if}
  {#if hint}
    <p class="text-xs text-muted-foreground">{hint}</p>
  {/if}
</div>
