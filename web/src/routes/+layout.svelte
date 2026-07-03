<script lang="ts">
  import '../app.css';
  import { onNavigate } from '$app/navigation';
  import { isAuthenticated, ensureSession } from '$lib/stores/auth.svelte';
  import Nav from '$lib/components/Nav.svelte';

  let { children } = $props();

  // Hydrate from the session cookie on the first navigation so the
  // guard below has reliable state.
  onNavigate(async (nav) => {
    if (typeof window === 'undefined') return;
    if (nav.to?.route.id === '/login' || nav.to?.route.id === null) return;
    await ensureSession();
    if (!isAuthenticated()) {
      // replaceState so the back button doesn't trap the user on /login
      // after they sign in.
      window.location.replace('/login');
    }
  });
</script>

<div class="min-h-screen flex flex-col">
  {#if isAuthenticated()}
    <Nav />
  {/if}
  <div class="flex-1">
    {@render children?.()}
  </div>
</div>