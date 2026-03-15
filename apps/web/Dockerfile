FROM node:20-alpine AS base
RUN npm install -g pnpm@9

FROM base AS deps
WORKDIR /app
# Copy workspace files
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml .npmrc ./
COPY packages/core/package.json ./packages/core/
COPY apps/web/package.json ./apps/web/
RUN pnpm install --frozen-lockfile

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/packages/core/node_modules ./packages/core/node_modules
COPY --from=deps /app/apps/web/node_modules ./apps/web/node_modules
COPY packages/core ./packages/core
COPY apps/web ./apps/web
COPY package.json pnpm-workspace.yaml .npmrc ./
RUN pnpm --filter @timetracking/core build
RUN pnpm --filter @timetracking/web build

FROM base AS runner
WORKDIR /app
ENV NODE_ENV=production
# Copy deps and build output
COPY --from=deps /app/node_modules ./node_modules
COPY --from=deps /app/packages/core/node_modules ./packages/core/node_modules
COPY --from=deps /app/apps/web/node_modules ./apps/web/node_modules
COPY --from=builder /app/packages/core/dist ./packages/core/dist
COPY --from=builder /app/apps/web/.next ./apps/web/.next
COPY package.json pnpm-workspace.yaml .npmrc ./
COPY packages/core/package.json ./packages/core/
COPY apps/web/package.json ./apps/web/
EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"
WORKDIR /app/apps/web
CMD ["node_modules/.bin/next", "start"]
