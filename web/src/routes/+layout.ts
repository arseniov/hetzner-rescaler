// SPA: no SSR (the browser hydrates the shell client-side and fetches
// data via apiFetch), but `prerender` is OFF so every request flows
// through SvelteKit's SSR pipeline — which is the only path that
// actually runs `hooks.server.ts`. With `prerender = true` the
// adapter-node bundle serves prerendered HTML directly from disk via
// a polka middleware that sits BEFORE the hook, so the auth gate
// would never fire. Forcing `prerender = false` adds one synchronous
// DB read per request (the session lookup) but is the price of
// server-side auth enforcement.
export const ssr = false;
export const prerender = false;
