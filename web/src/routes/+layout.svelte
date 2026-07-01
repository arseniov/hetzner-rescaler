<script lang="ts">
  import '../app.css';
  import { onNavigate } from '$app/navigation';
  import { isAuthenticated } from '$lib/stores/auth.svelte';
  import Nav from '$lib/components/Nav.svelte';

  let { children } = $props();

  // Client-side route guard: redirect to /login if not authenticated.
  // /login is the only public route.
  onNavigate((nav) => {
    if (typeof window === 'undefined') return;
    if (!isAuthenticated() && nav.to?.route.id !== '/login') {
      // SvelteKit will navigate; we use replaceState so the back button
      // doesn't trap the user on /login after they sign in.
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
