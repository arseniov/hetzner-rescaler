<script lang="ts">
  import { Card } from 'flowbite-svelte';

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

<Card class="flex-1 border-0">
  <p class="text-sm text-gray-600 dark:text-gray-400">{label}</p>
  <p class="mt-1 text-3xl font-semibold text-gray-900 dark:text-white">
    {#if loading}
      <span class="inline-block animate-pulse text-gray-400">…</span>
    {:else}
      {display}
    {/if}
  </p>
  {#if hint}
    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{hint}</p>
  {/if}
</Card>