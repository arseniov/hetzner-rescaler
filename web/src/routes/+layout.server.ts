// Root server load: exposes operator-controlled feature flags to every
// page on the server side, so the UI can react without inlining env
// values into the client bundle (PUBLIC_* vars require src/env.ts
// declaration or $env/dynamic/public plumbing).
//
// Currently exposed:
//   signupEnabled: boolean — true when new-account creation is allowed.
//     `login/+page.svelte` reads this to hide the "switch to sign up"
//     toggle and, if the user somehow lands in signup mode (e.g. local
//     state from a previous session), show a "signups disabled" Alert
//     instead of the form. The value is driven by Better Auth's
//     `disableSignUp` config (which itself reads DISABLE_SIGNUP) —
//     re-reading here would drift if we forgot to keep them in sync,
//     so we read the env directly and trust the operator set it
//     consistently.
export const load = async () => {
  const signupEnabled = process.env.DISABLE_SIGNUP !== 'true';
  return { signupEnabled };
};