// bun:sqlite is Bun's built-in SQLite adapter — no native binding, no
// `node-gyp` build step, ~3-6x faster than better-sqlite3 per the Bun
// docs. We open the same `/data/db.sqlite` file the Go backend writes
// (`modernc.org/sqlite`), so both processes share the data via WAL.
import { Database } from 'bun:sqlite';
import { drizzle } from 'drizzle-orm/bun-sqlite';
import { mkdirSync } from 'node:fs';
import { dirname } from 'node:path';
import * as schema from './schema';

// Production (compose) overrides this with `/data/db.sqlite`. The
// fallback to `/tmp/build.db` makes local `bun run build` work without
// extra setup — SvelteKit's prerender worker imports this module
// (hooks.server.ts → auth.ts → db.ts) and would crash trying to open
// a file under a non-existent `/data/` directory.
const url = process.env.DATABASE_URL ?? '/tmp/build.db';

// Ensure the parent directory exists. /tmp always exists; /data is
// created by the compose volume mount. This is a defensive no-op in
// production but lets users point DATABASE_URL at any writable path.
const dir = dirname(url);
if (dir && dir !== '/' && dir !== '.') {
  try {
    mkdirSync(dir, { recursive: true });
  } catch {
    // ignore — bun:sqlite will surface a clearer error if the path is
    // genuinely unwritable (e.g. permission denied on /data/).
  }
}

// bun:sqlite opens synchronously (no callback API like `node-sqlite3`).
// `create: true` makes the file if it's missing — important for the
// Docker entrypoint where the volume may be empty on first boot.
const sqlite = new Database(url, { create: true });

// bun:sqlite has no `db.pragma(...)` helper (better-sqlite3 had one).
// Use `db.run("PRAGMA ...")` directly. WAL lets the Go backend read
// concurrently with Better Auth writes; foreign_keys is on by default
// in bun:sqlite but we set it explicitly for parity with the Go side
// (which enables it via `_pragma=foreign_keys(1)` in modernc.org/sqlite).
sqlite.run('PRAGMA journal_mode = WAL;');
sqlite.run('PRAGMA foreign_keys = ON;');

export const db = drizzle({ client: sqlite, schema });