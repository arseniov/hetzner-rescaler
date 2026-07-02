import Database from 'better-sqlite3';
import { drizzle } from 'drizzle-orm/better-sqlite3';
import * as schema from './schema';

const url = process.env.DATABASE_URL ?? '/data/db.sqlite';

// better-sqlite3 connects synchronously. WAL mode lets the Go side
// (`modernc.org/sqlite` opens the same file with WAL+foreign_keys) read
// while Better Auth writes; SQLite serializes writers fine at our
// concurrency (a handful of sign-up / sign-in events).
const sqlite = new Database(url);
sqlite.pragma('journal_mode = WAL');
sqlite.pragma('foreign_keys = ON');

export const db = drizzle(sqlite, { schema });