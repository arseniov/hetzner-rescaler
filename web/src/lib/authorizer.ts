import { Authorizer } from '@authorizerdev/authorizer-js';

const PUBLIC_AUTHORIZER_URL =
  (import.meta.env.PUBLIC_AUTHORIZER_URL as string | undefined) ||
  'http://localhost:8080';

let _client: Authorizer | null = null;
function client(): Authorizer {
  if (!_client) {
    _client = new Authorizer({
      authorizerURL: PUBLIC_AUTHORIZER_URL,
      redirectURL: window.location.origin
    });
  }
  return _client;
}

// authorizer-js declares its own `Headers` type (Record<string, string>),
// distinct from the DOM Headers class. Use a plain record.
type AuthHeaders = Record<string, string>;

function authHeaders(token: string): AuthHeaders {
  return { Authorization: `Bearer ${token}` };
}

export interface AuthSession {
  token: string;
  expiresAt: number;
}

export const auth = {
  async login(email: string, password: string): Promise<AuthSession> {
    const res = await client().login({ email, password });
    if (!res?.data?.access_token) {
      throw new Error(res?.errors?.[0]?.message ?? 'login failed');
    }
    return {
      token: res.data.access_token,
      expiresAt: Date.now() + 60 * 60 * 1000
    };
  },

  async signup(email: string, password: string): Promise<AuthSession> {
    const res = await client().signup({
      email,
      password,
      confirm_password: password
    });
    if (!res?.data?.access_token) {
      throw new Error(res?.errors?.[0]?.message ?? 'signup failed');
    }
    return {
      token: res.data.access_token,
      expiresAt: Date.now() + 60 * 60 * 1000
    };
  },

  async logout(token: string): Promise<void> {
    try {
      await client().logout(authHeaders(token));
    } catch {
      // ignore logout errors — local session is cleared regardless
    }
  },

  async me(token: string): Promise<{ email: string } | null> {
    try {
      const res = await client().getSession(authHeaders(token));
      return res?.data?.user?.email ? { email: res.data.user.email } : null;
    } catch {
      return null;
    }
  }
};
