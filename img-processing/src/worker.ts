import { s3, write } from "bun";
import { Worker } from "bullmq";
import { eq } from "drizzle-orm";
import sharp from "sharp";

import { db } from "./db";
import { images } from "./db/schema";
import type { ImageOperation } from "./operations";
import { redisConnection, type ImageJobData, type ImageJobName } from "./queue";

function processImage(input: Buffer, operation: ImageOperation) {
  const image = sharp(input);

  switch (operation) {
    case "grayscale":
      return image.grayscale().toBuffer();
    case "sharpen":
      return image.sharpen().toBuffer();
    case "blur":
      return image.blur().toBuffer();
    case "webp":
      return image.webp().toBuffer();
  }
}

const worker = new Worker<ImageJobData, void, ImageJobName>(
  "image-processing",
  async (job) => {
    const [image] = await db
      .select()
      .from(images)
      .where(eq(images.id, job.data.imageId))
      .limit(1);

    if (!image) {
      throw new Error(`Image ${job.data.imageId} not found`);
    }

    await db
      .update(images)
      .set({
        status: "processing",
        errorMessage: null,
        updatedAt: new Date(),
      })
      .where(eq(images.id, image.id));

    try {
      const original = await s3.file(image.originalS3Key).arrayBuffer();
      const operation = image.operation as ImageOperation;
      const processed = await processImage(Buffer.from(original), operation);
      const result =
        operation === "webp"
          ? processed
          : await sharp(processed).png().toBuffer();
      const extension = operation === "webp" ? "webp" : "png";
      const resultS3Key = `results/${image.id}.${extension}`;

      await write(s3.file(resultS3Key), result);

      await db
        .update(images)
        .set({
          status: "completed",
          resultS3Key,
          errorMessage: null,
          updatedAt: new Date(),
        })
        .where(eq(images.id, image.id));
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);

      await db
        .update(images)
        .set({
          status: "failed",
          errorMessage: message,
          updatedAt: new Date(),
        })
        .where(eq(images.id, image.id));

      throw error;
    }
  },
  {
    connection: redisConnection,
    concurrency: 2,
  },
);

worker.on("completed", (job) => {
  console.log(`Image job ${job.id} completed`);
});

worker.on("failed", (job, error) => {
  console.error(`Image job ${job?.id} failed:`, error);
});

console.log("Image worker started");

async function shutdown() {
  await worker.close();
}

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
