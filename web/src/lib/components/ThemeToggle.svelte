<script lang="ts">
  import { theme } from '$lib/stores/theme.svelte';
  import { m } from '$lib/paraglide/messages.js';
</script>

<button
  type="button"
  class="text-sm rounded p-2 hover:bg-gray-100 dark:hover:bg-gray-700"
  aria-label={theme.current === 'dark' ? m.theme_toggle_light() : m.theme_toggle_dark()}
  onclick={() => theme.toggle()}
>
  {theme.current === 'dark' ? '☀' : '🌙'}
</button>