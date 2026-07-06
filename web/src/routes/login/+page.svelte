<script lang="ts">
  import { goto } from '$app/navigation';
  import { m } from '$lib/paraglide/messages.js';
  import { signIn, signUp } from '$lib/stores/auth.svelte';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Label from '$lib/components/ui/label.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  // signupEnabled is set by +layout.server.ts from the operator-
  // controlled DISABLE_SIGNUP env var. When false, sign-up mode is
  // hidden in the UI AND the server's Better Auth endpoint refuses
  // /api/auth/sign-up/email — defense in depth, since the toggle
  // (mode === 'signup') is also belt-and-braces.
  let { data }: { data: { signupEnabled: boolean } } = $props();

  let mode = $state<'signin' | 'signup'>('signin');
  let email = $state('');
  let password = $state('');
  let error = $state<string | null>(null);
  let busy = $state(false);

  // Derived: the UI should never allow visiting signup mode when
  // signups are disabled, even if some code path tried to flip `mode`.
  let signupBlocked = $derived(!data.signupEnabled);

  // The server-side hook sends logged-in visitors away from /login
  // with a 303 → /. So this component only ever renders for
  // unauthenticated users — no client-side session check at mount time.
  async function submit() {
    error = null;
    if (!email || !password) {
      error = m.login_error_required();
      return;
    }
    busy = true;
    try {
      if (mode === 'signin') await signIn(email, password);
      else await signUp(email, password, email.split('@')[0] ?? 'user');
      await goto('/');
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      busy = false;
    }
  }
</script>

<svelte:head>
  <title>Sign in · Hetzner Rescaler</title>
</svelte:head>

<!--
  Login. Single panel, single job, centred on the canvas. No marketing
  copy above the form — the operator knows what they're signing in to.
-->
<main class="flex min-h-screen items-center justify-center bg-background px-4 py-10">
  <div class="w-full max-w-sm">
    <!-- Brand mark + product name. Display face carries the title. -->
    <header class="mb-8 text-center">
      <h1 class="font-display text-2xl font-semibold tracking-tight text-foreground">
        {m.app_title()}
      </h1>
      <p class="mt-1 font-mono text-xs uppercase tracking-wider text-muted-foreground">
        {mode === 'signin' ? m.login_title() : m.login_signup_title()}
      </p>
    </header>

    <div class="rounded-md border border-border bg-card p-6">
      {#if error}
        <Alert variant="destructive" class="mb-4">{error}</Alert>
      {/if}

      {#if signupBlocked && mode === 'signup'}
        <Alert variant="warning">
          <span class="font-medium text-foreground">{m.login_signup_disabled_title()}</span>
          <span class="mt-1 block text-sm text-muted-foreground">
            {m.login_signup_disabled_message()}
          </span>
        </Alert>
      {:else}
        <form onsubmit={(e) => { e.preventDefault(); submit(); }} class="space-y-4">
          <div class="flex flex-col gap-1.5">
            <Label for="email">{m.login_email_label()}</Label>
            <Input
              id="email"
              type="email"
              bind:value={email}
              required
              autocomplete="email"
              placeholder="ops@example.com"
            />
          </div>
          <div class="flex flex-col gap-1.5">
            <Label for="password">{m.login_password_label()}</Label>
            <Input
              id="password"
              type="password"
              bind:value={password}
              required
              autocomplete="current-password"
            />
          </div>

          <Button type="submit" variant="primary" size="lg" disabled={busy || signupBlocked} class="w-full">
            {#if busy}
              <span class="opacity-70">{m.login_submit()}…</span>
            {:else}
              {mode === 'signin' ? m.login_submit() : m.login_signup_submit()}
            {/if}
          </Button>
        </form>
      {/if}

      {#if !signupBlocked}
        <div class="mt-6 border-t border-border pt-4 text-center">
          <button
            type="button"
            class="font-mono text-xs uppercase tracking-wider text-muted-foreground hover:text-foreground transition-colors"
            onclick={() => { mode = mode === 'signin' ? 'signup' : 'signin'; error = null; }}
          >
            {mode === 'signin' ? m.login_switch_to_signup() : m.login_switch_to_signin()}
          </button>
        </div>
      {/if}
    </div>
  </div>
</main>
