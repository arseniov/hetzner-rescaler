// Dynamic route — cannot be prerendered without a known id at build time.
// Falls back to the SPA shell (configured in svelte.config.js) and renders
// on the client based on $page.params.id.
export const prerender = false;
