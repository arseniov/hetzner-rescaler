import { createAuthClient } from 'better-auth/svelte';

// The baseURL here matches BETTER_AUTH_URL — the origin the browser sees
// after the reverse proxy. The client posts to /api/auth/* relative
// paths, which SvelteKit's hooks.server.ts dispatches into the Better
// Auth handler.
export const authClient = createAuthClient({});

export type AuthClient = typeof authClient;