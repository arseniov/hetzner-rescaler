import { describe, it, expect, vi, beforeEach } from 'vitest';

// `vi.mock(...)` factory bodies are hoisted above all imports, so any
// variable they reference must ALSO be hoisted via `vi.hoisted`.
// Otherwise the factory runs in a temporal-dead-zone and emits a
// "Cannot access … before initialization" error at module-load time.
const mocks = vi.hoisted(() => {
  const getSession = vi.fn();
  const svelteKitHandler = vi.fn();
  const redirect = vi.fn((status: number, location: string) => {
    const err = new Error(`REDIRECT:${status}:${location}`) as Error & {
      status?: number;
      location?: string;
    };
    err.status = status;
    err.location = location;
    throw err;
  });
  return { getSession, svelteKitHandler, redirect };
});

vi.mock('$app/environment', () => ({ building: false }));
vi.mock('$lib/server/auth', () => ({
  auth: { api: { getSession: mocks.getSession } }
}));
vi.mock('better-auth/svelte-kit', () => ({
  svelteKitHandler: mocks.svelteKitHandler
}));
vi.mock('@sveltejs/kit', async () => {
  const actual = await vi.importActual<typeof import('@sveltejs/kit')>('@sveltejs/kit');
  return {
    ...actual,
    redirect: (...args: unknown[]) => mocks.redirect(...(args as [number, string]))
  };
});

// Now safe to import the hook — every external it touches is stubbed.
import { handle } from './hooks.server';

type Event = Parameters<typeof handle>[0]['event'];
type Resolve = Parameters<typeof handle>[0]['resolve'];

function makeEvent(pathname: string, routeId: string | null = pathname): Event {
  const url = new URL(`http://localhost${pathname}`);
  const headers = new Headers();
  return {
    url,
    request: new Request(url, { headers }),
    locals: {} as App.Locals,
    // `event.route.id` is what the hook uses to recognise a SvelteKit
    // page route vs a static asset. Tests default to a real route id
    // (matching the pathname); static-asset tests pass `null`.
    route: { id: routeId } as { id: string | null },
    params: {},
    isDataRequest: false,
    isSubRequest: false,
    isRemoteRequest: false,
    cookies: undefined as never,
    fetch: undefined as never,
    getClientAddress: () => '127.0.0.1',
    platform: undefined,
    setHeaders: () => {}
  } as unknown as Event;
}

function makeResolve(): Resolve {
  return vi.fn(async () => new Response('resolved', { status: 200 })) as unknown as Resolve;
}

beforeEach(() => {
  mocks.getSession.mockReset();
  mocks.svelteKitHandler.mockReset();
  mocks.redirect.mockClear();
});

describe('hooks.server.handle', () => {
  it('redirects anonymous users on / to /login (303)', async () => {
    mocks.getSession.mockResolvedValue(null);
    const event = makeEvent('/');
    const resolve = makeResolve();

    await expect(handle({ event, resolve })).rejects.toThrow(/REDIRECT:303:\/login/);
    expect(resolve).not.toHaveBeenCalled();
    expect(mocks.getSession).toHaveBeenCalledTimes(1);
    expect((event.locals as App.Locals).user).toBeNull();
    expect((event.locals as App.Locals).session).toBeNull();
  });

  it('redirects anonymous users on /projects to /login', async () => {
    mocks.getSession.mockResolvedValue(null);
    const event = makeEvent('/projects');
    const resolve = makeResolve();

    await expect(handle({ event, resolve })).rejects.toThrow(/REDIRECT:303:\/login/);
  });

  it('lets anonymous users stay on /login (no redirect)', async () => {
    mocks.getSession.mockResolvedValue(null);
    const event = makeEvent('/login');
    const resolve = makeResolve();

    const response = await handle({ event, resolve });
    expect(response).toBeInstanceOf(Response);
    expect(resolve).toHaveBeenCalledTimes(1);
    expect(mocks.svelteKitHandler).not.toHaveBeenCalled();
  });

  it('redirects logged-in users on /login to /', async () => {
    mocks.getSession.mockResolvedValue({
      user: { id: 'u1', email: 'a@b.c', name: 'a' },
      session: { id: 's1', userId: 'u1' }
    });
    const event = makeEvent('/login');
    const resolve = makeResolve();

    await expect(handle({ event, resolve })).rejects.toThrow(/REDIRECT:303:\/$/);
  });

  it('passes logged-in users through on / and attaches user to locals', async () => {
    const user = { id: 'u1', email: 'a@b.c', name: 'a' };
    const session = { id: 's1', userId: 'u1' };
    mocks.getSession.mockResolvedValue({ user, session });
    const event = makeEvent('/');
    const resolve = makeResolve();

    const response = await handle({ event, resolve });
    expect(response).toBeInstanceOf(Response);
    expect(resolve).toHaveBeenCalledTimes(1);
    expect((event.locals as App.Locals).user).toEqual(user);
    expect((event.locals as App.Locals).session).toEqual(session);
  });

  it('always passes /api/auth/* through without consulting the session', async () => {
    const event = makeEvent('/api/auth/sign-in');
    const resolve = makeResolve();
    mocks.svelteKitHandler.mockImplementation(async () => new Response('auth-handler', { status: 200 }));

    const response = await handle({ event, resolve });
    expect(response).toBeInstanceOf(Response);
    expect(mocks.svelteKitHandler).toHaveBeenCalledTimes(1);
    expect(mocks.getSession).not.toHaveBeenCalled();
    expect(resolve).not.toHaveBeenCalled();
  });

  it('always passes /_app/* built assets through without consulting the session', async () => {
    mocks.getSession.mockResolvedValue(null);
    const event = makeEvent('/_app/immutable/entry/start.AAA.js');
    const resolve = makeResolve();

    const response = await handle({ event, resolve });
    expect(response).toBeInstanceOf(Response);
    expect(resolve).toHaveBeenCalledTimes(1);
    expect(mocks.getSession).not.toHaveBeenCalled();
  });

  it('passes non-route paths through (avoids redirecting missing static assets to /login)', async () => {
    mocks.getSession.mockResolvedValue(null);
    const event = makeEvent('/favicon.png', null);
    const resolve = makeResolve();

    const response = await handle({ event, resolve });
    expect(response).toBeInstanceOf(Response);
    expect(resolve).toHaveBeenCalledTimes(1);
    expect(mocks.getSession).not.toHaveBeenCalled();
  });

  it('forwards the request headers to auth.api.getSession', async () => {
    mocks.getSession.mockResolvedValue(null);
    const event = makeEvent('/');
    (event.request as Request).headers.set('cookie', 'better-auth.session_token=abc');
    const resolve = makeResolve();

    await expect(handle({ event, resolve })).rejects.toThrow(/REDIRECT:303:\/login/);
    expect(mocks.getSession).toHaveBeenCalledWith({ headers: event.request.headers });
  });
});
