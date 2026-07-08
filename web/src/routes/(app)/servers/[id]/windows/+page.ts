// The shell layout sets `prerender = true` and the adapter is static with
// `fallback: 'index.html'`. Dynamic [id] routes cannot be prerendered
// without knowing the id at build time, so opt out — they fall through to
// the SPA shell and resolve at runtime via $page.params.id.
export const prerender = false;