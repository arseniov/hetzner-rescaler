# ---- Stage 1: install deps & build the SvelteKit Node bundle ----
# We use oven/bun instead of node + npm because:
#   1. The SvelteKit web/ uses bun:sqlite (Bun's built-in SQLite adapter)
#      for runtime DB access — better-sqlite3 is gone entirely.
#   2. bun install --frozen-lockfile is faster and lockfile-deterministic.
#   3. The runner uses `bun build` to execute adapter-node's output;
#      Bun's Node compatibility shim is enough for SvelteKit's runtime.
FROM oven/bun:1 AS builder
WORKDIR /app

# Build-time args the SvelteKit bundle embeds.
ARG PUBLIC_INTERNAL_TOKEN=docker-placeholder
ENV PUBLIC_INTERNAL_TOKEN=$PUBLIC_INTERNAL_TOKEN \
    # Placeholders so the prerender worker can import db.ts (opens
    # DATABASE_URL) and auth.ts (requires BETTER_AUTH_SECRET). The
    # runtime image gets real values from the compose env_file.
    DATABASE_URL=/tmp/build.db \
    BETTER_AUTH_SECRET=build-only-placeholder-secret-must-be-at-least-32-chars

COPY web/package.json ./
COPY web/bun.lock ./
RUN bun install --frozen-lockfile

COPY web/ ./
RUN bun run build

# ---- Stage 2: runtime ----
FROM oven/bun:1-slim AS runner
WORKDIR /app

ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0 \
    DATABASE_URL=/data/db.sqlite

COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/bun.lock ./bun.lock
COPY --from=builder /app/drizzle ./drizzle
COPY --from=builder /app/src ./src
COPY web/drizzle.config.ts ./drizzle.config.ts

EXPOSE 3000

# Run Drizzle migrations on every container start before launching the
# SvelteKit server. The custom migrator (src/lib/server/migrate.ts)
# uses drizzle-orm/bun-sqlite/migrator and is idempotent — it tracks
# applied migrations in `__drizzle_migrations` and skips ones already
# applied, so this is safe to run on every boot, not just first-time
# setup. `&&` ensures the server only starts if migrations succeed (a
# failed migration will exit non-zero and Docker will restart the
# container, surfacing the error in `docker compose logs`).
#
# `bun build/index.js` runs the adapter-node bundle through Bun's
# runtime. We deliberately do NOT use `bun build` (Bun's bundler CLI)
# — that command expects an entrypoint and would error out. `start`
# is mapped to the same command in package.json.
CMD ["sh", "-c", "bun run db:migrate && bun run start"]