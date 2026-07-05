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
  // `bun:sqlite` (and any future `bun:*` builtin) is provided by the
  // Bun runtime at load time — it has no npm package to bundle. Without
  // the regex below, Rollup aborts with "Rollup failed to resolve
  // import 'bun:sqlite'". `optimizeDeps.exclude` keeps `vite dev` from
  // trying to pre-bundle it; `build.rollupOptions.external` keeps the
  // production SSR bundle from trying to bundle it. Both are needed
  // because Vite uses different resolution paths for dev (esbuild via
  // optimizeDeps) vs build (Rollup).
  optimizeDeps: {
    exclude: ['bun:sqlite', 'bun:sql', 'bun:fs']
  },
  build: {
    rollupOptions: {
      external: [/^bun:/]
    }
  },
  ssr: {
    // `db.ts` is only ever imported from server code (hooks.server.ts
    // → auth.ts → db.ts), so externalizing it for SSR is safe and
    // prevents Rollup from trying to inline the better-auth / drizzle
    // server graph into the prerender worker.
    external: ['bun:sqlite']
  },
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