import {
  integer,
  pgTable,
  text,
  timestamp,
  pgEnum,
  uuid,
} from "drizzle-orm/pg-core";

export const imageStatus = pgEnum("image_status", [
  "queued",
  "processing",
  "completed",
  "failed",
  "deleted",
]);

export const images = pgTable("images", {
  id: uuid("id").primaryKey().defaultRandom(),
  operation: text("operation").notNull(),
  originalS3Key: text("original_s3_key").notNull(),
  resultS3Key: text("result_s3_key"),
  status: imageStatus("status").notNull().default("queued"),
  errorMessage: text("error_message"),
  createdAt: timestamp("created_at", { withTimezone: true })
    .notNull()
    .defaultNow(),
  updatedAt: timestamp("updated_at", { withTimezone: true })
    .notNull()
    .defaultNow(),
});
