// Vitest setup: polyfills jsdom is missing for flowbite-svelte's
// Popper initialization. matchMedia is a no-op in jsdom.
if (typeof window !== 'undefined' && !window.matchMedia) {
  // @ts-expect-error - jsdom doesn't ship matchMedia; we polyfill
  window.matchMedia = (query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => {},
    removeListener: () => {},
    addEventListener: () => {},
    removeEventListener: () => {},
    dispatchEvent: () => false
  });
}