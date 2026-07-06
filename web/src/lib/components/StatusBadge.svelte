<script lang="ts">
  // Status union covers two vocabularies:
  //   - Event-derived fallback (ok / degraded / unknown) — used when
  //     the API didn't return a live Hetzner status.
  //   - Hetzner Cloud live states (running / initializing / starting /
  //     stopping / off / deleting) — populated from `server.status`
  //     when the API call to Hetzner succeeded.
  // Both render with the same vocabulary (dot + ALL-CAPS label) so a
  // row looks uniform; the colour encodes health, the label is the
  // exact source word so the operator can grep logs.
  export type Status =
    | 'ok' | 'degraded' | 'unknown'
    | 'running' | 'initializing' | 'starting'
    | 'stopping' | 'off' | 'deleting';

  interface Props { status: Status }
  let { status }: Props = $props();

  // Compact status indicator — a single dot in the relevant token
  // colour, followed by an ALL-CAPS label. No filled background: this
  // is a row in a data table, not a callout.
  const tone: Record<Status, string> = {
    ok: 'bg-success',
    degraded: 'bg-warning',
    unknown: 'bg-muted-foreground',
    running: 'bg-success',
    initializing: 'bg-warning',
    starting: 'bg-warning',
    stopping: 'bg-warning',
    off: 'bg-muted-foreground',
    deleting: 'bg-destructive'
  };
  // Label is the raw source word (uppercased). We don't try to
  // humanise Hetzner states — operators want the verbatim token.
  const text: Record<Status, string> = {
    ok: 'OK',
    degraded: 'DEGRADED',
    unknown: 'UNKNOWN',
    running: 'RUNNING',
    initializing: 'INITIALIZING',
    starting: 'STARTING',
    stopping: 'STOPPING',
    off: 'OFF',
    deleting: 'DELETING'
  };
</script>

<span class="inline-flex items-center gap-2 font-mono text-xs font-medium uppercase tracking-wider text-foreground">
  <span class="inline-block size-1.5 rounded-full {tone[status]}"></span>
  {text[status]}
</span>
