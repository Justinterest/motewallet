# syntax=docker/dockerfile:1

FROM node:22-alpine AS deps

RUN corepack enable && corepack prepare pnpm@9.15.9 --activate

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

FROM node:22-alpine AS builder

RUN corepack enable && corepack prepare pnpm@9.15.9 --activate

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

# 空字符串 = 浏览器同源请求 /api/...，由 nginx 反代到 backend
ARG NEXT_PUBLIC_API_URL=
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL} \
    NEXT_TELEMETRY_DISABLED=1

RUN pnpm build

FROM node:22-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=3001 \
    HOSTNAME=0.0.0.0

RUN addgroup -g 10001 -S nodejs \
    && adduser -S -u 10001 -G nodejs nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3001

CMD ["node", "server.js"]
