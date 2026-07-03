import { createAuthClient } from 'better-auth/svelte';

// Empty client config is intentional: BETTER_AUTH_URL and trusted origins
// are configured server-side in src/lib/server/auth.ts, and routing is
// mounted in src/hooks.server.ts via svelteKitHandler. The client just
// posts to whatever /api/auth/* SvelteKit dispatches.
export const authClient = createAuthClient({});

export type AuthClient = typeof authClient;
