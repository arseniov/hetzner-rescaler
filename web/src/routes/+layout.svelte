<script lang="ts">
  import '../app.css';
  import { onNavigate, afterNavigate } from '$app/navigation';
  import { page } from '$app/stores';
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

  // Short label for the mobile top bar. Derived from the current URL
  // path — keeps the bar useful without forcing each page to pass a
  // title down. Falls back to the app name on the dashboard.
  let pageTitle = $derived.by(() => {
    const path = $page.url.pathname;
    if (path === '/') return 'Dashboard';
    if (path === '/projects') return 'Projects';
    if (path.startsWith('/projects/')) return 'Project';
    if (path === '/servers') return 'Servers';
    if (path.startsWith('/servers/')) {
      if (path.endsWith('/edit')) return 'Edit server';
      if (path.endsWith('/windows')) return 'Windows';
      return 'Server';
    }
    if (path === '/events') return 'Events';
    if (path === '/status/health') return 'System health';
    if (path === '/status/servers') return 'Server status';
    return '';
  });

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
    <!--
      MOBILE TOP BAR — always rendered on mobile when the operator is
      signed in. Holds the hamburger that toggles the sidebar drawer.

      Why a sticky bar instead of a fixed-position button floating over
      content:
        1. The bar reserves its own height (h-12), so a page's `h1` sits
           below it and never overlaps the hamburger — no per-page
           padding hack needed.
        2. `sticky top-0` keeps the bar pinned while the user scrolls a
           long list (servers, events) so navigation is always one tap
           away.
        3. The bar also carries the page label so the operator always
           knows where they are without scanning the drawer.

      The bar hides from `md` upward because the sidebar becomes a
      fixed column on desktop and a floating trigger would compete
      with it.
    -->
    <header
      class="sticky top-0 z-20 flex h-12 items-center gap-3 border-b border-border bg-background/95 px-3 backdrop-blur md:hidden"
    >
      <button
        type="button"
        onclick={toggleDrawer}
        aria-label="Open navigation"
        aria-expanded={drawerOpen}
        class="inline-flex size-8 items-center justify-center rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
      >
        <Menu class="size-4" strokeWidth={1.5} aria-hidden="true" />
      </button>
      <span class="font-display text-sm font-semibold tracking-tight text-foreground truncate">
        {pageTitle}
      </span>
    </header>

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