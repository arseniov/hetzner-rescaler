import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      // adapter-node emits a build/ directory runnable via `node build`.
      out: 'build',
      precompress: false,
      envPrefix: ''
    }),
    alias: { $lib: 'src/lib' },
    prerender: {
      handleHttpError: 'warn',
      handleMissingId: 'warn'
    }
  }
};

export default config;