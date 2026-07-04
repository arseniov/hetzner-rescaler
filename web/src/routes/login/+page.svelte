<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Button, Card, Input, Label, Alert } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { ensureSession, isAuthenticated, signIn, signUp } from '$lib/stores/auth.svelte';

  let mode = $state<'signin' | 'signup'>('signin');
  let email = $state('');
  let password = $state('');
  let error = $state<string | null>(null);
  let busy = $state(false);

  onMount(async () => {
    await ensureSession();
    if (isAuthenticated()) await goto('/');
  });

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
  <Card class="w-full max-w-md">
    <h2 class="text-2xl font-semibold mb-4 text-gray-900 dark:text-white">
      {mode === 'signin' ? m.login_title() : m.login_signup_title()}
    </h2>

    {#if error}
      <Alert color="red" class="mb-3">{error}</Alert>
    {/if}

    <form onsubmit={(e) => { e.preventDefault(); submit(); }} class="space-y-3">
      <Label>
        {m.login_email_label()}
        <Input type="email" bind:value={email} required autocomplete="email" />
      </Label>
      <Label>
        {m.login_password_label()}
        <Input type="password" bind:value={password} required autocomplete="current-password" />
      </Label>

      <Button type="submit" disabled={busy} class="w-full">
        {mode === 'signin' ? m.login_submit() : m.login_signup_submit()}
      </Button>
    </form>

    <button
      type="button"
      class="mt-4 text-sm text-blue-600 dark:text-blue-400 hover:underline"
      onclick={() => { mode = mode === 'signin' ? 'signup' : 'signin'; error = null; }}
    >
      {mode === 'signin' ? m.login_switch_to_signup() : m.login_switch_to_signin()}
    </button>
  </Card>
</div>