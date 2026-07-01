import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    port: 5173,
    proxy: {
      // During `bun run dev`, proxy /api/* to the Go backend on :8080.
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: false
      }
    }
  },
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}']
  }
});