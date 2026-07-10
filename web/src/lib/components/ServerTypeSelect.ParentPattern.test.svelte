<!--
  Test wrapper for ServerTypeSelect. Simulates the edit-page pattern:
  the parent has a `server` $state and passes `location={server?.location}`
  as a derived expression. Used to verify the $effect in ServerTypeSelect
  fires when the parent's derived location flips from undefined to a
  real value, NOT when location is passed as a direct prop.

  Lives next to the .test file so the test can import it.
-->
<script lang="ts">
  import ServerTypeSelect from './ServerTypeSelect.svelte';
  import type { Server } from '$lib/types';

  let { server }: { server: Server | null } = $props();
</script>

<ServerTypeSelect
  value=""
  {server}
  location={server?.location}
  id="parent-derived"
/>
