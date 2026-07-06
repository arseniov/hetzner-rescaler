// Vitest setup: matchMedia is a no-op in jsdom. The dropdown /
// dialog primitives we use (bits-ui) call matchMedia during mount
// to resolve responsive state, so we polyfill it before the
// component tree initialises.
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