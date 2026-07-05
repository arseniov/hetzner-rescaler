// Custom Drizzle migrator using bun:sqlite. Replaces `drizzle-kit
// migrate` (which imports better-sqlite3 internally and cannot run on
// Bun). Reads the generated SQL files from `./drizzle/` and applies
// any un-applied migrations tracked in `__drizzle_migrations`.
//
// Invoked by `bun run db:migrate` (see package.json) and by the Docker
// entrypoint (`CMD ["sh", "-c", "bun run db:migrate && bun build"]`).

import { migrate } from 'drizzle-orm/bun-sqlite/migrator';
import { db } from './db';

migrate(db, { migrationsFolder: './drizzle' });
console.log('Migrations applied successfully.');