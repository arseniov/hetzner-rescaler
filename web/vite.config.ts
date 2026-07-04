import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { paraglideVitePlugin } from '@inlang/paraglide-js';
import { svelteTesting } from '@testing-library/svelte/vite';

export default defineConfig({
  plugins: [
    paraglideVitePlugin({
      project: './project.inlang',
      outdir: './src/lib/paraglide'
    }),
    svelteTesting(),
    sveltekit()
  ],
  server: {
    port: 5173,
    proxy: {
      // Everything under /api/* (except /api/auth/*) goes to the Go
      // backend on :8080. /api/auth/* is handled by Better Auth in
      // hooks.server.ts — Vite must let SvelteKit handle those routes.
      // We use `bypass` instead of two separate proxy entries so Vite
      // explicitly does NOT proxy Better Auth paths.
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: false,
        bypass: (req: any) => {
          const url = req.url ?? '';
          return url.startsWith('/api/auth/') || url === '/api/auth';
        }
      }
    }
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}'],
    environment: 'jsdom',
    setupFiles: ['./vitest.setup.ts']
  }
});