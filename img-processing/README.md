# Bun + Hono + PostgreSQL

```sh
docker compose up -d
Copy-Item .env.example .env
bun install
bun run db:migrate
bun run dev
```

Open `http://localhost:3000`.
