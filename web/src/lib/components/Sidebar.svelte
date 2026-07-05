<script lang="ts">
  import { page } from '$app/stores';
  import { Sidebar, SidebarGroup, SidebarItem, SidebarWrapper, CloseButton } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { isAuthenticated, signOut } from '$lib/stores/auth.svelte';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import { goto } from '$app/navigation';
  import ThemeToggle from '$lib/components/ThemeToggle.svelte';

  let { isOpen = false, closeSidebar = () => {} }: { isOpen?: boolean; closeSidebar?: () => void } =
    $props();

  async function handleSignOut() {
    // Close the SSE connection before clearing the auth session so we
    // don't keep the EventSource attached across the logout boundary.
    eventsStream.disconnect();
    await signOut();
    await goto('/login');
  }

  let path = $derived($page.url.pathname);

  function isActive(prefix: string): boolean {
    return path === prefix || path.startsWith(prefix + '/');
  }
</script>

<Sidebar
  {isOpen}
  {closeSidebar}
  backdrop={true}
  position="fixed"
  transitionParams={{ x: -50, duration: 150 }}
  activeUrl={path}
  ariaLabel="Primary navigation"
  classes={{ active: 'bg-neutral-tertiary text-heading', nonactive: 'text-body' }}
  class="z-40"
>
  <SidebarWrapper class="bg-white dark:bg-gray-800 h-full flex flex-col overflow-y-auto">
    <div class="flex items-center justify-between px-4 py-4 border-b border-gray-200 dark:border-gray-700">
      <h2 class="font-semibold text-lg text-gray-900 dark:text-white">
        {m.app_title()}
      </h2>
      <CloseButton onclick={closeSidebar} class="md:hidden" ariaLabel="Close navigation" />
    </div>

    <ul class="space-y-1 p-2">
      <SidebarItem href="/" label={m.sidebar_dashboard()} active={path === '/'} />
    </ul>

    <SidebarGroup class="space-y-1 p-2">
      <li class="px-3 pt-3 pb-1 text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">
        {m.sidebar_status()}
      </li>
      <SidebarItem href="/status/health" label={m.sidebar_status_health()} active={isActive('/status/health')} />
      <SidebarItem href="/status/servers" label={m.sidebar_status_servers()} active={isActive('/status/servers')} />
    </SidebarGroup>

    <ul class="space-y-1 p-2">
      <SidebarItem href="/projects" label={m.sidebar_projects()} active={isActive('/projects')} />
      <SidebarItem href="/servers" label={m.sidebar_servers()} active={isActive('/servers')} />
      <SidebarItem href="/events" label={m.sidebar_events()} active={isActive('/events')} />
    </ul>

    <div class="mt-auto p-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
      <ThemeToggle />
      <button
        class="text-sm text-gray-700 dark:text-gray-200 hover:underline"
        onclick={handleSignOut}
      >
        {m.sidebar_signout()}
      </button>
    </div>
  </SidebarWrapper>
</Sidebar>
