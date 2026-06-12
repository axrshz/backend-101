CREATE TYPE "public"."image_status" AS ENUM('queued', 'processing', 'completed', 'failed', 'deleted');--> statement-breakpoint
CREATE TABLE "images" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"operation" text NOT NULL,
	"original_s3_key" text NOT NULL,
	"result_s3_key" text,
	"status" "image_status" DEFAULT 'queued' NOT NULL,
	"error_message" text,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"updated_at" timestamp with time zone DEFAULT now() NOT NULL
);
