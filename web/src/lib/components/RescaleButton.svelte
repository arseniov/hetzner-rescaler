<script lang="ts">
  import { api, ApiError } from '$lib/api';
  import Button from '$lib/components/ui/button.svelte';
  import Dialog from '$lib/components/ui/dialog.svelte';
  import Alert from '$lib/components/ui/alert.svelte';

  interface Props {
    serverId: number;
    direction: 'up' | 'down';
    label: string;
    confirmCopy?: string;
    onComplete?: () => void;
  }

  let { serverId, direction, label, confirmCopy, onComplete }: Props = $props();
  let open = $state(false);
  let submitting = $state(false);
  let error = $state<string | null>(null);

  async function submit() {
    submitting = true;
    error = null;
    try {
      await api.rescale(serverId, { direction, confirm: true });
      open = false;
      onComplete?.();
    } catch (err) {
      error = err instanceof ApiError ? err.message : String(err);
    } finally {
      submitting = false;
    }
  }
</script>

<Button variant={direction === 'down' ? 'outline' : 'default'} onclick={() => (open = true)}>
  {label}
</Button>

<Dialog
  bind:open
  title="Confirm rescale"
  description={confirmCopy ?? `This will rescale the server ${direction === 'up' ? 'up' : 'down'}. The action may take a few minutes.`}
>
  {#if error}<Alert variant="destructive">{error}</Alert>{/if}
  {#snippet footer()}
    <Button variant="outline" onclick={() => (open = false)} disabled={submitting}>Cancel</Button>
    <Button variant="destructive" onclick={submit} disabled={submitting}>
      {submitting ? 'Rescaling…' : `Rescale ${direction}`}
    </Button>
  {/snippet}
</Dialog>