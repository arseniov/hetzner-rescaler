<script lang="ts">
  import '../app.css';
  import { onNavigate, afterNavigate } from '$app/navigation';
  import { browser } from '$app/environment';
  import { page } from '$app/stores';
  import { Menu, PanelLeftClose, PanelLeftOpen } from 'lucide-svelte';
  import { ensureSession, isAuthenticated } from '$lib/stores/auth.svelte';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';

  let { children } = $props();

  // Two pieces of sidebar state, each meaningful on its own viewport:
  //   - drawerOpen:        mobile-only Dialog drawer
  //   - sidebarCollapsed:  desktop-only "hide the fixed column"
  // The single toggle button in the header dispatches to whichever
  // one is appropriate for the current viewport — see toggleSidebar().
  let drawerOpen = $state(false);
  let sidebarCollapsed = $state(false);

  const closeDrawer = () => (drawerOpen = false);
  const closeDesktopSidebar = () => (sidebarCollapsed = true);

  // The viewport dispatch happens at click time. We don't bind to a
  // resize listener because the toggle is a deliberate operator action
  // (not a reactive system): the viewport check only governs which
  // state flips when the button is pressed, not which state is read on
  // every render.
  function toggleSidebar() {
    if (browser && window.matchMedia('(min-width: 768px)').matches) {
      sidebarCollapsed = !sidebarCollapsed;
    } else {
      drawerOpen = !drawerOpen;
    }
  }

  // Short label for the top bar. Derived from the current URL path —
  // keeps the bar useful without forcing each page to pass a title
  // down. Falls back to the app name on the dashboard.
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
      TOP BAR — visible on every screen size, sticky at the top. On
      mobile it carries the page title so the operator always knows
      where they are without opening the drawer. On desktop it carries
      only the sidebar toggle: the page's own <h1> sits in the content
      column as the primary heading, so duplicating it here would be
      redundant.

      The sidebar toggle is always rendered (no `md:hidden`). The button
      dispatches to the mobile drawer state or the desktop collapse
      state based on viewport — see toggleSidebar().
    -->
    <header
      class="sticky top-0 z-20 flex h-12 items-center gap-3 border-b border-border bg-background/95 px-3 backdrop-blur"
    >
      <button
        type="button"
        onclick={toggleSidebar}
        aria-label="Toggle navigation"
        aria-expanded={!sidebarCollapsed}
        class="inline-flex size-8 items-center justify-center rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
      >
        {#if sidebarCollapsed}
          <PanelLeftOpen class="size-4" strokeWidth={1.5} aria-hidden="true" />
        {:else}
          <PanelLeftClose class="size-4" strokeWidth={1.5} aria-hidden="true" />
        {/if}
      </button>
      <!-- Page title — shown on every screen so the operator knows
           where they are. The dashboard h1 is also visible in the
           content column on desktop; that's intentional (page title
           in the chrome is a navigation aid, the in-page h1 is the
           section heading). -->
      <span class="font-display text-sm font-semibold tracking-tight text-foreground truncate">
        {pageTitle}
      </span>
    </header>

    <!--
      Sidebar. The component handles both modes:
        - Mobile drawer: `isOpen` flips the Dialog open/closed.
        - Desktop fixed column: hidden when `collapsed` is true.

      The component receives both props so the same instance handles
      both viewports without duplication.
    -->
    <Sidebar
      isOpen={drawerOpen}
      collapsed={sidebarCollapsed}
      closeSidebar={closeDrawer}
      closeDesktopSidebar={closeDesktopSidebar}
    />

    <!--
      Main content column. On desktop the sidebar reserves w-64 of the
      left edge; `md:ml-64` clears that space. When the sidebar is
      collapsed on desktop, the margin drops to 0 so the content
      expands. On mobile the sidebar overlays the canvas (a Dialog
      drawer), so no offset is ever needed.
    -->
    <main
      class="min-w-0 min-h-screen {sidebarCollapsed ? '' : 'md:ml-64'} transition-[margin]"
    >
      {@render children?.()}
    </main>
  {:else}
    <main>
      {@render children?.()}
    </main>
  {/if}
</div>