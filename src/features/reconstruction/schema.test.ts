import { describe, expect, it } from "vitest";
import { sampleArtifact } from "./sampleArtifact";
import { artifactSchema } from "./schema";

describe("artifactSchema", () => {
  it("accepts the bundled sample artifact", () => {
    const parsed = artifactSchema.safeParse(sampleArtifact);

    expect(parsed.success).toBe(true);
    expect(sampleArtifact.vehicleTrack.length).toBeGreaterThan(3);
  });
});
