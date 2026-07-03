import type { Handle } from '@sveltejs/kit';
import { svelteKitHandler } from 'better-auth/svelte-kit';
import { auth } from '$lib/server/auth';

// Better Auth handles every /api/auth/* request — sign-up, sign-in,
// sign-out, session, callbacks, etc. Everything else falls through to
// SvelteKit's normal request flow (page rendering, +server.ts routes,
// or proxying via vite.config during dev).
export const handle: Handle = async ({ event, resolve }) => {
  return svelteKitHandler({ event, resolve, auth, building: false });
};
