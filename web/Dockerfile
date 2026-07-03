# ---- Stage 1: install deps & build the SvelteKit Node bundle ----
# We use Node 22 (not Bun) because better-sqlite3's native module is
# not supported on Bun's runtime (oven-sh/bun#4290), and SvelteKit's
# prerender worker imports db.ts which loads better-sqlite3.
# npm install does not read web/bun.lock; it resolves from package.json
# directly. bun.lock stays in the repo for local dev tooling.
FROM node:22-bookworm-slim AS builder
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
# --legacy-peer-deps: better-auth@1.6.23 declares peerOptional
# better-sqlite3@^12, but the project pins ^11.5.0. npm 7+ rejects
# the peerOptional mismatch; this flag restores the pre-v7 lenient
# algorithm and matches the version bun installed in dev.
RUN npm install --legacy-peer-deps

COPY web/ ./
RUN npm run build

# ---- Stage 2: runtime ----
FROM node:22-bookworm-slim AS runner
WORKDIR /app

ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0 \
    DATABASE_URL=/data/db.sqlite

COPY --from=builder /app/build ./build
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/drizzle ./drizzle

EXPOSE 3000

# node build/ runs the adapter-node server. SvelteKit honours PORT/HOST.
CMD ["node", "build"]