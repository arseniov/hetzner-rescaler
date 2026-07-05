<script lang="ts">
  import '../app.css';
  import { onNavigate, afterNavigate } from '$app/navigation';
  import { SidebarButton, uiHelpers } from 'flowbite-svelte';
  import { isAuthenticated, ensureSession } from '$lib/stores/auth.svelte';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';

  let { children } = $props();

  // Shared sidebar state: owned here so the mobile toggle button
  // (rendered in this layout) and the Sidebar component can talk to
  // the same instance of `isOpen`. `isOpen` is exposed as a getter
  // that re-reads the underlying $state on every access, so passing
  // it directly to <Sidebar> is reactive. `close` is exposed for
  // use by the Sidebar's CloseButton and the route-change effect.
  const sidebar = uiHelpers();
  const closeSidebar = sidebar.close;

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

  // Close the mobile drawer on every route change so navigating
  // via the sidebar links (which still close the drawer themselves
  // via the SidebarItem click handler) also covers programmatic
  // navigations and history-driven changes.
  afterNavigate(() => {
    closeSidebar();
  });
</script>

<div class="min-h-screen">
  {#if isAuthenticated()}
    <!-- Mobile-only hamburger that opens the sidebar as a drawer.
         The SidebarButton component already hides itself at `md` and
         above via its default theme, so the `md:hidden` here is a
         defensive duplicate. Positioned at the top-left corner. -->
    <SidebarButton
      onclick={sidebar.toggle}
      class="md:hidden fixed top-3 left-3 z-50"
      aria-label="Open navigation"
    />

    <Sidebar isOpen={sidebar.isOpen} {closeSidebar} />

    <!-- On desktop the sidebar is `fixed` and `w-64`. `md:ml-64` clears
         the page content so it doesn't sit underneath the sidebar.
         On mobile the sidebar is hidden by default and overlays the
         page when open, so no offset is needed. -->
    <main class="md:ml-64 min-w-0 min-h-screen">
      {@render children?.()}
    </main>
  {:else}
    <main>
      {@render children?.()}
    </main>
  {/if}
</div>
