import { defineConfig } from 'drizzle-kit';

const url = process.env.DATABASE_URL ?? '/data/db.sqlite';

export default defineConfig({
  schema: './src/lib/server/schema.ts',
  out: './drizzle',
  dialect: 'sqlite',
  // Use the bun-sqlite driver so drizzle-kit introspects via bun:sqlite
  // rather than reaching for better-sqlite3 (which Bun cannot load).
  // The CLI is invoked via `bun x drizzle-kit generate` from package.json.
  driver: 'bun-sqlite',
  dbCredentials: { url },
  verbose: true,
  strict: true
});