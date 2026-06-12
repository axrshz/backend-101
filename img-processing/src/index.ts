import { Hono } from "hono";

import { db } from "./db";
import { users } from "./db/schema";

const app = new Hono();

app.get("/", async (c) => {
  const allUsers = await db.select().from(users);

  return c.json({
    message: "Bun, Hono, Drizzle, and PostgreSQL are working.",
    users: allUsers,
  });
});

export default {
  port: process.env.PORT ?? 3000,
  fetch: app.fetch,
};
