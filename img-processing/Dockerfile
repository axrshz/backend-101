FROM oven/bun:1.3.14

WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY drizzle ./drizzle
COPY drizzle.config.ts ./
COPY public ./public
COPY src ./src

USER bun

EXPOSE 3000

CMD ["bun", "run", "start"]
