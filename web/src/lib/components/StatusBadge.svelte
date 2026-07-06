<script lang="ts">
  export type Status = 'ok' | 'degraded' | 'unknown';
  interface Props { status: Status }
  let { status }: Props = $props();

  // Compact status indicator — a single dot in the relevant token
  // colour, followed by an ALL-CAPS label. No flowbite Badge, no
  // filled background: this is a row in a data table, not a callout.
  const tone: Record<Status, string> = {
    ok: 'bg-success',
    degraded: 'bg-warning',
    unknown: 'bg-muted-foreground'
  };
  const text: Record<Status, string> = {
    ok: 'OK',
    degraded: 'DEGRADED',
    unknown: 'UNKNOWN'
  };
</script>

<span class="inline-flex items-center gap-2 font-mono text-xs font-medium uppercase tracking-wider text-foreground">
  <span class="inline-block size-1.5 rounded-full {tone[status]}"></span>
  {text[status]}
</span>
