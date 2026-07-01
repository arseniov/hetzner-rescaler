import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const buildDir = path.resolve(__dirname, '..', 'internal', 'web', 'build');

const gitignorePath = path.join(buildDir, '.gitignore');
const gitkeepPath = path.join(buildDir, '.gitkeep');
const gitignoreContent = '*\n!.gitignore\n!.gitkeep\n';

// adapter-static calls rimraf(assets) and rimraf(pages) before writing,
// which deletes the tracked .gitignore and .gitkeep in internal/web/build/.
// //go:embed in the Go side points at that directory, so the markers must
// be re-created after the adapter finishes — package.json scripts that
// touch the markers BEFORE the build get undone by the adapter.
const restoreMarkers = {
  name: 'restore-embed-markers',
  apply: 'build' as const,
  closeBundle() {
    fs.mkdirSync(buildDir, { recursive: true });
    fs.writeFileSync(gitignorePath, gitignoreContent);
    fs.writeFileSync(gitkeepPath, '');
  }
};

export default defineConfig({
  plugins: [sveltekit(), restoreMarkers],
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