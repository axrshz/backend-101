import { Hono } from "hono";
import { s3, write } from "bun";

import { db } from "./db";
import { images } from "./db/schema";
import { isImageOperation, operations } from "./operations";
import { eq } from "drizzle-orm";
import { imageQueue } from "./queue";

const app = new Hono();

app.get("/", async (c) => {
  const html = await Bun.file("public/index.html").text();
  return c.html(html);
});

app.post("/images", async (c) => {
  const body = await c.req.parseBody();
  const file = body.image;
  const operation = body.operation;

  if (!file || typeof file === "string") {
    return c.json({ error: "No image file provided" }, 400);
  }

  if (!isImageOperation(operation)) {
    return c.json(
      {
        error: "Invalid operation",
        allowedOperations: operations,
      },
      400,
    );
  }

  const uniqueId = crypto.randomUUID();
  const fileExtension = file.name.split(".").pop();
  const s3Path = `images/${uniqueId}.${fileExtension}`;

  try {
    const s3Ref = s3.file(s3Path);
    await write(s3Ref, file);

    const [image] = await db
      .insert(images)
      .values({
        operation,
        originalS3Key: s3Path,
        status: "queued",
      })
      .returning({
        id: images.id,
        status: images.status,
        originalS3Key: images.originalS3Key,
      });

    await imageQueue.add(
      "process-image",
      { imageId: image.id },
      {
        jobId: image.id,
      },
    );

    return c.json(
      {
        id: image.id,
        status: image.status,
        originalS3Key: image.originalS3Key,
      },
      201,
    );
  } catch (error) {
    console.error("Image creation error:", error);
    return c.json({ error: "Failed to create image" }, 500);
  }
});

app.get("/images/:id", async (c) => {
  try {
    const id = c.req.param("id");
    const [image] = await db
      .select({
        id: images.id,
        status: images.status,
        resultS3Key: images.resultS3Key,
        errorMessage: images.errorMessage,
      })
      .from(images)
      .where(eq(images.id, id))
      .limit(1);

    if (!image) {
      return c.json({ error: "Image not found" }, 404);
    }

    const resultUrl =
      image.status === "completed" && image.resultS3Key
        ? s3.file(image.resultS3Key).presign({
            method: "GET",
            expiresIn: 60 * 60,
          })
        : null;

    return c.json({
      id: image.id,
      status: image.status,
      resultUrl,
      errorMessage: image.status === "failed" ? image.errorMessage : null,
    });
  } catch (error) {
    console.error("Image lookup error:", error);
    return c.json({ error: "Failed to get image" }, 500);
  }
});

export default {
  port: process.env.PORT ?? 3000,
  fetch: app.fetch,
};
