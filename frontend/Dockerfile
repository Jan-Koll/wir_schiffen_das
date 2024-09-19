FROM oven/bun:latest AS builder

WORKDIR /app
COPY package*.json .
COPY bun.lockb .
RUN bun install --frozen-lockfile --ci
COPY . .
ENV NODE_ENV=production
RUN bun run build

FROM oven/bun
WORKDIR /app
COPY --from=builder /app/build build/
COPY package.json .
COPY bun.lockb .
RUN bun install --frozen-lockfile --ci --production

EXPOSE 3000
ENTRYPOINT ["bun", "./build/index.js"]

