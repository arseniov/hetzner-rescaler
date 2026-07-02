import { betterAuth } from 'better-auth';
import { drizzleAdapter } from 'better-auth/adapters/drizzle';
import { db } from './db';
import * as schema from './schema';

const secret = process.env.BETTER_AUTH_SECRET;
if (!secret || secret.length < 32) {
  throw new Error('BETTER_AUTH_SECRET must be set and at least 32 chars');
}

const baseURL = process.env.BETTER_AUTH_URL ?? 'http://localhost:8080';

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
    maxPasswordLength: 128
  },
  session: {
    // Default cookie settings (`httpOnly`, `sameSite=lax`, `secure` when
    // the request is https) are fine for our same-origin setup.
    expiresIn: 60 * 60 * 24 * 7, // 7 days
    updateAge: 60 * 60 * 24 // refresh once per day
  },
  advanced: {
    // Default cookie name `better-auth.session_token` is fine; we never
    // read it client-side.
  }
});

export type Auth = typeof auth;