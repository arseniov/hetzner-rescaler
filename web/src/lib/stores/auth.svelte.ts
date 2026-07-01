import { auth, type AuthSession } from '$lib/authorizer';

const STORAGE_KEY = 'rescaler.auth';

let _session = $state<AuthSession | null>(loadInitial());

function loadInitial(): AuthSession | null {
  if (typeof localStorage === 'undefined') return null;
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return null;
  try {
    const s = JSON.parse(raw) as AuthSession;
    if (s.expiresAt < Date.now()) {
      localStorage.removeItem(STORAGE_KEY);
      return null;
    }
    return s;
  } catch {
    return null;
  }
}

export function session(): AuthSession | null {
  return _session;
}

export function isAuthenticated(): boolean {
  return _session !== null && _session.expiresAt > Date.now();
}

export async function signIn(email: string, password: string): Promise<void> {
  const s = await auth.login(email, password);
  _session = s;
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
  }
}

export async function signUp(email: string, password: string): Promise<void> {
  const s = await auth.signup(email, password);
  _session = s;
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
  }
}

export async function signOut(): Promise<void> {
  if (_session) {
    try {
      await auth.logout(_session.token);
    } catch {
      /* ignore */
    }
  }
  _session = null;
  if (typeof localStorage !== 'undefined') {
    localStorage.removeItem(STORAGE_KEY);
  }
}
