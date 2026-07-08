<script lang="ts">
  import '../app.css';

  let { children } = $props();
</script>

<!--
  Root layout — applies the global stylesheet and renders children.
  Authenticated pages live under `(app)/+layout.svelte`, which owns
  the chrome (header, sidebar, navigation glue). Login and any future
  unauthenticated pages render straight through here.

  Splitting the chrome into the `(app)` group avoids a previous bug
  where `{#if isAuthenticated()}` racing navigation events made the
  header flicker off and on. Since the server-side hook redirects
  unauthenticated visitors to `/login` before this layout's children
  render, the authed layout is structurally always authed.
-->
{@render children?.()}
