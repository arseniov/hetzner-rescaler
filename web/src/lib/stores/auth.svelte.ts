import { authClient } from '$lib/auth-client';

export interface SessionUser {
  id: string;
  email: string;
  name: string;
}

// Holds the active session, hydrated by `ensureSession()` and cleared
// by `signOut()`. Better Auth itself uses cookies — this is a cache
// for the UI so components don't refetch on every read.
let _user = $state<SessionUser | null>(null);

export async function ensureSession(): Promise<SessionUser | null> {
  const { data } = await authClient.getSession();
  if (data?.user) {
    _user = { id: data.user.id, email: data.user.email, name: data.user.name };
  } else {
    _user = null;
  }
  return _user;
}

export function currentUser(): SessionUser | null {
  return _user;
}

export function isAuthenticated(): boolean {
  return _user !== null;
}

export async function signIn(email: string, password: string): Promise<void> {
  const { data, error } = await authClient.signIn.email({ email, password });
  if (error) throw new Error(error.message ?? 'sign-in failed');
  if (!data?.user) throw new Error('sign-in returned no user');
  _user = { id: data.user.id, email: data.user.email, name: data.user.name };
}

export async function signUp(email: string, password: string, name: string): Promise<void> {
  const { data, error } = await authClient.signUp.email({ email, password, name });
  if (error) throw new Error(error.message ?? 'sign-up failed');
  if (!data?.user) throw new Error('sign-up returned no user');
  _user = { id: data.user.id, email: data.user.email, name: data.user.name };
}

export async function signOut(): Promise<void> {
  try {
    await authClient.signOut();
  } finally {
    _user = null;
  }
}
