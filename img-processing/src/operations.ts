export const operations = ["grayscale", "sharpen", "blur", "webp"] as const;

export type ImageOperation = (typeof operations)[number];

export function isImageOperation(value: unknown): value is ImageOperation {
  return (
    typeof value === "string" &&
    operations.includes(value as ImageOperation)
  );
}
