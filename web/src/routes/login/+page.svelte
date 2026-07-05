<script lang="ts">
  import { goto } from '$app/navigation';
  import { Button, Card, Input, Label, Alert } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { signIn, signUp } from '$lib/stores/auth.svelte';

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
  // The form is also unconditionally hidden in that state below.
  let signupBlocked = $derived(!data.signupEnabled);

  // The server-side hook (web/src/hooks.server.ts) sends logged-in
  // visitors away from /login with a 303 → /. So this component only
  // ever renders for unauthenticated users — no client-side session
  // check is needed at mount time.
  async function submit() {
    error = null;
    if (!email || !password) {
      error = m.login_error_required();
      return;
    }
    busy = true;
    try {
      if (mode === 'signin') await signIn(email, password);
      else await signUp(email, password, email.split('@')[0]);
      await goto('/');
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      busy = false;
    }
  }
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 dark:bg-gray-900 p-4">
  <Card class="w-full max-w-md border-0">
    <h2 class="text-2xl font-semibold mb-4 text-gray-900 dark:text-white">
      {mode === 'signin' ? m.login_title() : m.login_signup_title()}
    </h2>

    {#if error}
      <Alert color="red" class="mb-3">{error}</Alert>
    {/if}

    {#if signupBlocked && mode === 'signup'}
      <Alert color="yellow" class="mb-3">
        <span class="font-medium">{m.login_signup_disabled_title()}</span>
        <span class="block mt-1 text-sm">{m.login_signup_disabled_message()}</span>
      </Alert>
    {:else}
      <form onsubmit={(e) => { e.preventDefault(); submit(); }} class="space-y-3">
        <Label>
          {m.login_email_label()}
          <Input type="email" bind:value={email} required autocomplete="email" />
        </Label>
        <Label>
          {m.login_password_label()}
          <Input type="password" bind:value={password} required autocomplete="current-password" />
        </Label>

        <Button type="submit" disabled={busy || signupBlocked} class="w-full">
          {mode === 'signin' ? m.login_submit() : m.login_signup_submit()}
        </Button>
      </form>
    {/if}

    {#if !signupBlocked}
      <button
        type="button"
        class="mt-4 text-sm text-blue-600 dark:text-blue-400 hover:underline"
        onclick={() => { mode = mode === 'signin' ? 'signup' : 'signin'; error = null; }}
      >
        {mode === 'signin' ? m.login_switch_to_signup() : m.login_switch_to_signin()}
      </button>
    {/if}
  </Card>
</div>