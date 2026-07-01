<script lang="ts">
  import { signIn, signUp, isAuthenticated } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';
  import Button from '$lib/components/ui/button.svelte';
  import Input from '$lib/components/ui/input.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  let mode = $state<'login' | 'signup'>('login');
  let email = $state('');
  let password = $state('');
  let error = $state<string | null>(null);
  let submitting = $state(false);

  // If already authenticated, bounce to dashboard.
  $effect(() => {
    if (isAuthenticated()) goto('/');
  });

  async function submit(e: SubmitEvent) {
    e.preventDefault();
    error = null;
    submitting = true;
    try {
      if (mode === 'login') await signIn(email, password);
      else await signUp(email, password);
      await goto('/');
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      submitting = false;
    }
  }
</script>

<div class="flex min-h-screen items-center justify-center px-4">
  <form onsubmit={submit} class="w-full max-w-sm space-y-4 rounded-lg border border-border bg-background p-6 shadow">
    <h1 class="text-2xl font-semibold">
      {mode === 'login' ? 'Sign in' : 'Create account'}
    </h1>

    {#if error}
      <Alert variant="destructive">{error}</Alert>
    {/if}

    <label class="block text-sm">
      Email
      <Input type="email" bind:value={email} required autocomplete="email" class="mt-1" />
    </label>

    <label class="block text-sm">
      Password
      <Input type="password" bind:value={password} required minlength={8} autocomplete={mode === 'login' ? 'current-password' : 'new-password'} class="mt-1" />
    </label>

    <Button type="submit" disabled={submitting} class="w-full">
      {submitting ? 'Working…' : (mode === 'login' ? 'Sign in' : 'Create account')}
    </Button>

    <p class="text-center text-xs text-muted-foreground">
      {#if mode === 'login'}
        No account?
        <button type="button" class="underline" onclick={() => (mode = 'signup')}>Create one</button>
      {:else}
        Already have an account?
        <button type="button" class="underline" onclick={() => (mode = 'login')}>Sign in</button>
      {/if}
    </p>
  </form>
</div>
