import { betterAuth } from 'better-auth';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { db } from './db';
import * as schema from './schema';

const secret = process.env.BETTER_AUTH_SECRET;
if (!secret || secret.length < 32) {
  throw new Error('BETTER_AUTH_SECRET must be set and at least 32 chars');
}

const baseURL = process.env.BETTER_AUTH_URL ?? 'http://localhost:8080';

// DISABLE_SIGNUP defaults to unset (false) so the admin can create
// their account first, then set DISABLE_SIGNUP=true to lock the door.
// Better Auth treats the empty/anything-other-than-true case as
// 'signups allowed' — we read process.env directly so an operator
// can flip it without rebuilding the image.
const signupDisabled = process.env.DISABLE_SIGNUP === 'true';

export const auth = betterAuth({
  secret,
  baseURL,
  // The browser hits the same origin (:8080) regardless of internal
  // topology; Caddy handles the upstream split. Better Auth's CORS
  // checks only need to whitelist the origin the browser sees.
  trustedOrigins: [baseURL],
  database: drizzleAdapter(db, {
    provider: 'sqlite',
    schema: {
      user: schema.user,
      session: schema.session,
      account: schema.account,
      verification: schema.verification
    }
  }),
  emailAndPassword: {
    enabled: true,
    autoSignIn: true,
    minPasswordLength: 8,
    maxPasswordLength: 128,
    // Gate new-account creation on the operator-controlled env var.
    // When true, POST /api/auth/sign-up/email returns 403 and Better
    // Auth's signUp.email client helper rejects with the same error.
    // The UI also hides the toggle (see /login/+page.svelte).
    disableSignUp: signupDisabled
  },
  session: {
    // Default cookie settings (`httpOnly`, `sameSite=lax`, `secure` when
    // the request is https) are fine for our same-origin setup.
    expiresIn: 60 * 60 * 24 * 7, // 7 days
    updateAge: 60 * 60 * 24 // refresh once per day
  }
  // Default cookie name `better-auth.session_token` is fine; we never
  // read it client-side, so we leave Better Auth's advanced cookie defaults
  // in place rather than overriding them here.
});

export type Auth = typeof auth;
