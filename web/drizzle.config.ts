import type { Config } from 'drizzle-kit';

const url = process.env.DATABASE_URL ?? '/data/db.sqlite';

export default {
  schema: './src/lib/server/schema.ts',
  out: './drizzle',
  dialect: 'sqlite',
  dbCredentials: { url },
  verbose: true,
  strict: true
} satisfies Config;