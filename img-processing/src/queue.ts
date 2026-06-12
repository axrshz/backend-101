import { Queue } from "bullmq";

const redisUrl = process.env.REDIS_URL;

if (!redisUrl) {
  throw new Error("REDIS_URL is required");
}

export type ImageJobData = {
  imageId: string;
};

export type ImageJobName = "process-image";

const parsedRedisUrl = new URL(redisUrl);

export const redisConnection = {
  host: parsedRedisUrl.hostname,
  port: Number(parsedRedisUrl.port || 6379),
  username: parsedRedisUrl.username || undefined,
  password: parsedRedisUrl.password || undefined,
  db: parsedRedisUrl.pathname
    ? Number(parsedRedisUrl.pathname.slice(1) || 0)
    : 0,
  tls: parsedRedisUrl.protocol === "rediss:" ? {} : undefined,
  maxRetriesPerRequest: null,
};

export const imageQueue = new Queue<ImageJobData, void, ImageJobName>(
  "image-processing",
  {
    connection: redisConnection,
    defaultJobOptions: {
      attempts: 3,
      backoff: {
        type: "exponential",
        delay: 1_000,
      },
      removeOnComplete: 100,
      removeOnFail: 500,
    },
  },
);
