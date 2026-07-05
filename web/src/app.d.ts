// See https://svelte.dev/docs/kit/types#app for the App namespace.

// Locals are populated by web/src/hooks.server.ts. The auth hook attaches
// the validated Better Auth session user + session row to every request;
// load functions can rely on `locals.user` being either a real user or
// `null` (the redirect to /login fires for the latter before the response).
declare global {
  namespace App {
    interface Locals {
      user: import('better-auth').User | null;
      session: import('better-auth').Session | null;
    }
  }
}

export {};
