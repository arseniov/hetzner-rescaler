<script lang="ts">
  import { isAuthenticated, signOut } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/ui/button.svelte';

  async function handleSignOut() {
    await signOut();
    await goto('/login');
  }
</script>

<nav class="flex items-center justify-between border-b border-border px-6 py-3">
  <a href="/" class="font-semibold text-lg">Hetzner Rescaler</a>
  <div class="flex items-center gap-4">
    {#if isAuthenticated()}
      <a href="/" class="text-sm hover:underline">Dashboard</a>
      <a href="/projects" class="text-sm hover:underline">Projects</a>
      <a href="/servers" class="text-sm hover:underline">Servers</a>
      <a href="/events" class="text-sm hover:underline">Events</a>
      <Button variant="ghost" size="sm" onclick={handleSignOut}>Sign out</Button>
    {/if}
  </div>
</nav>
