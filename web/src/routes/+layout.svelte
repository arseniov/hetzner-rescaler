<script lang="ts">
  import '../app.css';
  import { onNavigate, afterNavigate } from '$app/navigation';
  import { Menu } from 'lucide-svelte';
  import { ensureSession, isAuthenticated } from '$lib/stores/auth.svelte';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';

  let { children } = $props();

  // Mobile drawer state. Owned at the layout so the hamburger (rendered
  // here) and the Sidebar component share the same instance. `open` is
  // local $state; Sidebar reads it as a prop.
  let drawerOpen = $state(false);
  const closeDrawer = () => (drawerOpen = false);
  const toggleDrawer = () => (drawerOpen = !drawerOpen);

  // The server-side hook (web/src/hooks.server.ts) already gates every
  // page by session, so this client-side nav handler only needs to keep
  // the in-memory `_user` cache fresh and connect the SSE stream.
  onNavigate(async (nav) => {
    if (typeof window === 'undefined') return;
    if (nav.to?.route.id === '/login') return;
    await ensureSession();
    // connect() is idempotent — short-circuits if EventSource is
    // already attached.
    eventsStream.connect();
  });

  // Close the mobile drawer on every route change so navigating via the
  // sidebar links (which close the drawer themselves via the link click
  // handlers) also covers programmatic navigations and history-driven
  // changes.
  afterNavigate(() => {
    drawerOpen = false;
  });
</script>

<div class="min-h-screen bg-background text-foreground">
  {#if isAuthenticated()}
    <!-- Mobile-only hamburger that opens the sidebar drawer. Hidden
         from `md` upward because the sidebar becomes a fixed column. -->
    <button
      type="button"
      onclick={toggleDrawer}
      aria-label="Open navigation"
      aria-expanded={drawerOpen}
      class="fixed left-3 top-3 z-30 inline-flex size-9 items-center justify-center rounded-md border border-border bg-card text-foreground hover:bg-muted md:hidden"
    >
      <Menu class="size-4" strokeWidth={1.5} aria-hidden="true" />
    </button>

    <Sidebar isOpen={drawerOpen} closeSidebar={closeDrawer} />

    <!-- On desktop the sidebar is fixed at w-64; `md:ml-64` clears the
         content column. On mobile the sidebar overlays the canvas (a
         dialog drawer), so no offset is needed. -->
    <main class="min-w-0 min-h-screen md:ml-64">
      {@render children?.()}
    </main>
  {:else}
    <main>
      {@render children?.()}
    </main>
  {/if}
</div>
