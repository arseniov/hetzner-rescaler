<script lang="ts">
  import { Card } from 'flowbite-svelte';

  interface Props {
    label: string;
    value: string | number;
    hint?: string;
  }
  let { label, value, hint }: Props = $props();
</script>

<Card class="flex-1">
  <p class="text-sm text-gray-600 dark:text-gray-400">{label}</p>
  <p class="mt-1 text-3xl font-semibold text-gray-900 dark:text-white">{value}</p>
  {#if hint}
    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{hint}</p>
  {/if}
</Card>