<script lang="ts">
  import { page } from '$app/stores';
  import { Sidebar, SidebarGroup, SidebarItem, SidebarWrapper } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { isAuthenticated, signOut } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';

  async function handleSignOut() {
    await signOut();
    await goto('/login');
  }

  let path = $derived($page.url.pathname);

  function isActive(prefix: string): boolean {
    return path === prefix || path.startsWith(prefix + '/');
  }
</script>

<Sidebar
  isOpen={true}
  disableBreakpoints={true}
  alwaysOpen={true}
  activateClickOutside={false}
  backdrop={false}
  position="fixed"
  ariaLabel="Primary navigation"
  class="h-screen w-64 border-r border-gray-200 dark:border-gray-700"
>
  <SidebarWrapper class="bg-white dark:bg-gray-800 h-full flex flex-col overflow-y-auto">
    <h2 class="px-4 py-4 font-semibold text-lg text-gray-900 dark:text-white border-b border-gray-200 dark:border-gray-700">
      {m.app_title()}
    </h2>

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

    <div class="mt-auto p-4 border-t border-gray-200 dark:border-gray-700">
      {#if isAuthenticated()}
        <button
          type="button"
          class="text-sm text-gray-700 hover:underline dark:text-gray-200"
          onclick={handleSignOut}
        >
          {m.sidebar_signout()}
        </button>
      {/if}
    </div>
  </SidebarWrapper>
</Sidebar>