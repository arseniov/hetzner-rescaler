<script lang="ts">
  import '../app.css';
  import { onNavigate } from '$app/navigation';
  import { isAuthenticated, ensureSession } from '$lib/stores/auth.svelte';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';

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
      return;
    }
    // Open the SSE stream once we know the user is authenticated.
    // The connect() call is idempotent — it short-circuits if the
    // EventSource is already attached.
    eventsStream.connect();
  });
</script>

<div class="flex min-h-screen">
  {#if isAuthenticated()}
    <Sidebar />
  {/if}
  <main class="flex-1 min-w-0">
    {@render children?.()}
  </main>
</div>