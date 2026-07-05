import type { Handle } from '@sveltejs/kit';
import { redirect } from '@sveltejs/kit';
import { building } from '$app/environment';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import { auth } from '$lib/server/auth';

// Paths that bypass the session check entirely. /api/auth/* is owned by
// Better Auth itself (sign-in, sign-up, get-session, callbacks); /_app/*
// is SvelteKit's compiled JS/CSS asset bundle, which the browser fetches
// while loading the redirected /login HTML.
const PASSTHROUGH_PREFIXES = ['/api/auth/', '/_app/'] as const;

function isPassthrough(pathname: string): boolean {
  return PASSTHROUGH_PREFIXES.some((prefix) => pathname.startsWith(prefix));
}

export const handle: Handle = async ({ event, resolve }) => {
  const pathname = event.url.pathname;

  // Server-side auth gate. /api/auth/* and /_app/* always pass through;
  // the prerender phase runs `building === true` and must not see a real
  // session (no request cookies, no DB lookup). Requests with no
  // matching SvelteKit route (e.g. `/favicon.png` or a 404 candidate)
  // also fall through so the response is a clean 404 rather than a
  // misleading 303 to /login.
  if (!building && !isPassthrough(pathname) && event.route.id) {
    const session = await auth.api.getSession({ headers: event.request.headers });
    const hasUser = !!session?.user;

    event.locals.user = hasUser && session ? session.user : null;
    event.locals.session = hasUser && session ? session.session : null;

    const onLogin = pathname === '/login';

    if (hasUser) {
      // Signed-in visitor hitting /login → send them to the dashboard.
      if (onLogin) throw redirect(303, '/');
    } else {
      // Anonymous visitor hitting anything but /login → send them to log in.
      if (!onLogin) throw redirect(303, '/login');
    }
  }

  // /api/auth/* is dispatched to Better Auth; every other path (including
  // /login, the dashboard, /api/* proxies handled by SvelteKit) continues
  // through resolve(). SvelteKit will return prerendered HTML when
  // available and render SSR/load otherwise.
  if (pathname.startsWith('/api/auth/')) {
    return svelteKitHandler({ event, resolve, auth, building });
  }

  return resolve(event);
};
