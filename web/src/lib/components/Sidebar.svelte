<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { m } from '$lib/paraglide/messages.js';
  import { Dialog, DropdownMenu } from 'bits-ui';
  import {
    LayoutDashboard,
    Activity,
    Server,
    FolderKanban,
    Boxes,
    ScrollText,
    LogOut,
    ChevronDown,
    Sun,
    Moon,
    X
  } from 'lucide-svelte';
  import { signOut, currentUser } from '$lib/stores/auth.svelte';
  import { eventsStream } from '$lib/stores/eventsStream.svelte';
  import { theme } from '$lib/stores/theme.svelte';
  import { cn } from '$lib/utils';

  type Props = {
    isOpen?: boolean;
    closeSidebar?: () => void;
  };

  let { isOpen = false, closeSidebar = () => {} }: Props = $props();

  // Shape of the sidebar: one Dashboard item at the top, then a
  // "Status" group, then primary nav items, then auth/utility.
  type NavItem = { href: string; label: string; icon: typeof LayoutDashboard };

  const primary: NavItem[] = [
    { href: '/', label: m.sidebar_dashboard(), icon: LayoutDashboard },
    { href: '/projects', label: m.sidebar_projects(), icon: FolderKanban },
    { href: '/servers', label: m.sidebar_servers(), icon: Boxes },
    { href: '/events', label: m.sidebar_events(), icon: ScrollText }
  ];

  const status: NavItem[] = [
    { href: '/status/health', label: m.sidebar_status_health(), icon: Activity },
    { href: '/status/servers', label: m.sidebar_status_servers(), icon: Server }
  ];

  // Authenticated user is held by the auth store; a small derived
  // letter avatar makes the user menu recognizable without putting a
  // decorative icon above section labels.
  let userEmail = $derived(currentUser()?.email ?? '');
  let userInitial = $derived(userEmail ? userEmail[0]!.toUpperCase() : '·');

  let path = $derived($page.url.pathname);
  function isActive(href: string): boolean {
    return href === '/' ? path === '/' : path === href || path.startsWith(href + '/');
  }

  async function handleSignOut() {
    eventsStream.disconnect();
    await signOut();
    await goto('/login');
  }
</script>

{#snippet nav_body()}
  <div
    class="flex h-full w-64 flex-col border-r border-border bg-card"
    data-testid="sidebar-nav-body"
  >
    <!-- Brand line. Industrial display face, single character-height
         row. No logo, no tagline — the project title is the brand. -->
    <div class="flex items-center border-b border-border px-4 py-4">
      <h1 class="font-display text-base font-semibold tracking-tight text-foreground">
        {m.app_title()}
      </h1>
    </div>

    <!-- Primary nav. -->
    <nav aria-label="Primary" class="flex flex-col gap-0.5 p-2">
      {#each primary as item (item.href)}
        {@const active = isActive(item.href)}
        <a
          href={item.href}
          onclick={closeSidebar}
          class={cn(
            'group flex h-8 items-center gap-3 rounded-sm px-3 text-sm transition-colors',
            active
              ? 'bg-muted text-foreground'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          )}
        >
          <item.icon
            class={cn(
              'size-4 shrink-0',
              active ? 'text-foreground' : 'text-muted-foreground group-hover:text-foreground'
            )}
            strokeWidth={1.5}
            aria-hidden="true"
          />
          <span class="truncate">{item.label}</span>
        </a>
      {/each}
    </nav>

    <!-- Status group. The section label is muted; page labels switch to
         active text only on the currently selected page. -->
    <div class="flex flex-col gap-0.5 p-2 pt-2">
      <h2
        class="px-3 pb-1 text-[11px] font-medium uppercase tracking-wider text-muted-foreground"
      >
        {m.sidebar_status()}
      </h2>
      {#each status as item (item.href)}
        {@const active = isActive(item.href)}
        <a
          href={item.href}
          onclick={closeSidebar}
          class={cn(
            'group flex h-8 items-center gap-3 rounded-sm px-3 text-sm transition-colors',
            active
              ? 'bg-muted text-foreground'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground'
          )}
        >
          <item.icon
            class={cn(
              'size-4 shrink-0',
              active ? 'text-foreground' : 'text-muted-foreground group-hover:text-foreground'
            )}
            strokeWidth={1.5}
            aria-hidden="true"
          />
          <span class="truncate">{item.label}</span>
        </a>
      {/each}
    </div>

    <div class="mt-auto border-t border-border p-2">
      <div class="flex items-center justify-between gap-1 px-2 py-1">
        <!-- Theme toggle: outline-only icon button. The accent color
             never decorates a control — only state. -->
        <button
          type="button"
          onclick={() => theme.toggle()}
          aria-label={theme.current === 'dark' ? m.theme_toggle_light() : m.theme_toggle_dark()}
          class="inline-flex size-8 items-center justify-center rounded-sm text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
        >
          {#if theme.current === 'dark'}
            <Sun class="size-4" strokeWidth={1.5} aria-hidden="true" />
          {:else}
            <Moon class="size-4" strokeWidth={1.5} aria-hidden="true" />
          {/if}
        </button>

        <!-- User menu. Triggered by a circle initial so the affordance
             stays quiet. Uses bits-ui DropdownMenu for keyboard nav. -->
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class="inline-flex h-8 items-center gap-1 rounded-sm px-2 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
            aria-label="Account menu"
          >
            <span
              class="inline-flex size-5 items-center justify-center rounded-full bg-muted font-mono text-[10px] font-medium text-foreground"
            >
              {userInitial}
            </span>
            <ChevronDown class="size-3.5" strokeWidth={1.5} aria-hidden="true" />
          </DropdownMenu.Trigger>
          <DropdownMenu.Portal>
            <DropdownMenu.Content
              side="top"
              align="end"
              sideOffset={6}
              class="z-50 min-w-[12rem] rounded-md border border-border bg-popover p-1 text-popover-foreground"
            >
              <DropdownMenu.Item
                onselect={handleSignOut}
                class="flex h-8 cursor-pointer items-center gap-2 rounded-sm px-2 text-sm text-foreground outline-none data-[highlighted]:bg-muted"
              >
                <LogOut class="size-4 text-muted-foreground" strokeWidth={1.5} aria-hidden="true" />
                {m.sidebar_signout()}
              </DropdownMenu.Item>
            </DropdownMenu.Content>
          </DropdownMenu.Portal>
        </DropdownMenu.Root>
      </div>
    </div>
  </div>
{/snippet}

<!--
  DESKTOP: fixed-position sidebar visible from `md` and up. Hairline
  right border, hairline background contrast. Pointer-events always
  on (no overlay).
-->
<aside class="fixed inset-y-0 left-0 z-30 hidden md:block" aria-label="Primary navigation">
  {@render nav_body()}
</aside>

<!--
  MOBILE: same content rendered as a bits-ui Dialog drawer. The
  Dialog handles focus trap, escape-to-close, scroll lock, and
  ARIA attributes — the only styling we own is positioning.
-->
<Dialog.Root open={isOpen} onOpenChange={(o) => !o && closeSidebar()}>
  <Dialog.Portal>
    <Dialog.Overlay
      class="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm"
    />
    <Dialog.Content
      class="fixed inset-y-0 left-0 z-50 outline-none"
      aria-describedby={undefined}
    >
      {@render nav_body()}
      <button
        type="button"
        onclick={closeSidebar}
        aria-label="Close navigation"
        class="absolute right-3 top-3 inline-flex size-9 items-center justify-center rounded-md bg-card text-muted-foreground hover:text-foreground md:hidden"
      >
        <X class="size-4" strokeWidth={1.5} aria-hidden="true" />
      </button>
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
