<script lang="ts">
  import { Badge } from 'flowbite-svelte';

  export type Status = 'ok' | 'degraded' | 'unknown';
  interface Props { status: Status }
  let { status }: Props = $props();

  const color: Record<Status, 'green' | 'yellow' | 'gray'> = {
    ok: 'green',
    degraded: 'yellow',
    unknown: 'gray'
  };
  const text: Record<Status, string> = {
    ok: 'OK',
    degraded: 'DEGRADED',
    unknown: 'UNKNOWN'
  };
</script>

<Badge color={color[status]}>{text[status]}</Badge>